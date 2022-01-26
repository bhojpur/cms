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
	"strings"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/worker"
)

type workerJobLogger struct {
	job worker.BhojpurJobInterface
}

func (job workerJobLogger) Print(results ...interface{}) {
	job.job.AddLog(fmt.Sprint(results...))
}

// BhojpurWorkerArgument used for publish job's argument
type BhojpurWorkerArgument struct {
	IDs []string
	worker.Schedule
}

// SetWorker set publish's worker
func (publish *Publish) SetWorker(w *worker.Worker) {
	publish.WorkerScheduler = w
	publish.registerWorkerJob()
}

func (publish *Publish) registerWorkerJob() {
	if w := publish.WorkerScheduler; w != nil {
		if w.Admin == nil {
			fmt.Println("Need to add worker to admin first before set worker")
			return
		}

		bhojpurWorkerArgumentResource := w.Admin.NewResource(&BhojpurWorkerArgument{})
		bhojpurWorkerArgumentResource.Meta(&admin.Meta{Name: "IDs", Type: "publish_job_argument", Valuer: func(record interface{}, context *appsvr.Context) interface{} {
			var values = map[*admin.Resource][][]string{}

			if workerArgument, ok := record.(*BhojpurWorkerArgument); ok {
				for _, id := range workerArgument.IDs {
					if keys := strings.Split(id, "__"); len(keys) >= 2 {
						name, id := keys[0], keys[1:]
						recordRes := w.Admin.GetResource(name)
						values[recordRes] = append(values[recordRes], id)
					}
				}
			}

			return values
		}})

		w.RegisterJob(&worker.Job{
			Name:       "Publish",
			Group:      "Publish",
			Permission: roles.Deny(roles.Read, roles.Anyone),
			Handler: func(argument interface{}, job worker.BhojpurJobInterface) error {
				if argu, ok := argument.(*BhojpurWorkerArgument); ok {
					records := publish.searchWithPublishIDs(publish.DraftDB(), w.Admin, argu.IDs)
					publish.Logger(&workerJobLogger{job: job}).Publish(records...)
				}
				return nil
			},
			Resource: bhojpurWorkerArgumentResource,
		})

		w.RegisterJob(&worker.Job{
			Name:       "Discard",
			Group:      "Publish",
			Permission: roles.Deny(roles.Read, roles.Anyone),
			Handler: func(argument interface{}, job worker.BhojpurJobInterface) error {
				if argu, ok := argument.(*BhojpurWorkerArgument); ok {
					records := publish.searchWithPublishIDs(publish.DraftDB(), w.Admin, argu.IDs)
					publish.Logger(&workerJobLogger{job: job}).Discard(records...)
				}
				return nil
			},
			Resource: bhojpurWorkerArgumentResource,
		})
	}
}
