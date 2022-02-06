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
	"fmt"
	"reflect"
	"strings"

	"github.com/bhojpur/cms/pkg/l10n"
	orm "github.com/bhojpur/orm/pkg/engine"
)

func initalizePosition(scope *orm.Scope) {
	if !scope.HasError() {
		if _, ok := scope.Value.(sortingInterface); ok {
			var lastPosition int
			scope.NewDB().Set("l10n:mode", "locale").Model(modelValue(scope.Value)).Select("position").Order("position DESC").Limit(1).Row().Scan(&lastPosition)
			scope.SetColumn("Position", lastPosition+1)
		}
	}
}

func reorderPositions(scope *orm.Scope) {
	if !scope.HasError() {
		if _, ok := scope.Value.(sortingInterface); ok {
			table := scope.TableName()
			var additionalSQL []string
			var additionalValues []interface{}
			// with l10n
			if locale, ok := scope.DB().Get("l10n:locale"); ok && locale.(string) != "" && l10n.IsLocalizable(scope) {
				additionalSQL = append(additionalSQL, "language_code = ?")
				additionalValues = append(additionalValues, locale)
			}
			additionalValues = append(additionalValues, additionalValues...)

			// with soft delete
			if scope.HasColumn("DeletedAt") {
				additionalSQL = append(additionalSQL, "deleted_at IS NULL")
			}

			var sql string
			if len(additionalSQL) > 0 {
				sql = fmt.Sprintf("UPDATE %v SET position = (SELECT COUNT(pos) + 1 FROM (SELECT DISTINCT(position) AS pos FROM %v WHERE %v) AS t2 WHERE t2.pos < %v.position) WHERE %v", table, table, strings.Join(additionalSQL, " AND "), table, strings.Join(additionalSQL, " AND "))
			} else {
				sql = fmt.Sprintf("UPDATE %v SET position = (SELECT COUNT(pos) + 1 FROM (SELECT DISTINCT(position) AS pos FROM %v) AS t2 WHERE t2.pos < %v.position)", table, table, table)
			}
			if scope.NewDB().Exec(sql, additionalValues...).Error == nil {
				// Create Publish Event
				createPublishEvent(scope.DB(), scope.Value)
			}
		}
	}
}

func modelValue(value interface{}) interface{} {
	reflectValue := reflect.Indirect(reflect.ValueOf(value))
	if reflectValue.IsValid() {
		typ := reflectValue.Type()

		if reflectValue.Kind() == reflect.Slice {
			typ = reflectValue.Type().Elem()
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
		}

		return reflect.New(typ).Interface()
	}
	return nil
}

func beforeQuery(scope *orm.Scope) {
	modelValue := modelValue(scope.Value)
	if _, ok := modelValue.(sortingDescInterface); ok {
		scope.Search.Order("position desc")
	} else if _, ok := modelValue.(sortingInterface); ok {
		scope.Search.Order("position")
	}
}

// RegisterCallbacks register callbacks into Bhojpur ORM DB instance
func RegisterCallbacks(db *orm.DB) {
	if db.Callback().Create().Get("sorting:initalize_position") == nil {
		db.Callback().Create().Before("orm:create").Register("sorting:initalize_position", initalizePosition)
	}
	if db.Callback().Delete().Get("sorting:reorder_positions") == nil {
		db.Callback().Delete().After("orm:after_delete").Register("sorting:reorder_positions", reorderPositions)
	}
	if db.Callback().Query().Get("sorting:sort_by_position") == nil {
		db.Callback().Query().Before("orm:query").Register("sorting:sort_by_position", beforeQuery)
	}
}
