package bhojpur_session

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
	"net/http/httptest"

	"github.com/bhojpur/application/pkg/utils"
	"github.com/bhojpur/cms/pkg/session"
	bhojpur_session "github.com/bhojpur/session/pkg/engine"
)

var writer utils.ContextKey = "gorilla_writer"

// New initialize session manager for BhojpurSession
func New(engine *bhojpur_session.Manager) *BhojpurSession {
	return &BhojpurSession{Manager: engine}
}

// BhojpurSession session manager struct for BhojpurSession
type BhojpurSession struct {
	*bhojpur_session.Manager
}

func (bhojpursession BhojpurSession) getSession(w http.ResponseWriter, req *http.Request) (bhojpur_session.Store, error) {
	return bhojpursession.Manager.SessionStart(w, req)
}

// Add value to session data, if value is not string, will marshal it into JSON encoding and save it into session data.
func (bhojpursession BhojpurSession) Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error {
	sess, _ := bhojpursession.getSession(w, req)
	defer sess.SessionRelease(req.Context(), w)

	if str, ok := value.(string); ok {
		return sess.Set(req.Context(), key, str)
	}
	result, _ := json.Marshal(value)
	return sess.Set(req.Context(), key, string(result))
}

// Pop value from session data
func (bhojpursession BhojpurSession) Pop(w http.ResponseWriter, req *http.Request, key string) string {
	sess, _ := bhojpursession.getSession(w, req)
	defer sess.SessionRelease(req.Context(), w)

	result := sess.Get(req.Context(), key)

	sess.Delete(req.Context(), key)
	if result != nil {
		return fmt.Sprint(result)
	}
	return ""
}

// Get value from session data
func (bhojpursession BhojpurSession) Get(req *http.Request, key string) string {
	sess, _ := bhojpursession.getSession(httptest.NewRecorder(), req)

	result := sess.Get(req.Context(), key)
	if result != nil {
		return fmt.Sprint(result)
	}
	return ""
}

// Flash add flash message to session data
func (bhojpursession BhojpurSession) Flash(w http.ResponseWriter, req *http.Request, message session.Message) error {
	var messages []session.Message
	if err := bhojpursession.Load(req, "_flashes", &messages); err != nil {
		return err
	}
	messages = append(messages, message)
	return bhojpursession.Add(w, req, "_flashes", messages)
}

// Flashes returns a slice of flash messages from session data
func (bhojpursession BhojpurSession) Flashes(w http.ResponseWriter, req *http.Request) []session.Message {
	var messages []session.Message
	bhojpursession.PopLoad(w, req, "_flashes", &messages)
	return messages
}

// Load get value from session data and unmarshal it into result
func (bhojpursession BhojpurSession) Load(req *http.Request, key string, result interface{}) error {
	value := bhojpursession.Get(req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

// PopLoad pop value from session data and unmarshal it into result
func (bhojpursession BhojpurSession) PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error {
	value := bhojpursession.Pop(w, req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

// Middleware returns a new session manager middleware instance
func (bhojpursession BhojpurSession) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), writer, w)
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}
