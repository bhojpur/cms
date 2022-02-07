package main

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
	"net/http"

	"github.com/bhojpur/cms/internal/demo"
	"github.com/bhojpur/cms/pkg/admin"
	orm "github.com/bhojpur/orm/pkg/engine"
	_ "github.com/mattn/go-sqlite3"
)

// Create a simple Bhojpur ORM backend domain model
type User struct {
	orm.Model
	Name string
}

// Create another Bhojpur ORM backend domain model
type Product struct {
	orm.Model
	Name        string
	Description string
}

func main() {
	fmt.Println("Bhojpur CMS demo application server, opening SQL database")
	demoAppDB, _ := orm.Open("sqlite3", "internal/demo.db")
	demoAppDB.AutoMigrate(&User{}, &Product{})

	fmt.Println("Configuring the Bhojpur CMS demo application views")
	demoFS := demo.AssetFS
	// Register view paths into the demoFS
	demoFS.RegisterPath("app/views")
	demoFS.RegisterPath("app/vendor/plugin/views")

	// Compile application templates under registered view paths into binary
	demoFS.Compile()

	// Get file content with registered name
	fileContent, ok := demoFS.Asset("internal/app/views/index.tmpl")
	if ok != nil {
		fmt.Errorf("While configuring demo application", ok)
	}
	if fileContent != nil {
		fmt.Println("Configuring a demo home/index.html")
	}

	fmt.Println("Configuring a Bhojpur CMS administrator dashboard")
	// Initialize Bhojpur CMS - Administrator's Dashboard
	Admin := admin.New(&admin.AdminConfig{DB: demoAppDB})

	// Allow to use Bhojpur CMS - Admin to manage User, Product
	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})

	fmt.Println("Configuring an HTTP service request multiplexer")
	// initialize an HTTP request multiplexer
	mux := http.NewServeMux()

	fmt.Println("Mounting administrator dashboard web user interface")
	// Mount Bhojpur CMS - Administrator's web user interface to mux
	Admin.MountTo("/admin", mux)

	fmt.Println("Bhojpur CMS demo server listening at http://localhost:3000")
	http.ListenAndServe(":3000", mux)
}

/*
func rendering() {
	Render := render.New(&render.Config{
		ViewPaths:     []string{"internal/demo"},
		DefaultLayout: "application", // default value is application
		FuncMapMaker: func(*Render, *http.Request, http.ResponseWriter) template.FuncMap {
			// genereate FuncMap that could be used when render template based on request info
		},
	})
	Render.Execute("index", context, request, writer)
}
*/
