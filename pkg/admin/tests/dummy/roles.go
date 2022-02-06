package dummy

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

	"github.com/bhojpur/application/pkg/roles"
)

const (
	Role_editor               = "Editor"
	Role_supervisor           = "Supervisor"
	Role_system_administrator = "admin"
	Role_developer            = "Developer"
)

// InitRoles initialize roles of the system.
func InitRoles() {
	roles.Register(Role_editor, func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && currentUser.(User).Role == Role_editor
	})
	roles.Register(Role_supervisor, func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && currentUser.(User).Role == Role_supervisor
	})
	roles.Register(Role_system_administrator, func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && currentUser.(User).Role == Role_system_administrator
	})
	roles.Register(Role_developer, func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && currentUser.(User).Role == Role_developer
	})
}
