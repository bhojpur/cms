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
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/cms/pkg/admin"
	orm "github.com/bhojpur/orm/pkg/engine"
)

const (
	// JobStatusScheduled job status scheduled
	JobStatusScheduled = "scheduled"
	// JobStatusCancelled job status cancelled
	JobStatusCancelled = "cancelled"
	// JobStatusNew job status new
	JobStatusNew = "new"
	// JobStatusRunning job status running
	JobStatusRunning = "running"
	// JobStatusDone job status done
	JobStatusDone = "done"
	// JobStatusException job status exception
	JobStatusException = "exception"
	// JobStatusKilled job status killed
	JobStatusKilled = "killed"
)

// New create Worker with Config
func New(config ...*Config) *Worker {
	var cfg = &Config{}
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.Job == nil {
		cfg.Job = &BhojpurJob{}
	}

	if cfg.Queue == nil {
		cfg.Queue = NewCronQueue()
	}

	return &Worker{Config: cfg}
}

// Config worker config
type Config struct {
	Queue Queue
	Job   BhojpurJobInterface
	Admin *admin.Admin
}

// Worker worker definition
type Worker struct {
	*Config
	JobResource *admin.Resource
	Jobs        []*Job
	mounted     bool
}

// ConfigureBhojpurResourceBeforeInitialize a method used to config Worker for Bhojpur admin
func (worker *Worker) ConfigureBhojpurResourceBeforeInitialize(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		res.GetAdmin().RegisterViewPath("github.com/bhojpur/cms/pkg/worker/views")
		res.UseTheme("worker")

		worker.Admin = res.GetAdmin()
		worker.JobResource = worker.Admin.NewResource(worker.Config.Job)
		worker.JobResource.UseTheme("worker")
		worker.JobResource.Meta(&admin.Meta{Name: "Name", Valuer: func(record interface{}, context *appsvr.Context) interface{} {
			return record.(BhojpurJobInterface).GetJobName()
		}})
		worker.JobResource.IndexAttrs("ID", "Name", "Status", "CreatedAt")
		worker.JobResource.Name = res.Name

		for _, status := range []string{JobStatusScheduled, JobStatusNew, JobStatusRunning, JobStatusDone, JobStatusException} {
			var status = status
			worker.JobResource.Scope(&admin.Scope{Name: status, Handler: func(db *orm.DB, ctx *appsvr.Context) *orm.DB {
				return db.Where("status = ?", status)
			}})
		}

		// default scope
		worker.JobResource.Scope(&admin.Scope{
			Handler: func(db *orm.DB, ctx *appsvr.Context) *orm.DB {
				if jobName := ctx.Request.URL.Query().Get("job"); jobName != "" {
					return db.Where("kind = ?", jobName)
				}

				if groupName := ctx.Request.URL.Query().Get("group"); groupName != "" {
					var jobNames []string
					for _, job := range worker.Jobs {
						if groupName == job.Group {
							jobNames = append(jobNames, job.Name)
						}
					}
					if len(jobNames) > 0 {
						return db.Where("kind IN (?)", jobNames)
					}
					return db.Where("kind IS NULL")
				}

				{
					var jobNames []string
					for _, job := range worker.Jobs {
						jobNames = append(jobNames, job.Name)
					}
					if len(jobNames) > 0 {
						return db.Where("kind IN (?)", jobNames)
					}
				}

				return db
			},
			Default: true,
		})

		// Auto Migration
		worker.Admin.DB.AutoMigrate(worker.Config.Job)

		// Configure jobs
		for _, job := range worker.Jobs {
			if job.Resource == nil {
				job.Resource = worker.Admin.NewResource(worker.JobResource.Value)
			}
		}
	}
}

// ConfigureBhojpurResource a method used to config Worker for Bhojpur admin
func (worker *Worker) ConfigureBhojpurResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		// Parse job
		cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		bhojpurJobID := cmdLine.String("bhojpur-job", "", "Bhojpur Job ID")
		runAnother := cmdLine.Bool("run-another", false, "Run another Bhojpur job")
		cmdLine.Parse(os.Args[1:])
		worker.mounted = true

		if *bhojpurJobID != "" {
			if *runAnother == true {
				if newJob := worker.saveAnotherJob(*bhojpurJobID); newJob != nil {
					newJobID := newJob.GetJobID()
					bhojpurJobID = &newJobID
				} else {
					fmt.Println("failed to clone job " + *bhojpurJobID)
					os.Exit(1)
				}
			}

			if err := worker.RunJob(*bhojpurJobID); err == nil {
				os.Exit(0)
			} else {
				fmt.Println(err)
				// os.Exit(1)
			}
		}

		// register view funcmaps
		worker.Admin.RegisterFuncMap("get_grouped_jobs", func(worker *Worker, context *admin.Context) map[string][]*Job {
			var groupedJobs = map[string][]*Job{}
			var groupName = context.Request.URL.Query().Get("group")
			var jobName = context.Request.URL.Query().Get("job")
			for _, job := range worker.Jobs {
				if !(job.HasPermission(roles.Read, context.Context) && job.HasPermission(roles.Create, context.Context)) {
					continue
				}

				if (groupName == "" || groupName == job.Group) && (jobName == "" || jobName == job.Name) {
					groupedJobs[job.Group] = append(groupedJobs[job.Group], job)
				}
			}
			return groupedJobs
		})

		// configure routes
		router := worker.Admin.GetRouter()
		controller := workerController{Worker: worker}
		jobParamIDName := worker.JobResource.ParamIDName()

		router.Get(res.ToParam(), controller.Index, &admin.RouteConfig{Resource: worker.JobResource})
		router.Get(res.ToParam()+"/new", controller.New, &admin.RouteConfig{Resource: worker.JobResource})
		router.Get(fmt.Sprintf("%v/%v", res.ToParam(), jobParamIDName), controller.Show, &admin.RouteConfig{Resource: worker.JobResource})
		router.Get(fmt.Sprintf("%v/%v/edit", res.ToParam(), jobParamIDName), controller.Show, &admin.RouteConfig{Resource: worker.JobResource})
		router.Post(fmt.Sprintf("%v/%v/run", res.ToParam(), jobParamIDName), controller.RunJob, &admin.RouteConfig{Resource: worker.JobResource})
		router.Post(res.ToParam(), controller.AddJob, &admin.RouteConfig{Resource: worker.JobResource})
		router.Put(fmt.Sprintf("%v/%v", res.ToParam(), jobParamIDName), controller.Update, &admin.RouteConfig{Resource: worker.JobResource})
		router.Delete(fmt.Sprintf("%v/%v", res.ToParam(), jobParamIDName), controller.KillJob, &admin.RouteConfig{Resource: worker.JobResource})
	}
}

