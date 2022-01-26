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
	"errors"
	"fmt"

	appsvr "github.com/bhojpur/application/pkg/engine"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// EventInterface defined methods needs for a publish event
type EventInterface interface {
	Publish(db *orm.DB, event PublishEventInterface) error
	Discard(db *orm.DB, event PublishEventInterface) error
}

var events = map[string]EventInterface{}

// RegisterEvent register publish event
func RegisterEvent(name string, event EventInterface) {
	events[name] = event
}

// PublishEvent default publish event model
type PublishEvent struct {
	orm.Model
	Name          string
	Description   string
	Argument      string `sql:"size:65532"`
	PublishStatus bool
	PublishedBy   string
}

func getCurrentUser(db *orm.DB) (string, bool) {
	if user, hasUser := db.Get("bhojpur:current_user"); hasUser {
		var currentUser string
		if primaryField := db.NewScope(user).PrimaryField(); primaryField != nil {
			currentUser = fmt.Sprintf("%v", primaryField.Field.Interface())
		} else {
			currentUser = fmt.Sprintf("%v", user)
		}

		return currentUser, true
	}

	return "", false
}

// Publish publish data
func (publishEvent *PublishEvent) Publish(db *orm.DB) error {
	if event, ok := events[publishEvent.Name]; ok {
		err := event.Publish(db, publishEvent)
		if err == nil {
			var updateAttrs = map[string]interface{}{"PublishStatus": PUBLISHED}
			if user, hasUser := getCurrentUser(db); hasUser {
				updateAttrs["PublishedBy"] = user
			}
			err = db.Model(publishEvent).Update(updateAttrs).Error
		}
		return err
	}
	return errors.New("event not found")
}

// Discard discard data
func (publishEvent *PublishEvent) Discard(db *orm.DB) error {
	if event, ok := events[publishEvent.Name]; ok {
		err := event.Discard(db, publishEvent)
		if err == nil {
			var updateAttrs = map[string]interface{}{"PublishStatus": PUBLISHED}
			if user, hasUser := getCurrentUser(db); hasUser {
				updateAttrs["PublishedBy"] = user
			}
			err = db.Model(publishEvent).Update(updateAttrs).Error
		}
		return err
	}
	return errors.New("event not found")
}

// VisiblePublishResource force to display publish event in publish drafts even it is hidden in the menus
func (PublishEvent) VisiblePublishResource(*appsvr.Context) bool {
	return true
}
