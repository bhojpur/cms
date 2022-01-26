package widget

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"html/template"
	"reflect"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/utils"
	"github.com/bhojpur/cms/pkg/admin"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// Context widget context
type Context struct {
	Widgets          *Widgets
	DB               *orm.DB
	AvailableWidgets []string
	Options          map[string]interface{}
	Body             string
	InlineEdit       bool
	SourceType       string
	SourceID         string
	FuncMaps         template.FuncMap
	WidgetSetting    BhojpurWidgetSettingInterface
}

// Get get option with name
func (context Context) Get(name string) (interface{}, bool) {
	if value, ok := context.Options[name]; ok {
		return value, true
	}

	return nil, false
}

// Set set option by name
func (context *Context) Set(name string, value interface{}) {
	if context.Options == nil {
		context.Options = map[string]interface{}{}
	}
	context.Options[name] = value
}

// GetDB set option by name
func (context *Context) GetDB() *orm.DB {
	if context.DB != nil {
		return context.DB
	}
	return context.Widgets.Config.DB
}

// Clone clone a context
func (context *Context) Clone() *Context {
	return &Context{
		Widgets:          context.Widgets,
		DB:               context.DB,
		AvailableWidgets: context.AvailableWidgets,
		Options:          context.Options,
		InlineEdit:       context.InlineEdit,
		FuncMaps:         context.FuncMaps,
		WidgetSetting:    context.WidgetSetting,
	}
}

// Render render widget based on context
func (context *Context) Render(widgetName string, widgetGroupName string) template.HTML {
	var (
		visibleScopes         []string
		widgets               = context.Widgets
		widgetSettingResource = widgets.WidgetSettingResource
		clone                 = context.Clone()
	)

	for _, scope := range registeredScopes {
		if scope.Visible(context) {
			visibleScopes = append(visibleScopes, scope.ToParam())
		}
	}

	if setting := context.findWidgetSetting(widgetName, append(visibleScopes, "default"), widgetGroupName); setting != nil {
		clone.WidgetSetting = setting
		adminContext := admin.Context{Admin: context.Widgets.Config.Admin, Context: &appsvr.Context{DB: context.DB}}

		var (
			widgetObj     = GetWidget(setting.GetSerializableArgumentKind())
			widgetSetting = widgetObj.Context(clone, setting.GetSerializableArgument(setting))
		)

		if clone.InlineEdit {
			prefix := widgets.Resource.GetAdmin().GetRouter().Prefix
			inlineEditURL := adminContext.URLFor(setting, widgetSettingResource)
			if widgetObj.InlineEditURL != nil {
				inlineEditURL = widgetObj.InlineEditURL(context)
			}

			return template.HTML(fmt.Sprintf(
				"<script data-prefix=\"%v\" src=\"%v/assets/javascripts/widget_check.js?theme=widget\"></script><div class=\"bhojpur-widget bhojpur-widget-%v\" data-widget-inline-edit-url=\"%v\" data-url=\"%v\">\n%v\n</div>",
				prefix,
				prefix,
				utils.ToParamString(widgetObj.Name),
				fmt.Sprintf("%v/%v/inline-edit", prefix, widgets.Resource.ToParam()),
				inlineEditURL,
				widgetObj.Render(widgetSetting, setting.GetTemplate()),
			))
		}

		return widgetObj.Render(widgetSetting, setting.GetTemplate())
	}

	return template.HTML("")
}

func (context *Context) findWidgetSetting(widgetName string, scopes []string, widgetGroupName string) BhojpurWidgetSettingInterface {
	var (
		db                    = context.GetDB()
		widgetSettingResource = context.Widgets.WidgetSettingResource
		setting               BhojpurWidgetSettingInterface
		settings              = widgetSettingResource.NewSlice()
	)

	if context.SourceID != "" {
		db.Order("source_id DESC").Where("name = ? AND scope IN (?) AND ((shared = ? AND source_type = ?) OR (source_type = ? AND source_id = ?))", widgetName, scopes, true, "", context.SourceType, context.SourceID).Find(settings)
	} else {
		db.Where("name = ? AND scope IN (?) AND source_type = ?", widgetName, scopes, "").Find(settings)
	}

	settingsValue := reflect.Indirect(reflect.ValueOf(settings))
	if settingsValue.Len() > 0 {
	OUTTER:
		for _, scope := range scopes {
			for i := 0; i < settingsValue.Len(); i++ {
				s := settingsValue.Index(i).Interface().(BhojpurWidgetSettingInterface)
				if s.GetScope() == scope {
					setting = s
					break OUTTER
				}
			}
		}
	}

	if context.SourceType == "" {
		if setting == nil {
			if widgetGroupName == "" {
				utils.ExitWithMsg("Widget: Can't Create Widget Without Widget Type")
				return nil
			}
			setting = widgetSettingResource.NewStruct().(BhojpurWidgetSettingInterface)
			setting.SetWidgetName(widgetName)
			setting.SetGroupName(widgetGroupName)
			setting.SetSerializableArgumentKind(widgetGroupName)
			db.Create(setting)
		} else if setting.GetGroupName() != widgetGroupName {
			setting.SetGroupName(widgetGroupName)
			db.Save(setting)
		}
	}

	return setting
}
