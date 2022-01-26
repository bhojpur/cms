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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"
)

type cronJob struct {
	JobID   string
	Pid     int
	Command string
	Delete  bool `json:"-"`
}

func (job cronJob) ToString() string {
	marshal, _ := json.Marshal(job)
	return fmt.Sprintf("## BEGIN BHOJPUR JOB %v # %v\n%v\n## END BHOJPUR JOB\n", job.JobID, string(marshal), job.Command)
}

// Cron implemented a worker Queue based on cronjob
type Cron struct {
	Jobs     []*cronJob
	CronJobs []string
	mutex    sync.Mutex `sql:"-"`
}

// NewCronQueue initialize a Cron queue
func NewCronQueue() *Cron {
	return &Cron{}
}

func (cron *Cron) parseJobs() []*cronJob {
	cron.mutex.Lock()

	cron.Jobs = []*cronJob{}
	cron.CronJobs = []string{}
	if out, err := exec.Command("crontab", "-l").Output(); err == nil {
		var inBhojpurJob bool
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if strings.HasPrefix(line, "## BEGIN BHOJPUR JOB") {
				inBhojpurJob = true
				if idx := strings.Index(line, "{"); idx > 1 {
					var job cronJob
					if json.Unmarshal([]byte(line[idx-1:]), &job) == nil {
						cron.Jobs = append(cron.Jobs, &job)
					}
				}
			}

			if !inBhojpurJob {
				cron.CronJobs = append(cron.CronJobs, line)
			}

			if strings.HasPrefix(line, "## END BHOJPUR JOB") {
				inBhojpurJob = false
			}
		}
	}
	return cron.Jobs
}

func (cron *Cron) writeCronJob() error {
	defer cron.mutex.Unlock()

	cmd := exec.Command("crontab", "-")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdin, _ := cmd.StdinPipe()
	for _, cronJob := range cron.CronJobs {
		stdin.Write([]byte(cronJob + "\n"))
	}

	for _, job := range cron.Jobs {
		if !job.Delete {
			stdin.Write([]byte(job.ToString() + "\n"))
		}
	}
	stdin.Close()
	return cmd.Run()
}

// Add a job to cron queue
func (cron *Cron) Add(job BhojpurJobInterface) (err error) {
	cron.parseJobs()
	defer cron.writeCronJob()

	var binaryFile string
	if binaryFile, err = filepath.Abs(os.Args[0]); err == nil {
		var jobs []*cronJob
		for _, cronJob := range cron.Jobs {
			if cronJob.JobID != job.GetJobID() {
				jobs = append(jobs, cronJob)
			}
		}

		if scheduler, ok := job.GetArgument().(Scheduler); ok && scheduler.GetScheduleTime() != nil {
			scheduleTime := scheduler.GetScheduleTime().In(time.Local)
			job.SetStatus(JobStatusScheduled)

			currentPath, _ := os.Getwd()
			jobs = append(jobs, &cronJob{
				JobID:   job.GetJobID(),
				Command: fmt.Sprintf("%d %d %d %d * cd %v; %v --bhojpur-job %v\n", scheduleTime.Minute(), scheduleTime.Hour(), scheduleTime.Day(), scheduleTime.Month(), currentPath, binaryFile, job.GetJobID()),
			})
		} else {
			cmd := exec.Command(binaryFile, "--bhojpur-job", job.GetJobID())
			if err = cmd.Start(); err == nil {
				jobs = append(jobs, &cronJob{JobID: job.GetJobID(), Pid: cmd.Process.Pid})
				cmd.Process.Release()
			}
		}
		cron.Jobs = jobs
	}

	return
}

// Run a job from cron queue
func (cron *Cron) Run(bhojpurJob BhojpurJobInterface) error {
	job := bhojpurJob.GetJob()

	if job.Handler != nil {
		go func() {
			sigint := make(chan os.Signal, 1)

			// interrupt signal sent from terminal
			signal.Notify(sigint, syscall.SIGINT)
			// sigterm signal sent from kubernetes
			signal.Notify(sigint, syscall.SIGTERM)

			i := <-sigint

			bhojpurJob.SetProgressText(fmt.Sprintf("Worker killed by signal %s", i.String()))
			bhojpurJob.SetStatus(JobStatusKilled)

			bhojpurJob.StopReferesh()
			os.Exit(int(reflect.ValueOf(i).Int()))
		}()

		bhojpurJob.StartReferesh()
		defer bhojpurJob.StopReferesh()

		err := job.Handler(bhojpurJob.GetSerializableArgument(bhojpurJob), bhojpurJob)
		if err == nil {
			cron.parseJobs()
			defer cron.writeCronJob()
			for _, cronJob := range cron.Jobs {
				if cronJob.JobID == bhojpurJob.GetJobID() {
					cronJob.Delete = true
				}
			}
		}
		return err
	}

	return errors.New("no handler found for job " + job.Name)
}

// Kill a job from cron queue
func (cron *Cron) Kill(job BhojpurJobInterface) (err error) {
	cron.parseJobs()
	defer cron.writeCronJob()

	for _, cronJob := range cron.Jobs {
		if cronJob.JobID == job.GetJobID() {
			if process, err := os.FindProcess(cronJob.Pid); err == nil {
				if err = process.Kill(); err == nil {
					cronJob.Delete = true
					return nil
				}
			}
			return err
		}
	}
	return errors.New("failed to find job")
}

// Remove a job from cron queue
func (cron *Cron) Remove(job BhojpurJobInterface) error {
	cron.parseJobs()
	defer cron.writeCronJob()

	for _, cronJob := range cron.Jobs {
		if cronJob.JobID == job.GetJobID() {
			if cronJob.Pid == 0 {
				cronJob.Delete = true
				return nil
			}
			return errors.New("failed to remove current job as it is running")
		}
	}
	return errors.New("failed to find job")
}
