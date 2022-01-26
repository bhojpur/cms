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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/audited"
	"github.com/bhojpur/cms/pkg/serializable_meta"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// BhojpurJobInterface is a interface, defined methods that needs for a Bhojpur job
type BhojpurJobInterface interface {
	GetJobID() string
	GetJobName() string
	GetStatus() string
	SetStatus(string) error
	GetJob() *Job
	SetJob(*Job)

	GetProgress() uint
	SetProgress(uint) error
	GetProgressText() string
	SetProgressText(string) error
	GetLogs() []string
	AddLog(string) error
	GetResultsTable() ResultsTable
	AddResultsRow(...TableCell) error

	StartReferesh()
	StopReferesh()

	GetArgument() interface{}
	serializable_meta.SerializableMetaInterface
}

// ResultsTable is a struct, including importing/exporting results
type ResultsTable struct {
	Name       string `json:"-"` // only used for generate string column in database
	TableCells [][]TableCell
}

// Scan used to scan value from database into itself
func (resultsTable *ResultsTable) Scan(data interface{}) error {
	switch values := data.(type) {
	case []byte:
		return json.Unmarshal(values, resultsTable)
	case string:
		return resultsTable.Scan([]byte(values))
	default:
		return errors.New("unsupported data type for Bhojpur Job error table")
	}
}

// Value used to read value from itself and save it into databae
func (resultsTable ResultsTable) Value() (driver.Value, error) {
	result, err := json.Marshal(resultsTable)
	return string(result), err
}

// TableCell including Value, Error for a data cell
type TableCell struct {
	Value string
	Error string
}

// BhojpurJob predefined Bhojpur job struct, which will be used for Worker, if it doesn't include a job resource
type BhojpurJob struct {
	orm.Model
	Status       string `sql:"default:'new'"`
	Progress     uint
	ProgressText string
	Log          string       `sql:"size:65532"`
	ResultsTable ResultsTable `sql:"size:65532"`

	mutex sync.Mutex `sql:"-"`

	stopReferesh bool `sql:"-"`
	inReferesh   bool `sql:"-"`

	// Add `valid:"-"`` to make the BhojpurJob work well with bhojpur/errors/pkg/validation
	// When the bhojpur/errors/pkg/validation auto exec the validate struct callback we get error
	// runtime: goroutine stack exceeds 1000000000-byte limit
	// fatal error: stack overflow
	Job *Job `sql:"-" valid:"-"`

	audited.AuditedModel
	serializable_meta.SerializableMeta
}

// GetJobID get job's ID from a job
func (job *BhojpurJob) GetJobID() string {
	return fmt.Sprint(job.ID)
}

// GetJobName get job's name from a Bhojpur job
func (job *BhojpurJob) GetJobName() string {
	return job.Kind
}

// GetStatus get job's status from a Bhojpur job
func (job *BhojpurJob) GetStatus() string {
	return job.Status
}

// SetStatus set job's status to a Bhojpur job instance
func (job *BhojpurJob) SetStatus(status string) error {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	job.Status = status
	if status == JobStatusDone {
		job.Progress = 100
	}

	if job.shouldCallSave() {
		return job.callSave()
	}

	return nil
}

func (job *BhojpurJob) shouldCallSave() bool {
	return !job.inReferesh || job.stopReferesh
}

func (job *BhojpurJob) StartReferesh() {
	job.mutex.Lock()
	defer job.mutex.Unlock()
	if !job.inReferesh {
		job.inReferesh = true
		job.stopReferesh = false

		go func() {
			job.referesh()
		}()
	}
}

func (job *BhojpurJob) StopReferesh() {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	err := job.callSave()
	if err != nil {
		log.Println(err)
	}

	job.stopReferesh = true
}

func (job *BhojpurJob) referesh() {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	err := job.callSave()
	if err != nil {
		log.Println(err)
	}

	if job.stopReferesh {
		job.inReferesh = false
		job.stopReferesh = false
	} else {
		time.AfterFunc(5*time.Second, job.referesh)
	}
}

func (job *BhojpurJob) callSave() error {
	worker := job.GetJob().Worker
	context := worker.Admin.NewContext(nil, nil).Context
	return worker.JobResource.CallSave(job, context)
}

// SetJob set `Job` for a Bhojpur job instance
func (job *BhojpurJob) SetJob(j *Job) {
	job.Kind = j.Name
	job.Job = j
}

// GetJob get predefined job for a Bhojpur job instance
func (job *BhojpurJob) GetJob() *Job {
	if job.Job != nil {
		return job.Job
	}
	return nil
}

// GetArgument get job's argument
func (job *BhojpurJob) GetArgument() interface{} {
	return job.GetSerializableArgument(job)
}

// GetSerializableArgumentResource get job's argument's resource
func (job *BhojpurJob) GetSerializableArgumentResource() *admin.Resource {
	if j := job.GetJob(); j != nil {
		return j.Resource
	}
	return nil
}

// GetProgress get Bhojpur job's progress
func (job *BhojpurJob) GetProgress() uint {
	return job.Progress
}

// SetProgress set Bhojpur job's progress
func (job *BhojpurJob) SetProgress(progress uint) error {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	if progress > 100 {
		progress = 100
	}
	job.Progress = progress

	if job.shouldCallSave() {
		return job.callSave()
	}

	return nil
}

// GetProgressText get Bhojpur job's progress text
func (job *BhojpurJob) GetProgressText() string {
	return job.ProgressText
}

// SetProgressText set Bhojpur job's progress text
func (job *BhojpurJob) SetProgressText(str string) error {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	job.ProgressText = str
	if job.shouldCallSave() {
		return job.callSave()
	}

	return nil
}

// GetLogs get Bhojpur job's logs
func (job *BhojpurJob) GetLogs() []string {
	return strings.Split(job.Log, "\n")
}

// AddLog add a log to Bhojpur job
func (job *BhojpurJob) AddLog(log string) error {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	fmt.Println(log)
	job.Log += "\n" + log
	if job.shouldCallSave() {
		return job.callSave()
	}

	return nil
}

// GetResultsTable get the job's process logs
func (job *BhojpurJob) GetResultsTable() ResultsTable {
	return job.ResultsTable
}

// AddResultsRow add a row of process results to a job
func (job *BhojpurJob) AddResultsRow(cells ...TableCell) error {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	job.ResultsTable.TableCells = append(job.ResultsTable.TableCells, cells)
	if job.shouldCallSave() {
		return job.callSave()
	}

	return nil
}
