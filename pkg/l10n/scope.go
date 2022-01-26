package l10n

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
	"reflect"

	"github.com/bhojpur/application/pkg/utils"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// IsLocalizable return model is localizable or not
func IsLocalizable(scope *orm.Scope) (IsLocalizable bool) {
	if scope.GetModelStruct().ModelType == nil {
		return false
	}
	_, IsLocalizable = reflect.New(scope.GetModelStruct().ModelType).Interface().(l10nInterface)
	return
}

type localeCreatableInterface interface {
	CreatableFromLocale()
}

type localeCreatableInterface2 interface {
	LocaleCreatable()
}

func isLocaleCreatable(scope *orm.Scope) (ok bool) {
	if _, ok = reflect.New(scope.GetModelStruct().ModelType).Interface().(localeCreatableInterface); ok {
		return
	}
	_, ok = reflect.New(scope.GetModelStruct().ModelType).Interface().(localeCreatableInterface2)
	return
}

func setLocale(scope *orm.Scope, locale string) {
	for _, field := range scope.Fields() {
		if field.Name == "LanguageCode" {
			field.Set(locale)
		}
	}
}

func getQueryLocale(scope *orm.Scope) (locale string, isLocale bool) {
	if str, ok := scope.DB().Get("l10n:locale"); ok {
		if locale, ok := str.(string); ok && locale != "" {
			return locale, locale != Global
		}
	}
	return Global, false
}

func getLocale(scope *orm.Scope) (locale string, isLocale bool) {
	if str, ok := scope.DB().Get("l10n:localize_to"); ok {
		if locale, ok := str.(string); ok && locale != "" {
			return locale, locale != Global
		}
	}

	return getQueryLocale(scope)
}

func isSyncField(field *orm.StructField) bool {
	if _, ok := utils.ParseTagOption(field.Tag.Get("l10n"))["SYNC"]; ok {
		return true
	}
	return false
}

func syncColumns(scope *orm.Scope) (columns []string) {
	for _, field := range scope.GetModelStruct().StructFields {
		if isSyncField(field) {
			columns = append(columns, field.DBName)
		}
	}
	return
}
