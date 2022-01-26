package worker

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
	"html/template"
	"net/http"

	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/responder"
)

type workerController struct {
	*Worker
}

func (wc workerController) Index(context *admin.Context) {
	context = context.NewResourceContext(wc.JobResource)
	result, err := context.FindMany()
	context.AddError(err)

	if context.HasError() {
		http.NotFound(context.Writer, context.Request)
	} else {
		responder.With("html", func() {
			context.Execute("index", result)
		}).With("json", func() {
			context.JSON("index", result)
		}).Respond(context.Request)
	}
}

func (wc workerController) Show(context *admin.Context) {
	job, err := wc.GetJob(context.ResourceID)
	context.AddError(err)
	context.Execute("show", job)
}

func (wc workerController) New(context *admin.Context) {
	context.Execute("new", wc.Worker)
}

func (wc workerController) Update(context *admin.Context) {
	if job, err := wc.GetJob(context.ResourceID); err == nil {
		if job.GetStatus() == JobStatusScheduled || job.GetStatus() == JobStatusNew {
			if job.GetJob().HasPermission(roles.Update, context.Context) {
				if context.AddError(wc.Worker.JobResource.Decode(context.Context, job)); !context.HasError() {
					context.AddError(wc.Worker.JobResource.CallSave(job, context.Context))
					context.AddError(wc.Worker.AddJob(job))
				}

				if !context.HasError() {
					context.Flash(string(context.Admin.T(context.Context, "bhojpur_worker.form.successfully_updated", "{{.Name}} was successfully updated", wc.Worker.JobResource)), "success")
				}

				context.Execute("edit", job)
				return
			}
		}

		context.AddError(errors.New("not allowed to update this job"))
	} else {
		context.AddError(err)
	}

	http.Redirect(context.Writer, context.Request, context.Request.URL.Path, http.StatusFound)
}

func (wc workerController) AddJob(context *admin.Context) {
	jobResource := wc.Worker.JobResource
	result := jobResource.NewStruct().(BhojpurJobInterface)
	job := wc.Worker.GetRegisteredJob(context.Request.Form.Get("job_name"))
	result.SetJob(job)

	if !job.HasPermission(roles.Create, context.Context) {
		context.AddError(errors.New("don't have permission to run job"))
	}

	if context.AddError(jobResource.Decode(context.Context, result)); !context.HasError() {
		// ensure job name is correct
		result.SetJob(job)
		context.AddError(jobResource.CallSave(result, context.Context))
		context.AddError(wc.Worker.AddJob(result))
	}

	if context.HasError() {
		responder.With("html", func() {
			context.Writer.WriteHeader(422)
			context.Execute("edit", result)
		}).With("json", func() {
			context.Writer.WriteHeader(422)
			context.JSON("index", map[string]interface{}{"errors": context.GetErrors()})
		}).Respond(context.Request)
		return
	}

	context.Flash(string(context.Admin.T(context.Context, "bhojpur_worker.form.successfully_created", "{{.Name}} was successfully created", jobResource)), "success")
	http.Redirect(context.Writer, context.Request, context.Request.URL.Path, http.StatusFound)
}

func (wc workerController) RunJob(context *admin.Context) {
	if newJob := wc.Worker.saveAnotherJob(context.ResourceID); newJob != nil {
		wc.Worker.AddJob(newJob)
	} else {
		context.AddError(errors.New("failed to clone job " + context.ResourceID))
	}

	http.Redirect(context.Writer, context.Request, context.URLFor(wc.Worker.JobResource), http.StatusFound)
}

func (wc workerController) KillJob(context *admin.Context) {
	var msg template.HTML
	bhojpurJob, err := wc.Worker.GetJob(context.ResourceID)

	if err == nil {
		if err = wc.Worker.KillJob(bhojpurJob.GetJobID()); err == nil {
			msg = context.Admin.T(context.Context, "bhojpur_worker.form.successfully_killed", "{{.Name}} was successfully killed", wc.JobResource)
		} else {
			msg = context.Admin.T(context.Context, "bhojpur_worker.form.failed_to_kill", "Failed to kill job {{.Name}}", wc.JobResource)
		}
	}

	if err == nil {
		responder.With("html", func() {
			context.Flash(string(msg), "success")
			http.Redirect(context.Writer, context.Request, context.Request.URL.Path, http.StatusFound)
		}).With("json", func() {
			context.JSON("ok", map[string]interface{}{"message": msg})
		}).Respond(context.Request)
	} else {
		responder.With("html", func() {
			context.Flash(string(msg), "error")
			http.Redirect(context.Writer, context.Request, context.Request.URL.Path, http.StatusFound)
		}).With("json", func() {
			context.Writer.WriteHeader(422)
			context.JSON("index", map[string]interface{}{"errors": []error{err}})
		})
	}
}
