package render

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
	"io/ioutil"
	"regexp"
	"testing"

	"net/http/httptest"
	"net/textproto"
)

func TestExecute(t *testing.T) {
	Render := New(nil, "test")

	request := httptest.NewRequest("GET", "/test", nil)
	responseWriter := httptest.NewRecorder()
	var context interface{}

	tmpl := Render.Layout("layout_for_test")
	tmpl.Execute("test", context, request, responseWriter)

	if textproto.TrimString(responseWriter.Body.String()) != "Template for test" {
		t.Errorf("The template isn't rendered")
	}
}

func TestErrorMessageWhenMissingLayout(t *testing.T) {
	Render := New(nil, "test")

	request := httptest.NewRequest("GET", "/test", nil)
	responseWriter := httptest.NewRecorder()
	var context interface{}

	nonExistLayout := "ThePlant"
	tmpl := Render.Layout(nonExistLayout)
	err := tmpl.Execute("test", context, request, responseWriter)
	if err != nil {
		t.Error("we don't return error, we render the error on page instead")
	}

	bodyBytes, err1 := ioutil.ReadAll(responseWriter.Result().Body)
	if err1 != nil {
		t.Fatal(err1)
	}
	bodyString := string(bodyBytes)

	errorRegexp := "Failed to render page.+" + nonExistLayout + ".*"

	if matched, _ := regexp.MatchString(errorRegexp, bodyString); !matched {
		t.Errorf("Missing layout error message is incorrect")
	}
}

func TestErrorMessageWhenLayoutContainsError(t *testing.T) {
	Render := New(nil, "test")

	request := httptest.NewRequest("GET", "/test", nil)
	responseWriter := httptest.NewRecorder()
	var context interface{}

	layoutContainsError := "layout_contains_error"
	tmpl := Render.Layout(layoutContainsError)
	err := tmpl.Execute("test", context, request, responseWriter)

	if err != nil {
		t.Error("we don't return error, we render the error on page instead")
	}

	bodyBytes, err1 := ioutil.ReadAll(responseWriter.Result().Body)
	if err1 != nil {
		t.Fatal(err1)
	}
	bodyString := string(bodyBytes)

	errorRegexp := "Failed to render page.+" + layoutContainsError + ".*"

	if matched, _ := regexp.MatchString(errorRegexp, bodyString); !matched {
		t.Errorf("Missing layout error message is incorrect")
	}
}
