package admin_test

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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bhojpur/cms/pkg/admin"
	. "github.com/bhojpur/cms/tests/dummy"
	orm "github.com/bhojpur/orm/pkg/engine"
)

var (
	server       *httptest.Server
	db           *orm.DB
	Admin        *admin.Admin
	adminHandler http.Handler
)

func init() {
	Admin = NewDummyAdmin()
	adminHandler = Admin.NewServeMux("/admin")
	db = Admin.DB
	server = httptest.NewServer(adminHandler)
}

func TestMain(m *testing.M) {
	// Create universal logged-in user for test.
	createLoggedInUser()
	retCode := m.Run()

	os.Exit(retCode)
}

func createLoggedInUser() *User {
	user := User{Name: LoggedInUserName, Role: Role_system_administrator}
	if err := db.Save(&user).Error; err != nil {
		panic(err)
	}

	return &user
}