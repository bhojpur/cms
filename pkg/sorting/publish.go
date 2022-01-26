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
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bhojpur/cms/pkg/publish"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type changedSortingPublishEvent struct {
	Table       string
	PrimaryKeys []string
}

func (e changedSortingPublishEvent) Publish(db *orm.DB, event publish.PublishEventInterface) (err error) {
	if event, ok := event.(*publish.PublishEvent); ok {
		scope := db.NewScope("")
		if err = json.Unmarshal([]byte(event.Argument), &e); err == nil {
			var conditions []string
			originalTable := scope.Quote(publish.OriginalTableName(e.Table))
			draftTable := scope.Quote(publish.DraftTableName(e.Table))
			for _, primaryKey := range e.PrimaryKeys {
				conditions = append(conditions, fmt.Sprintf("%v.%v = %v.%v", originalTable, primaryKey, draftTable, primaryKey))
			}
			sql := fmt.Sprintf("UPDATE %v SET position = (select position FROM %v WHERE %v);", originalTable, draftTable, strings.Join(conditions, " AND "))
			return db.Exec(sql).Error
		}
		return err
	}
	return errors.New("invalid publish event")
}

func (e changedSortingPublishEvent) Discard(db *orm.DB, event publish.PublishEventInterface) (err error) {
	if event, ok := event.(*publish.PublishEvent); ok {
		scope := db.NewScope("")
		if err = json.Unmarshal([]byte(event.Argument), &e); err == nil {
			var conditions []string
			originalTable := scope.Quote(publish.OriginalTableName(e.Table))
			draftTable := scope.Quote(publish.DraftTableName(e.Table))
			for _, primaryKey := range e.PrimaryKeys {
				conditions = append(conditions, fmt.Sprintf("%v.%v = %v.%v", originalTable, primaryKey, draftTable, primaryKey))
			}
			sql := fmt.Sprintf("UPDATE %v SET position = (select position FROM %v WHERE %v);", draftTable, originalTable, strings.Join(conditions, " AND "))
			return db.Exec(sql).Error
		}
		return err
	}
	return errors.New("invalid publish event")
}

func init() {
	publish.RegisterEvent("changed_sorting", changedSortingPublishEvent{})
}
