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
	"time"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/utils"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/serializable_meta"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// BhojpurWidgetSettingInterface CMS widget setting interface
type BhojpurWidgetSettingInterface interface {
	GetWidgetName() string
	SetWidgetName(string)
	GetGroupName() string
	SetGroupName(string)
	GetScope() string
	SetScope(string)
	GetTemplate() string
	SetTemplate(string)
	GetSourceType() string
	SetSourceType(string)
	GetSourceID() string
	SetSourceID(string)
	GetShared() bool
	SetShared(bool)
	serializable_meta.SerializableMetaInterface
}

// BhojpurWidgetSetting default CMS widget setting struct
type BhojpurWidgetSetting struct {
	Name        string `orm:"primary_key;size:100"`
	Scope       string `orm:"primary_key;size:128;default:'default'"`
	SourceType  string `orm:"primary_key;default:'';size:100"`
	SourceID    string `orm:"primary_key;default:'';size:100"`
	Description string
	Shared      bool
	WidgetType  string
	GroupName   string
	Template    string
	serializable_meta.SerializableMeta
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ResourceName get widget setting's resource name
func (widgetSetting *BhojpurWidgetSetting) ResourceName() string {
	return "Widget Content"
}

// GetSerializableArgumentKind get serializable kind
func (widgetSetting *BhojpurWidgetSetting) GetSerializableArgumentKind() string {
	if widgetSetting.WidgetType != "" {
		return widgetSetting.WidgetType
	}
	return widgetSetting.Kind
}

// SetSerializableArgumentKind set serializable kind
func (widgetSetting *BhojpurWidgetSetting) SetSerializableArgumentKind(name string) {
	widgetSetting.WidgetType = name
	widgetSetting.Kind = name
}

// GetWidgetName get widget setting's group name
func (widgetSetting BhojpurWidgetSetting) GetWidgetName() string {
	return widgetSetting.Name
}

// SetWidgetName set widget setting's group name
func (widgetSetting *BhojpurWidgetSetting) SetWidgetName(name string) {
	widgetSetting.Name = name
}

// GetGroupName get widget setting's group name
func (widgetSetting BhojpurWidgetSetting) GetGroupName() string {
	return widgetSetting.GroupName
}

// SetGroupName set widget setting's group name
func (widgetSetting *BhojpurWidgetSetting) SetGroupName(groupName string) {
	widgetSetting.GroupName = groupName
}

// GetScope get widget's scope
func (widgetSetting BhojpurWidgetSetting) GetScope() string {
	return widgetSetting.Scope
}

// SetScope set widget setting's scope
func (widgetSetting *BhojpurWidgetSetting) SetScope(scope string) {
	widgetSetting.Scope = scope
}

// GetSourceType get widget's source type
func (widgetSetting BhojpurWidgetSetting) GetSourceType() string {
	return widgetSetting.SourceType
}

// SetSourceType set widget setting's souce type
func (widgetSetting *BhojpurWidgetSetting) SetSourceType(sourceType string) {
	widgetSetting.SourceType = sourceType
}

// GetSourceID get widget's source ID
func (widgetSetting BhojpurWidgetSetting) GetSourceID() string {
	return widgetSetting.SourceID
}

// SetSourceID set widget setting's source id
func (widgetSetting *BhojpurWidgetSetting) SetSourceID(sourceID string) {
	widgetSetting.SourceID = sourceID
}

// GetShared get widget's source ID
func (widgetSetting BhojpurWidgetSetting) GetShared() bool {
	return widgetSetting.Shared
}

// SetShared set widget setting's source id
func (widgetSetting *BhojpurWidgetSetting) SetShared(shared bool) {
	widgetSetting.Shared = shared
}

// GetTemplate get used widget template
func (widgetSetting BhojpurWidgetSetting) GetTemplate() string {
	if widget := GetWidget(widgetSetting.GetSerializableArgumentKind()); widget != nil {
		for _, value := range widget.Templates {
			if value == widgetSetting.Template {
				return value
			}
		}

		// return first value of defined widget templates
		for _, value := range widget.Templates {
			return value
		}
	}
	return ""
}

// SetTemplate set used widget's template
func (widgetSetting *BhojpurWidgetSetting) SetTemplate(template string) {
	widgetSetting.Template = template
}

// GetSerializableArgumentResource get setting's argument's resource
func (widgetSetting *BhojpurWidgetSetting) GetSerializableArgumentResource() *admin.Resource {
	widget := GetWidget(widgetSetting.GetSerializableArgumentKind())
	if widget != nil {
		return widget.Setting
	}
	return nil
}

// ConfigureBhojpurResource a method used to config Widget for CMS admin
func (widgetSetting *BhojpurWidgetSetting) ConfigureBhojpurResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		res.Meta(&admin.Meta{Name: "PreviewIcon", Valuer: func(result interface{}, context *appsvr.Context) interface{} {
			if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
				if widget := GetWidget(setting.GetSerializableArgumentKind()); widget != nil {
					return template.HTML(fmt.Sprintf("<img class='bhojpur-preview-icon' src='%v'/>", widget.PreviewIcon))
				}
			}
			return ""
		}})

		res.Meta(&admin.Meta{Name: "Name", Type: "string"})
		res.Meta(&admin.Meta{Name: "DisplayName", Label: "Name", Type: "readonly", FieldName: "Name"})
		res.Meta(&admin.Meta{Name: "Description", Type: "string"})

		res.Meta(&admin.Meta{
			Name: "Scope",
			Type: "hidden",
			Valuer: func(result interface{}, context *appsvr.Context) interface{} {
				if scope := context.Request.URL.Query().Get("widget_scope"); scope != "" {
					return scope
				}

				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					if scope := setting.GetScope(); scope != "" {
						return scope
					}
				}

				return "default"
			},
			Setter: func(result interface{}, metaValue *resource.MetaValue, context *appsvr.Context) {
				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					setting.SetScope(utils.ToString(metaValue.Value))
				}
			},
		})

		res.Meta(&admin.Meta{
			Name: "SourceType",
			Type: "hidden",
			Valuer: func(result interface{}, context *appsvr.Context) interface{} {
				if sourceType := context.Request.URL.Query().Get("widget_source_type"); sourceType != "" {
					return sourceType
				}

				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					if sourceType := setting.GetSourceType(); sourceType != "" {
						return sourceType
					}
				}
				return ""
			},
			Setter: func(result interface{}, metaValue *resource.MetaValue, context *appsvr.Context) {
				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					setting.SetSourceType(utils.ToString(metaValue.Value))
				}
			},
		})

		res.Meta(&admin.Meta{
			Name: "SourceID",
			Type: "hidden",
			Valuer: func(result interface{}, context *appsvr.Context) interface{} {
				if sourceID := context.Request.URL.Query().Get("widget_source_id"); sourceID != "" {
					return sourceID
				}

				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					if sourceID := setting.GetSourceID(); sourceID != "" {
						return sourceID
					}
				}
				return ""
			},
			Setter: func(result interface{}, metaValue *resource.MetaValue, context *appsvr.Context) {
				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					setting.SetSourceID(utils.ToString(metaValue.Value))
				}
			},
		})

		res.Meta(&admin.Meta{
			Name: "Widgets",
			Type: "select_one",
			Valuer: func(result interface{}, context *appsvr.Context) interface{} {
				if typ := context.Request.URL.Query().Get("widget_type"); typ != "" {
					return typ
				}

				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					widget := GetWidget(setting.GetSerializableArgumentKind())
					if widget == nil {
						return ""
					}
					return widget.Name
				}

				return ""
			},
			Collection: func(result interface{}, context *appsvr.Context) (results [][]string) {
				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					if setting.GetWidgetName() == "" {
						for _, widget := range registeredWidgets {
							results = append(results, []string{widget.Name, widget.Name})
						}
					} else {
						groupName := setting.GetGroupName()
						for _, group := range registeredWidgetsGroup {
							if group.Name == groupName {
								for _, widget := range group.Widgets {
									results = append(results, []string{widget, widget})
								}
							}
						}
					}

					if len(results) == 0 {
						results = append(results, []string{setting.GetSerializableArgumentKind(), setting.GetSerializableArgumentKind()})
					}
				}
				return
			},
			Setter: func(result interface{}, metaValue *resource.MetaValue, context *appsvr.Context) {
				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					setting.SetSerializableArgumentKind(utils.ToString(metaValue.Value))
				}
			},
		})

		res.Meta(&admin.Meta{
			Name: "Template",
			Type: "select_one",
			Valuer: func(result interface{}, context *appsvr.Context) interface{} {
				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					return setting.GetTemplate()
				}
				return ""
			},
			Collection: func(result interface{}, context *appsvr.Context) (results [][]string) {
				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					if widget := GetWidget(setting.GetSerializableArgumentKind()); widget != nil {
						for _, value := range widget.Templates {
							results = append(results, []string{value, value})
						}
					}
				}
				return
			},
			Setter: func(result interface{}, metaValue *resource.MetaValue, context *appsvr.Context) {
				if setting, ok := result.(BhojpurWidgetSettingInterface); ok {
					setting.SetTemplate(utils.ToString(metaValue.Value))
				}
			},
		})

		res.Meta(&admin.Meta{
			Name:  "Shared",
			Label: "Add to Container Library (can be reused on other pages)",
		})

		res.Action(&admin.Action{
			Name: "Preview",
			URL: func(record interface{}, context *admin.Context) string {
				return fmt.Sprintf("%v/%v/%v/!preview", context.Admin.GetRouter().Prefix, res.ToParam(), record.(BhojpurWidgetSettingInterface).GetWidgetName())
			},
			Modes: []string{"edit", "menu_item"},
		})

		res.AddProcessor(&resource.Processor{
			Name: "widget-default-value",
			Handler: func(value interface{}, metaValues *resource.MetaValues, context *appsvr.Context) error {
				if widgetSetting, ok := value.(BhojpurWidgetSettingInterface); ok {
					if widgetSetting.GetShared() {
						widgetSetting.SetSourceType("")
						widgetSetting.SetSourceID("")
					}
				}
				return nil
			},
		})

		res.UseTheme("widget")

		res.IndexAttrs("PreviewIcon", "Name", "Description", "CreatedAt", "UpdatedAt")
		res.OverrideShowAttrs(func() {
			res.ShowAttrs("PreviewIcon", "Name", "Scope", "WidgetType", "Template", "Description", "Value", "CreatedAt", "UpdatedAt")
		})
		res.EditAttrs(
			"DisplayName", "Description", "Scope", "Widgets", "Template",
			&admin.Section{
				Title: "Settings",
				Rows:  [][]string{{"Kind"}, {"SerializableMeta"}},
			},
			"Shared", "SourceType", "SourceID",
		)
		res.NewAttrs("Name", "Description", "Scope", "Widgets", "Template",
			&admin.Section{
				Title: "Settings",
				Rows:  [][]string{{"Kind"}, {"SerializableMeta"}},
			},
			"Shared", "SourceType", "SourceID",
		)

		searchHandler := res.SearchHandler
		res.SearchHandler = func(keyword string, context *appsvr.Context) *orm.DB {
			// don't include widgets have source_type in index page
			if context.ResourceID == "" {
				context.SetDB(context.GetDB().Where("source_type = ?", ""))
			}
			return searchHandler(keyword, context)
		}
	}
}
