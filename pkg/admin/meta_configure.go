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
	"errors"

	appsvr "github.com/bhojpur/application/pkg/engine"
)

// metaConfig meta config
type metaConfig struct {
}

// GetTemplate get customized template for meta
func (metaConfig) GetTemplate(context *Context, metaType string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

var defaultMetaConfigureMaps = map[string]func(*Meta){
	"date": func(meta *Meta) {
		if _, ok := meta.Config.(*DatetimeConfig); !ok || meta.Config == nil {
			meta.Config = &DatetimeConfig{}
			meta.Config.ConfigureBhojpurMeta(meta)
		}
	},

	"datetime": func(meta *Meta) {
		if _, ok := meta.Config.(*DatetimeConfig); !ok || meta.Config == nil {
			meta.Config = &DatetimeConfig{ShowTime: true}
			meta.Config.ConfigureBhojpurMeta(meta)
		}
	},

	"string": func(meta *Meta) {
		if meta.FormattedValuer == nil {
			meta.SetFormattedValuer(func(value interface{}, context *appsvr.Context) interface{} {
				switch str := meta.GetValuer()(value, context).(type) {
				case *string:
					if str != nil {
						return *str
					}
					return ""
				case string:
					return str
				default:
					return str
				}
			})
		}
	},

	"text": func(meta *Meta) {
		if meta.FormattedValuer == nil {
			meta.SetFormattedValuer(func(value interface{}, context *appsvr.Context) interface{} {
				switch str := meta.GetValuer()(value, context).(type) {
				case *string:
					if str != nil {
						return *str
					}
					return ""
				case string:
					return str
				default:
					return str
				}
			})
		}
	},

	"select_one": func(meta *Meta) {
		if metaConfig, ok := meta.Config.(*SelectOneConfig); !ok || metaConfig == nil {
			meta.Config = &SelectOneConfig{Collection: meta.Collection}
			meta.Config.ConfigureBhojpurMeta(meta)
		} else if meta.Collection != nil {
			metaConfig.Collection = meta.Collection
			meta.Config.ConfigureBhojpurMeta(meta)
		}
	},

	"select_many": func(meta *Meta) {
		if metaConfig, ok := meta.Config.(*SelectManyConfig); !ok || metaConfig == nil {
			meta.Config = &SelectManyConfig{Collection: meta.Collection}
			meta.Config.ConfigureBhojpurMeta(meta)
		} else if meta.Collection != nil {
			metaConfig.Collection = meta.Collection
			meta.Config.ConfigureBhojpurMeta(meta)
		}
	},

	"single_edit": func(meta *Meta) {
		if _, ok := meta.Config.(*SingleEditConfig); !ok || meta.Config == nil {
			meta.Config = &SingleEditConfig{}
			meta.Config.ConfigureBhojpurMeta(meta)
		}
	},

	"collection_edit": func(meta *Meta) {
		if _, ok := meta.Config.(*CollectionEditConfig); !ok || meta.Config == nil {
			meta.Config = &CollectionEditConfig{}
			meta.Config.ConfigureBhojpurMeta(meta)
		}
	},

	"rich_editor": func(meta *Meta) {
		if _, ok := meta.Config.(*RichEditorConfig); !ok || meta.Config == nil {
			meta.Config = &RichEditorConfig{}
			meta.Config.ConfigureBhojpurMeta(meta)
		}
	},
}
