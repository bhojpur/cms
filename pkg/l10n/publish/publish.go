package publish

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
	"net/http"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/l10n"
	"github.com/bhojpur/cms/pkg/publish"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type availableLocalesInterface interface {
	AvailableLocales() []string
}

type publishableLocalesInterface interface {
	PublishableLocales() []string
}

type editableLocalesInterface interface {
	EditableLocales() []string
}

func getPublishableLocales(req *http.Request, currentUser interface{}) []string {
	if user, ok := currentUser.(publishableLocalesInterface); ok {
		return user.PublishableLocales()
	}

	if user, ok := currentUser.(editableLocalesInterface); ok {
		return user.EditableLocales()
	}

	if user, ok := currentUser.(availableLocalesInterface); ok {
		return user.AvailableLocales()
	}
	return []string{l10n.Global}
}

// RegisterL10nForPublish register l10n language switcher for publish
func RegisterL10nForPublish(Publish *publish.Publish, Admin *admin.Admin) {
	searchHandler := Publish.SearchHandler
	Publish.SearchHandler = func(db *orm.DB, context *appsvr.Context) *orm.DB {
		if context != nil {
			if context.Request != nil && context.Request.URL.Query().Get("locale") == "" {
				publishableLocales := getPublishableLocales(context.Request, context.CurrentUser)
				return searchHandler(db, context).Set("l10n:mode", "unscoped").Scopes(func(db *orm.DB) *orm.DB {
					scope := db.NewScope(db.Value)
					if l10n.IsLocalizable(scope) {
						return db.Where(fmt.Sprintf("%v.language_code IN (?)", scope.QuotedTableName()), publishableLocales)
					}
					return db
				})
			}
			return searchHandler(db, context).Set("l10n:mode", "locale")
		}
		return searchHandler(db, context).Set("l10n:mode", "unscoped")
	}

	Admin.RegisterViewPath("github.com/bhojpur/cms/pkg/l10n/publish/views")

	Admin.RegisterFuncMap("publishable_locales", func(context admin.Context) []string {
		return getPublishableLocales(context.Request, context.CurrentUser)
	})
}
