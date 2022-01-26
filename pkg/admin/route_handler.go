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
	"strings"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/roles"
)

// RouteConfig config for admin routes
type RouteConfig struct {
	Resource       *Resource
	Permissioner   HasPermissioner
	PermissionMode roles.PermissionMode
	Values         map[interface{}]interface{}
}

type requestHandler func(c *Context)

type routeHandler struct {
	Path   string
	Handle requestHandler
	Config *RouteConfig
}

func newRouteHandler(path string, handle requestHandler, configs ...*RouteConfig) *routeHandler {
	handler := &routeHandler{
		Path:   "/" + strings.TrimPrefix(path, "/"),
		Handle: handle,
	}

	for _, config := range configs {
		handler.Config = config
	}

	if handler.Config == nil {
		handler.Config = &RouteConfig{}
	}

	if handler.Config.Resource != nil {
		if handler.Config.Permissioner == nil {
			handler.Config.Permissioner = handler.Config.Resource
		}

		handler.Config.Resource.mounted = true
	}
	return handler
}

// HasPermission check has permission to access router handler or not
func (handler routeHandler) HasPermission(permissionMode roles.PermissionMode, context *appsvr.Context) bool {
	if handler.Config.Permissioner == nil {
		return true
	}

	if handler.Config.PermissionMode != "" {
		return handler.Config.Permissioner.HasPermission(handler.Config.PermissionMode, context)
	}

	return handler.Config.Permissioner.HasPermission(permissionMode, context)
}
