package audited

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

	orm "github.com/bhojpur/orm/pkg/engine"
)

type auditableInterface interface {
	SetCreatedBy(createdBy interface{})
	GetCreatedBy() string
	SetUpdatedBy(updatedBy interface{})
	GetUpdatedBy() string
}

func isAuditable(scope *orm.Scope) (isAuditable bool) {
	if scope.GetModelStruct().ModelType == nil {
		return false
	}
	_, isAuditable = reflect.New(scope.GetModelStruct().ModelType).Interface().(auditableInterface)
	return
}

func getCurrentUser(scope *orm.Scope) (string, bool) {
	var user interface{}
	var hasUser bool

	user, hasUser = scope.DB().Get("audited:current_user")

	if !hasUser {
		user, hasUser = scope.DB().Get("bhojpur:current_user")
	}

	if hasUser {
		var currentUser string
		if primaryField := scope.New(user).PrimaryField(); primaryField != nil {
			currentUser = fmt.Sprintf("%v", primaryField.Field.Interface())
		} else {
			currentUser = fmt.Sprintf("%v", user)
		}

		return currentUser, true
	}

	return "", false
}

func assignCreatedBy(scope *orm.Scope) {
	if isAuditable(scope) {
		if user, ok := getCurrentUser(scope); ok {
			scope.SetColumn("CreatedBy", user)
		}
	}
}

func assignUpdatedBy(scope *orm.Scope) {
	if isAuditable(scope) {
		if user, ok := getCurrentUser(scope); ok {
			if attrs, ok := scope.InstanceGet("orm:update_attrs"); ok {
				updateAttrs := attrs.(map[string]interface{})
				updateAttrs["updated_by"] = user
				scope.InstanceSet("orm:update_attrs", updateAttrs)
			} else {
				scope.SetColumn("UpdatedBy", user)
			}
		}
	}
}

// RegisterCallbacks register callback into Bhojpur ORM DB
func RegisterCallbacks(db *orm.DB) {
	callback := db.Callback()
	if callback.Create().Get("audited:assign_created_by") == nil {
		callback.Create().After("orm:before_create").Register("audited:assign_created_by", assignCreatedBy)
	}
	if callback.Update().Get("audited:assign_updated_by") == nil {
		callback.Update().After("orm:before_update").Register("audited:assign_updated_by", assignUpdatedBy)
	}
}
