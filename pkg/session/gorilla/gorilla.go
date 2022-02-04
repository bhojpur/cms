package gorilla

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
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bhojpur/application/pkg/utils"
	session "github.com/bhojpur/cms/pkg/session"
	gorillaContext "github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

// New initialize session manager for gorilla
func New(sessionName string, store sessions.Store) *Gorilla {
	return &Gorilla{SessionName: sessionName, Store: store}
}

// Gorilla session manager struct for gorilla
type Gorilla struct {
	SessionName string
	Store       sessions.Store
}

var reader utils.ContextKey = "gorilla_reader"

func (gorilla Gorilla) getSession(req *http.Request) (*sessions.Session, error) {
	if r, ok := req.Context().Value(reader).(*http.Request); ok {
		return gorilla.Store.Get(r, gorilla.SessionName)
	}
	return gorilla.Store.Get(req, gorilla.SessionName)
}

func (gorilla Gorilla) saveSession(w http.ResponseWriter, req *http.Request) {
	if session, err := gorilla.getSession(req); err == nil {
		if err := session.Save(req, w); err != nil {
			fmt.Printf("No error should happen when saving session data, but got %v", err)
		}
	}
}

// Add value to session data, if value is not string, will marshal it into JSON encoding and save it into session data.
func (gorilla Gorilla) Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error {
	defer gorilla.saveSession(w, req)

	session, err := gorilla.getSession(req)
	if err != nil {
		return err
	}

	if str, ok := value.(string); ok {
		session.Values[key] = str
	} else {
		result, _ := json.Marshal(value)
		session.Values[key] = string(result)
	}

	return nil
}

// Pop value from session data
func (gorilla Gorilla) Pop(w http.ResponseWriter, req *http.Request, key string) string {
	defer gorilla.saveSession(w, req)

	if session, err := gorilla.getSession(req); err == nil {
		if value, ok := session.Values[key]; ok {
			delete(session.Values, key)
			return fmt.Sprint(value)
		}
	}
	return ""
}

// Get value from session data
func (gorilla Gorilla) Get(req *http.Request, key string) string {
	if session, err := gorilla.getSession(req); err == nil {
		if value, ok := session.Values[key]; ok {
			return fmt.Sprint(value)
		}
	}
	return ""
}

// Flash add flash message to session data
func (gorilla Gorilla) Flash(w http.ResponseWriter, req *http.Request, message session.Message) error {
	var messages []session.Message
	if err := gorilla.Load(req, "_flashes", &messages); err != nil {
		return err
	}
	messages = append(messages, message)
	return gorilla.Add(w, req, "_flashes", messages)
}

// Flashes returns a slice of flash messages from session data
func (gorilla Gorilla) Flashes(w http.ResponseWriter, req *http.Request) []session.Message {
	var messages []session.Message
	gorilla.PopLoad(w, req, "_flashes", &messages)
	return messages
}

// Load get value from session data and unmarshal it into result
func (gorilla Gorilla) Load(req *http.Request, key string, result interface{}) error {
	value := gorilla.Get(req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

// PopLoad pop value from session data and unmarshal it into result
func (gorilla Gorilla) PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error {
	value := gorilla.Pop(w, req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

// Middleware returns a new session manager middleware instance
func (gorilla Gorilla) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer gorillaContext.Clear(req)
		ctx := context.WithValue(req.Context(), reader, req)
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}
