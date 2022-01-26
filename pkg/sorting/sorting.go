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
	"fmt"
	"reflect"

	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/publish"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type sortingInterface interface {
	GetPosition() int
	SetPosition(int)
}

type sortingDescInterface interface {
	GetPosition() int
	SetPosition(int)
	SortingDesc()
}

// Sorting ascending mode
type Sorting struct {
	Position int `sql:"DEFAULT:NULL"`
}

// GetPosition get current position
func (position Sorting) GetPosition() int {
	return position.Position
}

// SetPosition set position, only set field value, won't save
func (position *Sorting) SetPosition(pos int) {
	position.Position = pos
}

// SortingDESC descending mode
type SortingDESC struct {
	Sorting
}

// SortingDesc make your model sorting desc by default
func (SortingDESC) SortingDesc() {}

func init() {
	admin.RegisterViewPath("github.com/bhojpur/cms/pkg/sorting/views")
}

func newModel(value interface{}) interface{} {
	return reflect.New(reflect.Indirect(reflect.ValueOf(value)).Type()).Interface()
}

func move(db *orm.DB, value sortingInterface, pos int) (err error) {
	var startedTransaction bool
	var tx = db.Set("publish:publish_event", true)
	if t := tx.Begin(); t.Error == nil {
		startedTransaction = true
		tx = t
	}

	scope := db.NewScope(value)
	for _, field := range scope.PrimaryFields() {
		// "version_name" is a "reserved" primary key, we always update all versions postion at the same time.
		// so don't count version name as a condition.
		if field.DBName != "id" && field.DBName != "version_name" {
			tx = tx.Where(fmt.Sprintf("%s = ?", field.DBName), field.Field.Interface())
		}
	}

	currentPos := value.GetPosition()
	var results *orm.DB
	if pos > 0 {
		results = tx.Table(scope.TableName()).Where("position > ? AND position <= ?", currentPos, currentPos+pos).
			UpdateColumn("position", orm.Expr("position - ?", 1))
	} else {
		results = tx.Table(scope.TableName()).Where("position < ? AND position >= ?", currentPos, currentPos+pos).
			UpdateColumn("position", orm.Expr("position + ?", 1))
	}

	if err = results.Error; err == nil && results.RowsAffected > 0 {
		// Use ID as the ONLY condition, so that we can update all version of one record's position.
		modelObj := reflect.Indirect(reflect.ValueOf(value))
		err = tx.Table(scope.TableName()).Where("id = ?", modelObj.FieldByName("ID").Interface().(uint)).UpdateColumn("position", currentPos+pos).Error
	}

	// Create Publish Event
	createPublishEvent(tx, value)

	if startedTransaction {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}
	return err
}

func createPublishEvent(db *orm.DB, value interface{}) (err error) {
	// Create Publish Event in Draft Mode
	if publish.IsDraftMode(db) && publish.IsPublishableModel(value) {
		scope := db.NewScope(value)
		var sortingPublishEvent = changedSortingPublishEvent{
			Table: scope.TableName(),
		}
		for _, field := range scope.PrimaryFields() {
			sortingPublishEvent.PrimaryKeys = append(sortingPublishEvent.PrimaryKeys, field.DBName)
		}

		var result []byte
		if result, err = json.Marshal(sortingPublishEvent); err == nil {
			err = db.New().Where("publish_status = ?", publish.DIRTY).Where(map[string]interface{}{
				"name":     "changed_sorting",
				"argument": string(result),
			}).Attrs(map[string]interface{}{
				"publish_status": publish.DIRTY,
				"description":    "Changed sort order for " + scope.GetModelStruct().ModelType.Name(),
			}).FirstOrCreate(&publish.PublishEvent{}).Error
		}
	}
	return
}

// MoveUp move position up
func MoveUp(db *orm.DB, value sortingInterface, pos int) error {
	return move(db, value, -pos)
}

// MoveDown move position down
func MoveDown(db *orm.DB, value sortingInterface, pos int) error {
	return move(db, value, pos)
}

// MoveTo move position to
func MoveTo(db *orm.DB, value sortingInterface, pos int) error {
	return move(db, value, pos-value.GetPosition())
}
