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
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/bhojpur/application/pkg/utils"
)

// Render find widget by name, render it based on current context
func (widgets *Widgets) Render(widgetName string, widgetGroupName string) template.HTML {
	return widgets.NewContext(nil).Render(widgetName, widgetGroupName)
}

// NewContext create new context for widgets
func (widgets *Widgets) NewContext(context *Context) *Context {
	if context == nil {
		context = &Context{}
	}

	if context.DB == nil {
		context.DB = widgets.Config.DB
	}

	if context.Options == nil {
		context.Options = map[string]interface{}{}
	}

	if context.FuncMaps == nil {
		context.FuncMaps = template.FuncMap{}
	}

	for key, fc := range widgets.funcMaps {
		if _, ok := context.FuncMaps[key]; !ok {
			context.FuncMaps[key] = fc
		}
	}

	context.Widgets = widgets
	return context
}

// Funcs return view functions map
func (context *Context) Funcs(funcMaps template.FuncMap) *Context {
	if context.FuncMaps == nil {
		context.FuncMaps = template.FuncMap{}
	}

	for key, fc := range funcMaps {
		context.FuncMaps[key] = fc
	}

	return context
}

// FuncMap return funcmap
func (context *Context) FuncMap() template.FuncMap {
	funcMap := template.FuncMap{}

	funcMap["render_widget"] = func(widgetName string, widgetGroupName ...string) template.HTML {
		var groupName string
		if len(widgetGroupName) == 0 {
			groupName = ""
		} else {
			groupName = widgetGroupName[0]
		}
		return context.Render(widgetName, groupName)
	}

	return funcMap
}

// Render register widget itself content
func (w *Widget) Render(context *Context, file string) template.HTML {
	if len(context.Body) > 0 {
		return template.HTML(context.Body)
	}

	var (
		err     error
		content []byte
		tmpl    *template.Template
	)

	if file == "" && len(w.Templates) > 0 {
		file = w.Templates[0]
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Get error when render file %v: %v", file, r)
			utils.ExitWithMsg(err)
		}
	}()

	if content, err = context.Widgets.AssetFS.Asset(file + ".tmpl"); err == nil {
		if tmpl, err = template.New(filepath.Base(file)).Funcs(context.FuncMaps).Parse(string(content)); err == nil {
			var result = bytes.NewBufferString("")
			if err = tmpl.Execute(result, context.Options); err == nil {
				return template.HTML(result.String())
			}
		}
	}

	return template.HTML(err.Error())
}

// RegisterViewPath register views directory
func (widgets *Widgets) RegisterViewPath(p string) {
	if filepath.IsAbs(p) {
		viewPaths = append(viewPaths, p)
		widgets.AssetFS.RegisterPath(p)
	} else {
		for _, gopath := range utils.GOPATH() {
			viewPaths = append(viewPaths, filepath.Join(gopath, "src", p))
			widgets.AssetFS.RegisterPath(filepath.Join(gopath, "src", p))
		}
	}
}

// LoadPreviewAssets will return assets tag used for Preview
func (widgets *Widgets) LoadPreviewAssets() template.HTML {
	tags := ""
	for _, asset := range widgets.Config.PreviewAssets {
		extension := filepath.Ext(asset)
		if extension == ".css" {
			tags += fmt.Sprintf("<link rel=\"stylesheet\" type=\"text/css\" href=\"%v\">\n", asset)
		} else if extension == ".js" {
			tags += fmt.Sprintf("<script src=\"%v\"></script>\n", asset)
		} else {
			tags += fmt.Sprintf("%v\n", asset)
		}
	}
	return template.HTML(tags)
}
