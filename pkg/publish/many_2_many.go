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
	orm "github.com/bhojpur/orm/pkg/engine"
)

type publishJoinTableHandler struct {
	orm.JoinTableHandler
}

func (handler publishJoinTableHandler) Table(db *orm.DB) string {
	if IsDraftMode(db) {
		return handler.TableName + "_draft"
	}
	return handler.TableName
}

func (handler publishJoinTableHandler) Add(h orm.JoinTableHandlerInterface, db *orm.DB, source1 interface{}, source2 interface{}) error {
	// production mode
	if !IsDraftMode(db) {
		if err := handler.JoinTableHandler.Add(h, db.Set(publishDraftMode, true), source1, source2); err != nil {
			return err
		}
	}
	return handler.JoinTableHandler.Add(h, db, source1, source2)
}

func (handler publishJoinTableHandler) Delete(h orm.JoinTableHandlerInterface, db *orm.DB, sources ...interface{}) error {
	// production mode
	if !IsDraftMode(db) {
		if err := handler.JoinTableHandler.Delete(h, db.Set(publishDraftMode, true), sources...); err != nil {
			return err
		}
	}
	return handler.JoinTableHandler.Delete(h, db, sources...)
}
