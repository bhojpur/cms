package render

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
	"net/http"
	"path/filepath"
	"strings"
)

// Template template struct
type Template struct {
	render             *Render
	layout             string
	usingDefaultLayout bool
	funcMap            template.FuncMap
}

// FuncMap get func maps from tmpl
func (tmpl *Template) funcMapMaker(req *http.Request, writer http.ResponseWriter) template.FuncMap {
	var funcMap = template.FuncMap{}

	for key, fc := range tmpl.render.funcMaps {
		funcMap[key] = fc
	}

	if tmpl.render.Config.FuncMapMaker != nil {
		for key, fc := range tmpl.render.Config.FuncMapMaker(tmpl.render, req, writer) {
			funcMap[key] = fc
		}
	}

	for key, fc := range tmpl.funcMap {
		funcMap[key] = fc
	}
	return funcMap
}

// Funcs register Funcs for tmpl
func (tmpl *Template) Funcs(funcMap template.FuncMap) *Template {
	tmpl.funcMap = funcMap
	return tmpl
}

// Render render tmpl
func (tmpl *Template) Render(templateName string, obj interface{}, request *http.Request, writer http.ResponseWriter) (template.HTML, error) {
	var (
		content []byte
		t       *template.Template
		err     error
		funcMap = tmpl.funcMapMaker(request, writer)
		render  = func(name string, objs ...interface{}) (template.HTML, error) {
			var (
				err           error
				renderObj     interface{}
				renderContent []byte
			)

			if len(objs) == 0 {
				// default obj
				renderObj = obj
			} else {
				// overwrite obj
				for _, o := range objs {
					renderObj = o
					break
				}
			}

			if renderContent, err = tmpl.findTemplate(name); err == nil {
				var partialTemplate *template.Template
				result := bytes.NewBufferString("")
				if partialTemplate, err = template.New(filepath.Base(name)).Funcs(funcMap).Parse(string(renderContent)); err == nil {
					if err = partialTemplate.Execute(result, renderObj); err == nil {
						return template.HTML(result.String()), err
					}
				}
			} else {
				err = fmt.Errorf("failed to find template: %v", name)
			}

			if err != nil {
				fmt.Println(err)
			}

			return "", err
		}
	)

	// funcMaps
	funcMap["render"] = render
	funcMap["yield"] = func() (template.HTML, error) { return render(templateName) }

	layout := tmpl.layout
	usingDefaultLayout := false

	if layout == "" && tmpl.usingDefaultLayout {
		usingDefaultLayout = true
		layout = tmpl.render.DefaultLayout
	}

	var tpl bytes.Buffer

	if layout != "" {
		content, err = tmpl.findTemplate(filepath.Join("layouts", layout))
		if err == nil {
			t, err = template.New("").Funcs(funcMap).Parse(string(content))
			if err != nil {
				goto OnError
			}

			err = t.Execute(&tpl, obj)
			if err != nil {
				goto OnError
			}

			return template.HTML(tpl.String()), nil
		} else if !usingDefaultLayout {
			goto OnError
		}
	}

	content, err = tmpl.findTemplate(templateName)
	if err != nil {
		goto OnError
	}
	t, err = template.New("").Funcs(funcMap).Parse(string(content))
	if err != nil {
		goto OnError
	}

	err = t.Execute(&tpl, obj)
	if err != nil {
		goto OnError
	}
	return template.HTML(tpl.String()), nil

OnError:
	err = fmt.Errorf("Failed to render page '%s' with template '%v.tmpl', got error: %v", templateName, filepath.Join("layouts", tmpl.layout), err)
	fmt.Println(err)
	// Display error in page directly
	return template.HTML(err.Error()), nil
}

// Execute execute tmpl
func (tmpl *Template) Execute(templateName string, obj interface{}, req *http.Request, w http.ResponseWriter) error {
	result, err := tmpl.Render(templateName, obj, req, w)
	if err == nil {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html")
		}

		_, err = w.Write([]byte(result))
	}
	return err
}

func (tmpl *Template) findTemplate(name string) ([]byte, error) {
	return tmpl.render.Asset(strings.TrimSpace(name) + ".tmpl")
}
