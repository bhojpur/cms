package transition

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
	"strings"

	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/audited"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// StateChangeLog a model that used to keep state change logs
type StateChangeLog struct {
	orm.Model
	ReferTable string
	ReferID    string
	From       string
	To         string
	Note       string `sql:"size:1024"`
	audited.AuditedModel
}

// GenerateReferenceKey generate reference key used for change log
func GenerateReferenceKey(model interface{}, db *orm.DB) string {
	var (
		scope         = db.NewScope(model)
		primaryValues []string
	)

	for _, field := range scope.PrimaryFields() {
		primaryValues = append(primaryValues, fmt.Sprint(field.Field.Interface()))
	}

	return strings.Join(primaryValues, "::")
}

// GetStateChangeLogs get state change logs
func GetStateChangeLogs(model interface{}, db *orm.DB) []StateChangeLog {
	var (
		changelogs []StateChangeLog
		scope      = db.NewScope(model)
	)

	db.Where("refer_table = ? AND refer_id = ?", scope.TableName(), GenerateReferenceKey(model, db)).Find(&changelogs)

	return changelogs
}

// GetLastStateChange gets last state change
func GetLastStateChange(model interface{}, db *orm.DB) *StateChangeLog {
	var (
		changelog StateChangeLog
		scope     = db.NewScope(model)
	)

	db.Where("refer_table = ? AND refer_id = ?", scope.TableName(), GenerateReferenceKey(model, db)).Last(&changelog)
	if changelog.To == "" {
		return nil
	}
	return &changelog
}

// ConfigureBhojpurResource used to configure transition for Bhojpur CMS admin
func (stageChangeLog *StateChangeLog) ConfigureBhojpurResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		if res.Permission == nil {
			res.Permission = roles.Deny(roles.Update, roles.Anyone).Deny(roles.Create, roles.Anyone)
		} else {
			res.Permission = res.Permission.Deny(roles.Update, roles.Anyone).Deny(roles.Create, roles.Anyone)
		}
	}
}
