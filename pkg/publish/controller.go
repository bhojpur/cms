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
	"strings"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/application/pkg/utils"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/worker"
)

const (
	// PublishPermission publish permission
	PublishPermission roles.PermissionMode = "publish"
)

type publishController struct {
	*Publish
}

type visiblePublishResourceInterface interface {
	VisiblePublishResource(*appsvr.Context) bool
}

func (pc *publishController) Preview(context *admin.Context) {
	type resource struct {
		*admin.Resource
		Value interface{}
	}

	var drafts = []resource{}

	draftDB := context.GetDB().Set(publishDraftMode, true).Unscoped()
	for _, res := range context.Admin.GetResources() {
		if visibleInterface, ok := res.Value.(visiblePublishResourceInterface); ok {
			if !visibleInterface.VisiblePublishResource(context.Context) {
				continue
			}
		} else if res.Config.Invisible {
			continue
		}

		if res.HasPermission(PublishPermission, context.Context) {
			results := res.NewSlice()
			if IsPublishableModel(res.Value) || IsPublishEvent(res.Value) {
				if pc.SearchHandler(draftDB.Where("publish_status = ?", DIRTY), context.Context).Find(results).RowsAffected > 0 {
					drafts = append(drafts, resource{
						Resource: res,
						Value:    results,
					})
				}
			}
		}
	}
	context.Execute("publish_drafts", drafts)
}

func (pc *publishController) Diff(context *admin.Context) {
	var (
		resourceID = context.Request.URL.Query().Get(":publish_unique_key")
		params     = strings.Split(resourceID, "__") // name__primary_keys
		res        = context.Admin.GetResource(params[0])
	)

	draft := res.NewStruct()
	pc.search(context.GetDB().Set(publishDraftMode, true), res, [][]string{params[1:]}).First(draft)

	production := res.NewStruct()
	pc.search(context.GetDB().Set(publishDraftMode, false), res, [][]string{params[1:]}).First(production)

	results := map[string]interface{}{"Production": production, "Draft": draft, "Resource": res}
	fmt.Fprintf(context.Writer, string(context.Render("publish_diff", results)))
}

func (pc *publishController) PublishOrDiscard(context *admin.Context) {
	var request = context.Request
	var ids = request.Form["checked_ids[]"]

	if scheduler := pc.Publish.WorkerScheduler; scheduler != nil {
		jobResource := scheduler.JobResource
		result := jobResource.NewStruct().(worker.BhojpurJobInterface)
		if request.Form.Get("publish_type") == "discard" {
			result.SetJob(scheduler.GetRegisteredJob("Discard"))
		} else {
			result.SetJob(scheduler.GetRegisteredJob("Publish"))
		}

		workerArgument := &BhojpurWorkerArgument{IDs: ids}
		if t, err := utils.ParseTime(request.Form.Get("scheduled_time"), context.Context); err == nil {
			workerArgument.ScheduleTime = &t
		}
		result.SetSerializableArgumentValue(workerArgument)

		jobResource.CallSave(result, context.Context)
		scheduler.AddJob(result)

		http.Redirect(context.Writer, context.Request, context.URLFor(jobResource), http.StatusFound)
	} else {
		records := pc.searchWithPublishIDs(context.GetDB().Set(publishDraftMode, true), context.Admin, ids)

		if request.Form.Get("publish_type") == "publish" {
			pc.Publish.Publish(records...)
		} else if request.Form.Get("publish_type") == "discard" {
			pc.Publish.Discard(records...)
		}

		http.Redirect(context.Writer, context.Request, context.Request.RequestURI, http.StatusFound)
	}
}

// ConfigureBhojpurResourceBeforeInitialize configure Bhojpur resource when initialize Bhojpur CMS admin
func (publish *Publish) ConfigureBhojpurResourceBeforeInitialize(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		res.GetAdmin().RegisterViewPath("github.com/bhojpur/cms/pkg/publish/views")
		res.UseTheme("publish")

		if event := res.GetAdmin().GetResource("PublishEvent"); event == nil {
			eventResource := res.GetAdmin().AddResource(&PublishEvent{}, &admin.Config{Invisible: true})
			eventResource.IndexAttrs("Name", "Description", "CreatedAt")
		}
	}
}

// ConfigureBhojpurResource configure Bhojpur resource for Bhojpur CMS admin
func (publish *Publish) ConfigureBhojpurResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		controller := publishController{publish}
		router := res.GetAdmin().GetRouter()
		router.Get(fmt.Sprintf("/%v/diff/:publish_unique_key", res.ToParam()), controller.Diff)
		router.Get(res.ToParam(), controller.Preview)
		router.Post(res.ToParam(), controller.PublishOrDiscard)

		res.GetAdmin().RegisterFuncMap("publish_unique_key", func(res *admin.Resource, record interface{}, context *admin.Context) string {
			var publishKeys = []string{res.ToParam()}
			var scope = publish.DB.NewScope(record)
			for _, primaryField := range scope.PrimaryFields() {
				publishKeys = append(publishKeys, fmt.Sprint(primaryField.Field.Interface()))
			}
			return strings.Join(publishKeys, "__")
		})

		res.GetAdmin().RegisterFuncMap("is_publish_event_resource", func(res *admin.Resource) bool {
			return IsPublishEvent(res.Value)
		})
	}
}
