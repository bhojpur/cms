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
	"testing"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/test/utils"
)

type User struct {
	Name string
	ID   uint64
}

var db = utils.TestDB()

func TestAddResource(t *testing.T) {
	admin := New(&appsvr.Config{DB: db})
	user := admin.AddResource(&User{})

	if user != admin.resources[0] {
		t.Error("resource not added")
	}

	if admin.GetMenus()[0].Name != "Users" {
		t.Error("resource not added to menu")
	}
}

func TestAddResourceWithInvisibleOption(t *testing.T) {
	admin := New(&appsvr.Config{DB: db})
	user := admin.AddResource(&User{}, &Config{Invisible: true})

	if user != admin.resources[0] {
		t.Error("resource not added")
	}

	if len(admin.GetMenus()) != 0 {
		t.Error("invisible resource registered in menu")
	}
}

func TestGetResource(t *testing.T) {
	admin := New(&appsvr.Config{DB: db})
	user := admin.AddResource(&User{})

	if admin.GetResource("User") != user {
		t.Error("resource not returned")
	}
}

func TestNewResource(t *testing.T) {
	admin := New(&appsvr.Config{DB: db})
	user := admin.NewResource(&User{})

	if user.Name != "User" {
		t.Error("default resource name didn't set")
	}
}

type UserWithCustomizedName struct{}

func (u *UserWithCustomizedName) ResourceName() string {
	return "CustomizedName"
}

func TestNewResourceWithCustomizedName(t *testing.T) {
	admin := New(&appsvr.Config{DB: db})
	user := admin.NewResource(&UserWithCustomizedName{})

	if user.Name != "CustomizedName" {
		t.Error("customize resource name didn't set")
	}
}
