package responder

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

// Package responder respond differently according to request's accepted mime type

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// Register mime type and format
//     responder.Register("application/json", "json")
func Register(mimeType string, format string) {
	mime.AddExtensionType("."+strings.TrimPrefix(format, "."), mimeType)
}

func init() {
	for mimeType, format := range map[string]string{
		"text/html":        "html",
		"application/json": "json",
		"application/xml":  "xml",
	} {
		Register(mimeType, format)
	}
}

// Responder is holder of registed response handlers, response `Request` based on its accepted mime type
type Responder struct {
	responds         map[string]func()
	DefaultResponder func()
}

// With could be used to register response handler for mime type formats, the formats could be string or []string
//     responder.With("html", func() {
//       writer.Write([]byte("this is a html request"))
//     }).With([]string{"json", "xml"}, func() {
//       writer.Write([]byte("this is a json or xml request"))
//     })
func With(formats interface{}, fc func()) *Responder {
	rep := &Responder{responds: map[string]func(){}}
	return rep.With(formats, fc)
}

// With could be used to register response handler for mime type formats, the formats could be string or []string
func (rep *Responder) With(formats interface{}, fc func()) *Responder {
	if f, ok := formats.(string); ok {
		rep.responds[f] = fc
	} else if fs, ok := formats.([]string); ok {
		for _, f := range fs {
			rep.responds[f] = fc
		}
	}

	if rep.DefaultResponder == nil {
		rep.DefaultResponder = fc
	}
	return rep
}

// Respond differently according to request's accepted mime type
func (rep *Responder) Respond(request *http.Request) {
	// get request format from url
	if ext := filepath.Ext(request.URL.Path); ext != "" {
		if respond, ok := rep.responds[strings.TrimPrefix(ext, ".")]; ok {
			respond()
			return
		}
	}

	// get request format from Accept
	for _, accept := range strings.Split(request.Header.Get("Accept"), ",") {
		if exts, err := mime.ExtensionsByType(accept); err == nil {
			for _, ext := range exts {
				if respond, ok := rep.responds[strings.TrimPrefix(ext, ".")]; ok {
					respond()
					return
				}
			}
		}
	}

	// use first format as default
	if rep.DefaultResponder != nil {
		rep.DefaultResponder()
	}
	return
}
