package publish2

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

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/cms/pkg/admin"
)

type controller struct {
	Resource *admin.Resource
}

type visiblePublishResourceInterface interface {
	VisiblePublishResource(*appsvr.Context) bool
}

func (ctr controller) Dashboard(context *admin.Context) {
	type resourceResult struct {
		Resource            *admin.Resource
		ComingOnlineResults interface{}
		GoingOfflineResults interface{}
	}

	var results = []resourceResult{}

	for _, res := range context.Admin.GetResources() {
		if IsSchedulableModel(res.Value) {
			if visibleInterface, ok := res.Value.(visiblePublishResourceInterface); ok {
				if !visibleInterface.VisiblePublishResource(context.Context) {
					continue
				}
			} else if res.Config.Invisible {
				continue
			}

			db := context.GetDB()
			result := resourceResult{Resource: res}

			comingOnlineData := res.NewSlice()
			if db.Set(VisibleMode, "on").Set(ScheduleMode, ComingOnlineMode).Set(VersionMode, VersionMultipleMode).Find(comingOnlineData).RowsAffected > 0 {
				result.ComingOnlineResults = comingOnlineData
			}

			goingOfflineData := res.NewSlice()
			if db.Set(VisibleMode, "on").Set(ScheduleMode, GoingOfflineMode).Set(VersionMode, VersionMultipleMode).Find(goingOfflineData).RowsAffected > 0 {
				result.GoingOfflineResults = goingOfflineData
			}

			if result.ComingOnlineResults != nil || result.GoingOfflineResults != nil {
				results = append(results, result)
			}
		}
	}

	context.Action = "index"
	context.Execute("publish2/dashboard", results)
}

func (ctr controller) Versions(context *admin.Context) {
	records := context.Resource.NewSlice()
	record := context.Resource.NewStruct()
	primaryQuerySQL, primaryParams := ctr.Resource.ToPrimaryQueryParams(context.ResourceID, context.Context)
	tx := context.GetDB().Set(admin.DisableCompositePrimaryKeyMode, "on").Set(VersionMode, VersionMultipleMode).Set(ScheduleMode, ModeOff).Set(VisibleMode, ModeOff)
	tx.Where(primaryQuerySQL, primaryParams...).First(record)

	scope := tx.NewScope(record)
	tx.Find(records, fmt.Sprintf("%v = ?", scope.PrimaryKey()), scope.PrimaryKeyValue())

	result := context.Funcs(template.FuncMap{
		"version_metas": func() (metas []*admin.Meta) {
			for _, name := range []string{"VersionName", "ScheduledStartAt", "ScheduledEndAt", "PublishReady", "PublishLiveNow"} {
				if meta := ctr.Resource.GetMeta(name); meta != nil && meta.HasPermission(roles.Read, context.Context) {
					metas = append(metas, meta)
				}
			}
			return
		},
	}).Render("publish2/versions", records)
	context.Writer.Write([]byte(result))
}
