package session

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
	"html/template"
	"net/http"
)

// ManagerInterface session manager interface
type ManagerInterface interface {
	// Add value to session data, if value is not string, will marshal it into JSON encoding and save it into session data.
	Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error
	// Get value from session data
	Get(req *http.Request, key string) string
	// Pop value from session data
	Pop(w http.ResponseWriter, req *http.Request, key string) string

	// Flash add flash message to session data
	Flash(w http.ResponseWriter, req *http.Request, message Message) error
	// Flashes returns a slice of flash messages from session data
	Flashes(w http.ResponseWriter, req *http.Request) []Message

	// Load get value from session data and unmarshal it into result
	Load(req *http.Request, key string, result interface{}) error
	// PopLoad pop value from session data and unmarshal it into result
	PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error

	// Middleware returns a new session manager middleware instance.
	Middleware(http.Handler) http.Handler
}

// Message message struct
type Message struct {
	Message template.HTML
	Type    string
}
