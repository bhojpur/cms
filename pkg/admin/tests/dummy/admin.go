package dummy

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

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/application/test/utils"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/media"
)

// NewDummyAdmin generate admin for Dummy Admin app
func NewDummyAdmin(keepData ...bool) *admin.Admin {
	var (
		db     = utils.TestDB()
		models = []interface{}{&User{}, &CreditCard{}, &Address{}, &Language{}, &Profile{}, &Phone{}, &Company{}}
		Admin  = admin.New(&admin.AdminConfig{Auth: DummyAuth{}, DB: db})
	)

	media.RegisterCallbacks(db)

	InitRoles()

	for _, value := range models {
		if len(keepData) == 0 {
			db.DropTableIfExists(value)
		}
		db.AutoMigrate(value)
	}

	c := Admin.AddResource(&Company{})
	c.Action(&admin.Action{
		Name: "Publish",
		Handler: func(argument *admin.ActionArgument) (err error) {
			fmt.Println("Publish company")
			return
		},
		Method:   "GET",
		Resource: c,
		Modes:    []string{"edit"},
	})
	c.Action(&admin.Action{
		Name:       "Preview",
		Permission: roles.Deny(roles.CRUD, Role_system_administrator),
		Handler: func(argument *admin.ActionArgument) (err error) {
			fmt.Println("Preview company")
			return
		},
		Method:   "GET",
		Resource: c,
		Modes:    []string{"edit"},
	})
	c.Action(&admin.Action{
		Name:       "Approve",
		Permission: roles.Allow(roles.Read, Role_system_administrator),
		Handler: func(argument *admin.ActionArgument) (err error) {
			fmt.Println("Approve company")
			return
		},
		Method: "GET",
		Modes:  []string{"edit"},
	})
	Admin.AddResource(&CreditCard{})

	Admin.AddResource(&Language{}, &admin.Config{Name: "语种 & 语言", Priority: -1})
	user := Admin.AddResource(&User{}, &admin.Config{Permission: roles.Allow(roles.CRUD, Role_system_administrator)})
	user.Meta(&admin.Meta{
		Name: "CreditCard",
		Type: "single_edit",
	})
	user.Meta(&admin.Meta{
		Name: "Languages",
		Type: "select_many",
		Collection: func(resource interface{}, context *appsvr.Context) (results [][]string) {
			if languages := []Language{}; !context.GetDB().Find(&languages).RecordNotFound() {
				for _, language := range languages {
					results = append(results, []string{fmt.Sprint(language.ID), language.Name})
				}
			}
			return
		},
	})

	return Admin
}

const LoggedInUserName = "Bhojpur"

type DummyAuth struct {
}

func (DummyAuth) LoginURL(ctx *admin.Context) string {
	return "/auth/login"
}

func (DummyAuth) LogoutURL(ctx *admin.Context) string {
	return "/auth/logout"
}

func (DummyAuth) GetCurrentUser(ctx *admin.Context) appsvr.CurrentUser {
	u := User{}

	if err := ctx.Admin.DB.Where("name = ?", LoggedInUserName).First(&u).Error; err != nil {
		fmt.Println("Cannot load logged in user", err.Error())
		return nil
	}

	return u
}