// SetQueue set worker's queue
func (worker *Worker) SetQueue(queue Queue) {
	worker.Queue = queue
}

// RegisterJob register a job into Worker
func (worker *Worker) RegisterJob(job *Job) error {
	if worker.mounted {
		debug.PrintStack()
		fmt.Printf("Job should be registered before Worker mounted into admin, but %v is registered after that", job.Name)
	}

	job.Worker = worker
	worker.Jobs = append(worker.Jobs, job)
	return nil
}

// GetRegisteredJob register a job into Worker
func (worker *Worker) GetRegisteredJob(name string) *Job {
	for _, job := range worker.Jobs {
		if job.Name == name {
			return job
		}
	}
	return nil
}

// GetJob get job with id
func (worker *Worker) GetJob(jobID string) (BhojpurJobInterface, error) {
	bhojpurJob := worker.JobResource.NewStruct().(BhojpurJobInterface)

	context := worker.Admin.NewContext(nil, nil)
	context.ResourceID = jobID
	context.Resource = worker.JobResource

	if err := worker.JobResource.FindOneHandler(bhojpurJob, nil, context.Context); err == nil {
		for _, job := range worker.Jobs {
			if job.Name == bhojpurJob.GetJobName() {
				bhojpurJob.SetJob(job)
				return bhojpurJob, nil
			}
		}
		return nil, fmt.Errorf("failed to load job: %v, unregistered job type: %v", jobID, bhojpurJob.GetJobName())
	}
	return nil, fmt.Errorf("failed to find job: %v", jobID)
}

// AddJob add job to worker
func (worker *Worker) AddJob(bhojpurJob BhojpurJobInterface) error {
	return worker.Queue.Add(bhojpurJob)
}

// RunJob run job with job id
func (worker *Worker) RunJob(jobID string) error {
	bhojpurJob, err := worker.GetJob(jobID)

	if bhojpurJob != nil && err == nil {
		defer func() {
			if r := recover(); r != nil {
				bhojpurJob.AddLog(string(debug.Stack()))
				bhojpurJob.SetProgressText(fmt.Sprint(r))
				bhojpurJob.SetStatus(JobStatusException)
			}
		}()

		if bhojpurJob.GetStatus() != JobStatusNew && bhojpurJob.GetStatus() != JobStatusScheduled {
			return errors.New("invalid job status, current status: " + bhojpurJob.GetStatus())
		}

		if err = bhojpurJob.SetStatus(JobStatusRunning); err == nil {
			if err = bhojpurJob.GetJob().GetQueue().Run(bhojpurJob); err == nil {
				if len(bhojpurJob.GetResultsTable().TableCells) > 0 {
					return bhojpurJob.SetStatus(JobStatusException)
				}
				return bhojpurJob.SetStatus(JobStatusDone)
			}

			bhojpurJob.SetProgressText(err.Error())
			bhojpurJob.SetStatus(JobStatusException)
		}
	}

	return err
}

func (worker *Worker) saveAnotherJob(jobID string) BhojpurJobInterface {
	jobResource := worker.JobResource
	newJob := jobResource.NewStruct().(BhojpurJobInterface)

	job, err := worker.GetJob(jobID)
	if err == nil {
		newJob.SetJob(job.GetJob())
		newJob.SetSerializableArgumentValue(job.GetArgument())
		context := worker.Admin.NewContext(nil, nil)
		if err := jobResource.CallSave(newJob, context.Context); err == nil {
			return newJob
		}
	}
	return nil
}

// KillJob kill job with job id
func (worker *Worker) KillJob(jobID string) error {
	if bhojpurJob, err := worker.GetJob(jobID); err == nil {
		if bhojpurJob.GetStatus() == JobStatusRunning {
			if err = bhojpurJob.GetJob().GetQueue().Kill(bhojpurJob); err == nil {
				bhojpurJob.SetStatus(JobStatusKilled)
				return nil
			}
			return err
		} else if bhojpurJob.GetStatus() == JobStatusScheduled || bhojpurJob.GetStatus() == JobStatusNew {
			bhojpurJob.SetStatus(JobStatusKilled)
			return worker.RemoveJob(jobID)
		} else {
			return errors.New("invalid job status")
		}
	} else {
		return err
	}
}

// RemoveJob remove job with job id
func (worker *Worker) RemoveJob(jobID string) error {
	bhojpurJob, err := worker.GetJob(jobID)
	if err == nil {
		return bhojpurJob.GetJob().GetQueue().Remove(bhojpurJob)
	}
	return err
}
