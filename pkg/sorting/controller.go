package sorting

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
	"net/http"
	"path"
	"strconv"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/cms/pkg/admin"
)

func updatePosition(context *admin.Context) {
	if result, err := context.FindOne(); err == nil {
		if position, ok := result.(sortingInterface); ok {
			if pos, err := strconv.Atoi(context.Request.Form.Get("to")); err == nil {
				var count int
				if _, ok := result.(sortingDescInterface); ok {
					var result = context.Resource.NewStruct()
					context.GetDB().Set("l10n:mode", "locale").Order("position DESC", true).First(result)
					count = result.(sortingInterface).GetPosition()
					pos = count - pos + 1
				}

				if MoveTo(context.GetDB(), position, pos) == nil {
					var pos = position.GetPosition()
					if _, ok := result.(sortingDescInterface); ok {
						pos = count - pos + 1
					}

					context.Writer.Write([]byte(fmt.Sprintf("%d", pos)))
					return
				}
			}
		}
	}
	context.Writer.WriteHeader(admin.HTTPUnprocessableEntity)
	context.Writer.Write([]byte("Error"))
}

// ConfigureBhojpurResource configure sorting for CMS admin
func (s *Sorting) ConfigureBhojpurResourceBeforeInitialize(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		res.UseTheme("sorting")

		if res.Permission == nil {
			res.Permission = roles.NewPermission()
		}

		role := res.Permission.Role
		if _, ok := role.Get("sorting_mode"); !ok {
			role.Register("sorting_mode", func(req *http.Request, currentUser interface{}) bool {
				return req.URL.Query().Get("sorting") != ""
			})
		}

		res.Meta(&admin.Meta{
			Name: "Position",
			Valuer: func(value interface{}, ctx *appsvr.Context) interface{} {
				db := ctx.GetDB()
				var count int
				var pos = value.(sortingInterface).GetPosition()

				if _, ok := modelValue(value).(sortingDescInterface); ok {
					if total, ok := db.Get("sorting_total_count"); ok {
						count = total.(int)
					} else {
						var result = res.NewStruct()
						db.New().Order("position DESC", true).First(result)
						count = result.(sortingInterface).GetPosition()
						db.InstantSet("sorting_total_count", count)
					}
					pos = count - pos + 1
				}

				primaryKey := ctx.GetDB().NewScope(value).PrimaryKeyValue()
				url := path.Join(ctx.Request.URL.Path, fmt.Sprintf("%v", primaryKey), "sorting/update_position")
				return template.HTML(fmt.Sprintf("<input type=\"number\" disabled class=\"bhojpur-sorting__position\" value=\"%v\" data-sorting-url=\"%v\" data-position=\"%v\">", pos, url, pos))
			},
			Permission: roles.Allow(roles.Read, "sorting_mode"),
		})
	}
}

func (s *Sorting) ConfigureBhojpurResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		Admin := res.GetAdmin()

		res.OverrideIndexAttrs(func() {
			res.IndexAttrs(res.IndexAttrs(), "Position")
		})

		res.OverrideNewAttrs(func() {
			res.NewAttrs(res.NewAttrs(), "-Position")
		})

		res.OverrideEditAttrs(func() {
			res.EditAttrs(res.EditAttrs(), "-Position")
		})

		res.OverrideShowAttrs(func() {
			res.ShowAttrs(res.ShowAttrs(), "-Position")
		})

		router := Admin.GetRouter()
		router.Post(fmt.Sprintf("/%v/%v/sorting/update_position", res.ToParam(), res.ParamIDName()), updatePosition)
	}
}
