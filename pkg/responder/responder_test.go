package responder_test

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
	"testing"

	"github.com/bhojpur/cms/pkg/responder"
)

func checkRespond(request *http.Request, format string, t *testing.T) {
	responder.With("html", func() {
		if format != "html" {
			t.Errorf("Should call %v, but called html", format)
		}
	}).With("json", func() {
		if format != "json" {
			t.Errorf("Should call %v, but called json", format)
		}
	}).With("xml", func() {
		if format != "xml" {
			t.Errorf("Should call %v, but called xml", format)
		}
	}).Respond(request)
}

func newRequestWithAcceptType(acceptType string) *http.Request {
	request, _ := http.NewRequest("GET", "", nil)
	request.Header.Add("Accept", acceptType)
	return request
}

func TestRespond(t *testing.T) {
	mimeMap := map[string]string{
		"text/html":        "html",
		"application/json": "json",
		"application/xml":  "xml",
	}

	for mimeType, format := range mimeMap {
		checkRespond(newRequestWithAcceptType(mimeType), format, t)
	}
}
