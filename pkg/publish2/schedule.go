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
	"time"

	"github.com/bhojpur/application/pkg/validations"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type Schedule struct {
	ScheduledStartAt *time.Time `orm:"index"`
	ScheduledEndAt   *time.Time `orm:"index"`
	ScheduledEventID *uint
}

func (schedule *Schedule) GetScheduledStartAt() *time.Time {
	return schedule.ScheduledStartAt
}

func (schedule *Schedule) SetScheduledStartAt(t *time.Time) {
	schedule.ScheduledStartAt = t
}

func (schedule *Schedule) GetScheduledEndAt() *time.Time {
	return schedule.ScheduledEndAt
}

func (schedule *Schedule) SetScheduledEndAt(t *time.Time) {
	schedule.ScheduledEndAt = t
}

func (schedule *Schedule) GetScheduledEventID() *uint {
	return schedule.ScheduledEventID
}

type ScheduledInterface interface {
	GetScheduledStartAt() *time.Time
	SetScheduledStartAt(*time.Time)
	GetScheduledEndAt() *time.Time
	SetScheduledEndAt(*time.Time)
	GetScheduledEventID() *uint
}

type ScheduledEvent struct {
	orm.Model
	Name             string
	ScheduledStartAt *time.Time
	ScheduledEndAt   *time.Time
}

func (scheduledEvent ScheduledEvent) ToParam() string {
	return "scheduled_events"
}

func (scheduledEvent ScheduledEvent) BeforeSave(tx *orm.DB) {
	if scheduledEvent.Name == "" {
		tx.AddError(validations.NewError(scheduledEvent, "Name", "Name can not be empty"))
	}
}
