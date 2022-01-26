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
	"fmt"
	"html/template"
	"testing"

	appsvr "github.com/bhojpur/application/pkg/engine"
	orm "github.com/bhojpur/orm/pkg/engine"
	"github.com/fatih/color"
)

type rawTestCase struct {
	HTML         string
	ExpectResult string
}

func TestFuncMaps(t *testing.T) {
	rawTestCases := []rawTestCase{
		{HTML: "<a href='#'>Hello</a>", ExpectResult: "Hello"},
		{HTML: "<a href='http://www.google.com'>Hello</a>", ExpectResult: "<a href=\"http://www.google.com\" rel=\"nofollow\">Hello</a>"},
		{HTML: "<a href='http://www.google.com' data-hint='Hello'>Hello</a>", ExpectResult: "<a href=\"http://www.google.com\" rel=\"nofollow\">Hello</a>"},
	}

	unsafeRawTestCases := []rawTestCase{
		{HTML: "<a href='http://g.cn'>Hello</a>", ExpectResult: "<a href='http://g.cn'>Hello</a>"},
		{HTML: "<a href='#' data-hint='Hello'>Hello</a>", ExpectResult: "<a href='#' data-hint='Hello'>Hello</a>"},
	}

	context := Context{
		Admin: New(&appsvr.Config{}),
	}
	funcMaps := context.FuncMap()

	for i, testcase := range rawTestCases {
		result := funcMaps["raw"].((func(string) template.HTML))(testcase.HTML)
		var hasError bool
		if result != template.HTML(testcase.ExpectResult) {
			t.Errorf(color.RedString(fmt.Sprintf("Admin FuncMap raw #%v: expect get %v, but got '%v'", i+1, testcase.ExpectResult, result)))
			hasError = true
		}
		if !hasError {
			fmt.Printf(color.GreenString(fmt.Sprintf("Admin FuncMap raw #%v: Success\n", i+1)))
		}
	}

	for i, testcase := range unsafeRawTestCases {
		result := funcMaps["unsafe_raw"].((func(string) template.HTML))(testcase.HTML)
		var hasError bool
		if result != template.HTML(testcase.ExpectResult) {
			t.Errorf(color.RedString(fmt.Sprintf("Admin FuncMap unsafe_raw #%v: expect get %v, but got '%v'", i+1, testcase.ExpectResult, result)))
			hasError = true
		}
		if !hasError {
			fmt.Printf(color.GreenString(fmt.Sprintf("Admin FuncMap unsafe_raw #%v: Success\n", i+1)))
		}
	}
}

type FakeStruct struct {
	orm.Model
	Name string
}

func TestIsEqual(t *testing.T) {
	c1 := FakeStruct{Name: "c1"}
	c1.ID = 1
	c2 := FakeStruct{Name: "c2"}
	c2.ID = 1

	context := Context{
		Admin: New(&appsvr.Config{}),
	}
	if !context.isEqual(c1, c2) {
		t.Error("same primary key is not equal")
	}

	c1.ID = 2
	if context.isEqual(c1, c2) {
		t.Error("different primary key is equal")
	}

	a := "a test"
	b := "another one"
	if context.isEqual(a, b) {
		t.Error("different string is equal")
	}

	c := 11
	d := 11
	if !context.isEqual(c, d) {
		t.Error("same int is not equal")
	}
}
