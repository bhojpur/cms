//go:build !iris
// +build !iris

package iris

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
	"github.com/bhojpur/cms/pkg/admin"
	orm "github.com/bhojpur/orm/pkg/engine"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kataras/iris/v12"
)

// Create a ORM-backend model
type User struct {
	orm.Model
	Name string
}

// Create another ORM-backend model
type Product struct {
	orm.Model
	Name        string
	Description string
}

func main() {
	DB, _ := orm.Open("sqlite3", "demo.db")
	DB.AutoMigrate(&User{}, &Product{})

	bhojpurPrefix := "/admin"
	// Initialize Bhojpur CMS - Admin module.
	Admin := admin.New(&admin.AdminConfig{DB: DB})

	// Allow to use Admin to manage User, Product
	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})
	// Create a Bhojpur handler and convert it to an iris one with `iris.FromStd`.
	handler := iris.FromStd(Admin.NewServeMux(bhojpurPrefix))

	// Initialize Iris.
	app := iris.New()
	// Mount routes for "/admin" and "/admin/:xxx/..."
	app.Any(bhojpurPrefix, handler)
	app.Any(bhojpurPrefix+"/{p:path}", handler)

	// Start the server.
	// Navigate at: http://localhost:9000/admin.
	app.Listen(":9000")
}
