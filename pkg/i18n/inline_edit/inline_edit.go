package inline_edit

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

	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/i18n"
)

func init() {
	admin.RegisterViewPath("github.com/bhojpur/cms/pkg/i18n/inline_edit/views")
}

// FuncMap generate func map for inline edit
func FuncMap(I18n *i18n.I18N, locale string, enableInlineEdit bool) template.FuncMap {
	return template.FuncMap{
		"t": InlineEdit(I18n, locale, enableInlineEdit),
	}
}

// InlineEdit enable inline edit
func InlineEdit(I18n *i18n.I18N, locale string, isInline bool) func(string, ...interface{}) template.HTML {
	return func(key string, args ...interface{}) template.HTML {
		// Get Translation Value
		var value template.HTML
		var defaultValue string
		if len(args) > 0 {
			if args[0] == nil {
				defaultValue = key
			} else {
				defaultValue = fmt.Sprint(args[0])
			}
			value = I18N.Default(defaultValue).T(locale, key, args[1:]...)
		} else {
			value = I18n.T(locale, key)
		}

		// Append inline-edit script/tag
		if isInline {
			var editType string
			if len(value) > 25 {
				editType = "data-type=\"textarea\""
			}
			prefix := I18n.Resource.GetAdmin().GetRouter().Prefix
			assetsTag := fmt.Sprintf("<script data-prefix=\"%v\" src=\"%v/assets/javascripts/i18n-checker.js?theme=i18n\"></script>", prefix, prefix)
			return template.HTML(fmt.Sprintf("%s<span class=\"bhojpur-i18n-inline\" %s data-locale=\"%s\" data-key=\"%s\">%s</span>", assetsTag, editType, locale, key, string(value)))
		}
		return value
	}
}
