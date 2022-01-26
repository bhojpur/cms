package scs

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
	"encoding/json"
	"net/http"

	"github.com/alexedwards/scs"
	session "github.com/bhojpur/session/pkg/engine"
)

// New initialize session manager for SCS
func New(manager *scs.Manager) *SCS {
	return &SCS{Manager: manager}
}

// SCS session manager struct for SCS
type SCS struct {
	*scs.Manager
}

// Add value to session data, if value is not string, will marshal it into JSON encoding and save it into session data.
func (scs SCS) Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error {
	scssession := scs.Manager.Load(req)

	if str, ok := value.(string); ok {
		return scssession.PutString(w, key, str)
	}
	result, _ := json.Marshal(value)
	return scssession.PutString(w, key, string(result))
}

// Pop value from session data
func (scs SCS) Pop(w http.ResponseWriter, req *http.Request, key string) string {
	scssession := scs.Manager.Load(req)
	result, _ := scssession.PopString(w, key)
	return result
}

// Get value from session data
func (scs SCS) Get(req *http.Request, key string) string {
	scssession := scs.Manager.Load(req)
	result, _ := scssession.GetString(key)
	return result
}

// Flash add flash message to session data
func (scs SCS) Flash(w http.ResponseWriter, req *http.Request, message session.Message) error {
	var messages []session.Message
	if err := scs.Load(req, "_flashes", &messages); err != nil {
		return err
	}
	messages = append(messages, message)
	return scs.Add(w, req, "_flashes", messages)
}

// Flashes returns a slice of flash messages from session data
func (scs SCS) Flashes(w http.ResponseWriter, req *http.Request) []session.Message {
	var messages []session.Message
	scs.PopLoad(w, req, "_flashes", &messages)
	return messages
}

// Load get value from session data and unmarshal it into result
func (scs SCS) Load(req *http.Request, key string, result interface{}) error {
	value := scs.Get(req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

// PopLoad pop value from session data and unmarshal it into result
func (scs SCS) PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error {
	value := scs.Pop(w, req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

// Middleware returns a new session manager middleware instance
func (scs SCS) Middleware(handler http.Handler) http.Handler {
	return scs.Manager.Multi(handler)
}
