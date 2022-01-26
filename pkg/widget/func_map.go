package widget

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
	"sort"
	"strings"

	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/cms/pkg/admin"
)

type GroupedWidgets struct {
	Group   string
	Widgets []*Widget
}

var funcMap = map[string]interface{}{
	"widget_available_scopes": func() []*Scope {
		if len(registeredScopes) > 0 {
			return append([]*Scope{{Name: "Default Visitor", Param: "default"}}, registeredScopes...)
		}
		return []*Scope{}
	},
	"widget_grouped_widgets": func(context *admin.Context) []*GroupedWidgets {
		groupedWidgetsSlice := []*GroupedWidgets{}

	OUTER:
		for _, w := range registeredWidgets {
			var roleNames = []interface{}{}
			for _, role := range context.Roles {
				roleNames = append(roleNames, role)
			}
			if w.Permission == nil || w.Permission.HasPermission(roles.Create, roleNames...) {
				for _, groupedWidgets := range groupedWidgetsSlice {
					if groupedWidgets.Group == w.Group {
						groupedWidgets.Widgets = append(groupedWidgets.Widgets, w)
						continue OUTER
					}
				}

				groupedWidgetsSlice = append(groupedWidgetsSlice, &GroupedWidgets{
					Group:   w.Group,
					Widgets: []*Widget{w},
				})
			}
		}

		sort.SliceStable(groupedWidgetsSlice, func(i, j int) bool {
			if groupedWidgetsSlice[i].Group == "" {
				return false
			}
			return strings.Compare(groupedWidgetsSlice[i].Group, groupedWidgetsSlice[j].Group) < 0
		})

		return groupedWidgetsSlice
	},
}
