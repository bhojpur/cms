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

// Package audited is used to log last UpdatedBy and CreatedBy for your models

import "fmt"

// AuditedModel make Model Auditable, embed `audited.AuditedModel` into your model as anonymous field to make the model auditable
//    type User struct {
//      orm.Model
//      audited.AuditedModel
//    }
type AuditedModel struct {
	CreatedBy string
	UpdatedBy string
}

// SetCreatedBy set created by
func (model *AuditedModel) SetCreatedBy(createdBy interface{}) {
	model.CreatedBy = fmt.Sprintf("%v", createdBy)
}

// GetCreatedBy get created by
func (model AuditedModel) GetCreatedBy() string {
	return model.CreatedBy
}

// SetUpdatedBy set updated by
func (model *AuditedModel) SetUpdatedBy(updatedBy interface{}) {
	model.UpdatedBy = fmt.Sprintf("%v", updatedBy)
}

// GetUpdatedBy get updated by
func (model AuditedModel) GetUpdatedBy() string {
	return model.UpdatedBy
}
