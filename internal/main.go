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
	fmt.Println("Configuring the Bhojpur CMS demo application views")
	loginFS := demo.AssetFS.NameSpace("login")
	pluginFS := demo.AssetFS.NameSpace("plugin")
	// Register view paths into the loginFS
	loginFS.RegisterPath("templates/app/views")
	pluginFS.RegisterPath("templates/app/vendor/plugin/views")

	// Compile application templates under registered view paths into binary
	loginFS.Compile()
	pluginFS.Compile()

	fmt.Println("Configuring a Bhojpur CMS demo login templates")
	// Get file content with registered name
	loginContent, err := loginFS.Asset("login.html")
	if err != nil {
		fmt.Errorf("While configuring demo application login", err)
	}
	if loginContent != nil {
		fmt.Println("Configuring a demo application login.html")
	}

	fmt.Println("Configuring a Bhojpur CMS demo plugin templates")
	pluginContent, err := pluginFS.Asset("index.tmpl")
	// Get file content with registered name
	if err != nil {
		fmt.Errorf("While configuring demo application plugin", err)
	}
	if pluginContent != nil {
		fmt.Println("Configuring a demo application plugin.tmpl")
	}

	fmt.Println("Bhojpur CMS demo application server, opening SQL database")
	demoAppDB, _ := orm.Open("sqlite3", "internal/demo.db")
	demoAppDB.AutoMigrate(&User{}, &Product{})

	fmt.Println("Configuring a Bhojpur CMS administrator dashboard")
	// Initialize Bhojpur CMS - Administrator's Dashboard
	Admin := admin.New(&admin.AdminConfig{DB: demoAppDB})

	// Allow to use Bhojpur CMS - Admin to manage User, Product
	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})

	fmt.Println("Configuring an HTTP service request multiplexer")
	// initialize an HTTP request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/", landingPageFunc)

	fmt.Println("Mounting administrator dashboard web user interface")
	// Mount Bhojpur CMS - Administrator's web user interface to mux
	Admin.MountTo("/admin", mux)

	fmt.Println("Bhojpur CMS demo server listening at http://localhost:3000")
	http.ListenAndServe(":3000", mux)
}

func landingPageFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("<html><head><title>Bhojpur CMS</title></head><body><h1>Welcome to Bhojpur CMS</h1><p>Login to <a href=\"/admin\">Administrator's Dashboard</a></p><br/><br><hr size=\"1px\"><center>Copyright &copy; 2018 by <a href=\"https://www.bhojpur-consulting.com\">Bhojpur Consulting Private Limited</a>, India. All rights reserved.</center></body></html>"))

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
