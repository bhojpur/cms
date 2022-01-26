package admin

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

	appsvr "github.com/bhojpur/application/pkg/engine"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// SettingsStorageInterface settings storage interface
type SettingsStorageInterface interface {
	Get(key string, value interface{}, context *Context) error
	Save(key string, value interface{}, res *Resource, user appsvr.CurrentUser, context *Context) error
}

func newSettings(db *orm.DB) SettingsStorageInterface {
	if db != nil {
		db.AutoMigrate(&BhojpurAdminSetting{})
	}
	return settings{}
}

// BhojpurAdminSetting admin settings
type BhojpurAdminSetting struct {
	orm.Model
	Key      string
	Resource string
	UserID   string
	Value    string `orm:"size:65532"`
}

type settings struct{}

// Get load admin settings
func (settings) Get(key string, value interface{}, context *Context) error {
	var (
		settings  = []BhojpurAdminSetting{}
		tx        = context.GetDB().New()
		resParams = ""
		userID    = ""
	)
	sqlCondition := fmt.Sprintf("%v = ? AND (resource = ? OR resource = ?) AND (user_id = ? OR user_id = ?)", tx.NewScope(nil).Quote("key"))

	if context.Resource != nil {
		resParams = context.Resource.ToParam()
	}

	if context.CurrentUser != nil {
		userID = ""
	}

	tx.Where(sqlCondition, key, resParams, "", userID, "").Order("user_id DESC, resource DESC, id DESC").Find(&settings)

	for _, setting := range settings {
		if err := json.Unmarshal([]byte(setting.Value), value); err != nil {
			return err
		}
	}

	return nil
}

// Save save admin settings
func (settings) Save(key string, value interface{}, res *Resource, user appsvr.CurrentUser, context *Context) error {
	var (
		tx          = context.GetDB().New()
		result, err = json.Marshal(value)
		resParams   = ""
		userID      = ""
	)

	if err != nil {
		return err
	}

	if res != nil {
		resParams = res.ToParam()
	}

	if user != nil {
		userID = ""
	}

	err = tx.Where(BhojpurAdminSetting{
		Key:      key,
		UserID:   userID,
		Resource: resParams,
	}).Assign(BhojpurAdminSetting{Value: string(result)}).FirstOrCreate(&BhojpurAdminSetting{}).Error

	return err
}
