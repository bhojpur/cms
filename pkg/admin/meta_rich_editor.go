package admin

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
	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/utils"
)

// RichEditorConfig rich editor meta config
type RichEditorConfig struct {
	AssetManager         *Resource
	DisableHTMLSanitizer bool
	Plugins              []RedactorPlugin
	Settings             map[string]interface{}
	metaConfig
}

// RedactorPlugin register redactor plugins into rich editor
type RedactorPlugin struct {
	Name   string
	Source string
}

// ConfigureBhojpurMeta configure rich editor meta
func (richEditorConfig *RichEditorConfig) ConfigureBhojpurMeta(metaor resource.Metaor) {
	if meta, ok := metaor.(*Meta); ok {
		meta.Type = "rich_editor"

		// Compatible with old rich editor setting
		if meta.Resource != nil {
			richEditorConfig.AssetManager = meta.Resource
			meta.Resource = nil
		}

		if !richEditorConfig.DisableHTMLSanitizer {
			setter := meta.GetSetter()
			meta.SetSetter(func(resource interface{}, metaValue *resource.MetaValue, context *appsvr.Context) {
				metaValue.Value = utils.HTMLSanitizer.Sanitize(utils.ToString(metaValue.Value))
				setter(resource, metaValue, context)
			})
		}

		if richEditorConfig.Settings == nil {
			richEditorConfig.Settings = map[string]interface{}{}
		}

		plugins := []string{"source"}
		for _, plugin := range richEditorConfig.Plugins {
			plugins = append(plugins, plugin.Name)
		}
		richEditorConfig.Settings["plugins"] = plugins
	}
}
