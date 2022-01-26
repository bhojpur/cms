package widget_test

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
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/test/utils"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/widget"
	orm "github.com/bhojpur/orm/pkg/engine"
	"github.com/fatih/color"
)

var db *orm.DB
var Widgets *widget.Widgets
var Admin *admin.Admin
var Server *httptest.Server

type bannerArgument struct {
	Title    string
	SubTitle string
}

func init() {
	db = utils.TestDB()
	//db.LogMode(true)
}

// Runner
func setup() {
	if err := db.DropTableIfExists(&widget.BhojpurWidgetSetting{}).Error; err != nil {
		panic(err)
	}
	db.AutoMigrate(&widget.BhojpurWidgetSetting{})
	mux := http.NewServeMux()
	Server = httptest.NewServer(mux)

	Widgets = widget.New(&widget.Config{
		DB: db,
	})
	Widgets.RegisterViewPath("github.com/bhojpur/cms/pkg/widget/test")

	Admin = admin.New(&appsvr.Config{DB: db})
	Admin.AddResource(Widgets)
	Admin.MountTo("/admin", mux)

	Widgets.RegisterWidget(&widget.Widget{
		Name:      "Banner",
		Templates: []string{"banner"},
		Setting:   Admin.NewResource(&bannerArgument{}),
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			if setting != nil {
				argument := setting.(*bannerArgument)
				context.Options["Title"] = argument.Title
				context.Options["SubTitle"] = argument.SubTitle
			}
			return context
		},
	})

	Widgets.RegisterScope(&widget.Scope{
		Name: "From Google",
		Visible: func(context *widget.Context) bool {
			if request, ok := context.Get("Request"); ok {
				_, ok := request.(*http.Request).URL.Query()["from_google"]
				return ok
			}
			return false
		},
	})

	Widgets.RegisterWidget(&widget.Widget{
		Name:    "NoTemplate",
		Setting: Admin.NewResource(&bannerArgument{}),
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			context.Body = "<h1>My html</h1>"
			return context
		},
	})
}

func reset() {
	db.DropTable(&widget.BhojpurWidgetSetting{})
	db.AutoMigrate(&widget.BhojpurWidgetSetting{})
	setup()
}

// Test DB's record after call Render
func TestRenderRecord(t *testing.T) {
	reset()
	var count int
	db.Model(&widget.BhojpurWidgetSetting{}).Where(widget.BhojpurWidgetSetting{Name: "HomeBanner", WidgetType: "Banner", Scope: "default", GroupName: "Banner"}).Count(&count)
	if count != 0 {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render Record TestCase: should don't exist widget setting")))
	}

	widgetContext := Widgets.NewContext(&widget.Context{})
	widgetContext.Render("HomeBanner", "Banner")
	db.Model(&widget.BhojpurWidgetSetting{}).Where(widget.BhojpurWidgetSetting{Name: "HomeBanner", WidgetType: "Banner", Scope: "default", GroupName: "Banner"}).Count(&count)
	if count == 0 {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render Record TestCase: should have default widget setting")))
	}

	http.PostForm(Server.URL+"/admin/widgets/HomeBanner",
		url.Values{"_method": {"PUT"},
			"BhojpurResource.Scope":       {"from_google"},
			"BhojpurResource.ActivatedAt": {"2018-03-26 10:10:42.433372925 +0800 CST"},
			"BhojpurResource.Widgets":     {"Banner"},
			"BhojpurResource.Template":    {"banner"},
			"BhojpurResource.Kind":        {"Banner"},
		})
	db.Model(&widget.BhojpurWidgetSetting{}).Where(widget.BhojpurWidgetSetting{Name: "HomeBanner", WidgetType: "Banner", Scope: "from_google"}).Count(&count)
	if count == 0 {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render Record TestCase: should have from_google widget setting")))
	}
}

// Runner
func TestRenderContext(t *testing.T) {
	reset()
	setting := &widget.BhojpurWidgetSetting{}
	db.Where(widget.BhojpurWidgetSetting{Name: "HomeBanner", WidgetType: "Banner", Scope: "default"}).FirstOrInit(setting)
	db.Create(setting)

	html := Widgets.Render("HomeBanner", "Banner")
	if !strings.Contains(string(html), "Hello, \n<h1></h1>\n<h2></h2>\n") {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result:\n %s\n", 1, html)))
	}

	widgetContext := Widgets.NewContext(&widget.Context{
		Options: map[string]interface{}{"CurrentUser": "Bhojpur"},
	})
	html = widgetContext.Render("HomeBanner", "Banner")
	if !strings.Contains(string(html), "Hello, Bhojpur\n<h1></h1>\n<h2></h2>\n") {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result:\n %s\n", 2, html)))
	}

	db.Where(widget.BhojpurWidgetSetting{Name: "HomeBanner", WidgetType: "Banner"}).FirstOrInit(setting)
	setting.SetSerializableArgumentValue(&bannerArgument{Title: "Title", SubTitle: "SubTitle"})
	err := db.Model(setting).Update(setting).Error
	if err != nil {
		panic(err)
	}

	html = widgetContext.Render("HomeBanner", "Banner")
	if !strings.Contains(string(html), "Hello, Bhojpur\n<h1>Title</h1>\n<h2>SubTitle</h2>\n") {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result:\n %s\n", 3, html)))
	}

}

func TestRenderNoTemplate(t *testing.T) {
	reset()

	html := Widgets.Render("abc", "NoTemplate")
	if !strings.Contains(string(html), "<h1>My html</h1>") {
		t.Errorf(color.RedString(fmt.Sprintf("\nWidget Render TestCase #%d: Failure Result:\n %s\n", 5, html)))
	}

}

func TestRegisterFuncMap(t *testing.T) {
	func1 := func() {}
	Widgets.RegisterFuncMap("func1", func1)
	context := Widgets.NewContext(nil)
	if _, ok := context.FuncMaps["func1"]; !ok {
		t.Errorf("func1 should be assigned to context")
	}

	context2 := context.Funcs(template.FuncMap{"func2": func() {}})
	if _, ok := context.FuncMaps["func2"]; !ok {
		t.Errorf("func2 should be assigned to context")
	}
	if _, ok := context2.FuncMaps["func2"]; !ok {
		t.Errorf("func2 should be assigned to context")
	}

	context3 := Widgets.NewContext(nil)
	if _, ok := context3.FuncMaps["func3"]; ok {
		t.Errorf("func3 should not be assigned to other contexts")
	}
}
