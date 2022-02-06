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
	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/utils"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// Filter filter definiation
type Filter struct {
	Name       string
	Label      string
	Type       string
	Operations []string // eq, cont, gt, gteq, lt, lteq
	Resource   *Resource
	Visible    func(context *Context) bool
	Handler    func(*orm.DB, *FilterArgument) *orm.DB
	Config     FilterConfigInterface
}

// SavedFilter saved filter settings
type SavedFilter struct {
	Name string
	URL  string
}

// FilterConfigInterface filter config interface
type FilterConfigInterface interface {
	ConfigureBhojpurAdminFilter(*Filter)
}

// FilterArgument filter argument that used in handler
type FilterArgument struct {
	Value    *resource.MetaValues
	Resource *Resource
	Context  *appsvr.Context
}

// Filter register filter for Bhojpur CMS resource
func (res *Resource) Filter(filter *Filter) {
	filter.Resource = res

	if filter.Label == "" {
		filter.Label = utils.HumanizeString(filter.Name)
	}

	if meta := res.GetMeta(filter.Name); meta != nil {
		if filter.Type == "" {
			filter.Type = meta.Type
		}

		if filter.Config == nil {
			if filterConfig, ok := meta.Config.(FilterConfigInterface); ok {
				filter.Config = filterConfig
			}
		}
	}

	if filter.Config != nil {
		filter.Config.ConfigureBhojpurAdminFilter(filter)
	}

	if filter.Handler == nil {
		// generate default handler
		filter.Handler = func(db *orm.DB, filterArgument *FilterArgument) *orm.DB {
			if metaValue := filterArgument.Value.Get("Value"); metaValue != nil {
				keyword := utils.ToString(metaValue.Value)
				if _, ok := filter.Config.(*SelectManyConfig); ok {
					if arr, ok := metaValue.Value.([]string); ok {
						keyword = strings.Join(arr, ",")
					}
				}

				if keyword != "" {
					field := filterField{FieldName: filter.Name}
					if operationMeta := filterArgument.Value.Get("Operation"); operationMeta != nil {
						if operation := utils.ToString(operationMeta.Value); operation != "" {
							field.Operation = operation
						}
					}
					if field.Operation == "" {
						if len(filter.Operations) > 0 {
							field.Operation = filter.Operations[0]
						} else {
							field.Operation = "contains"
						}
					}

					return filterResourceByFields(res, []filterField{field}, keyword, db, filterArgument.Context)
				}
			}
			return db
		}
	}

	if filter.Type != "" {
		res.filters = append(res.filters, filter)
	} else {
		utils.ExitWithMsg("Invalid filter definition %v for resource %v", filter.Name, res.Name)
	}
}

// GetFilters get registered filters
func (res *Resource) GetFilters() []*Filter {
	return res.filters
}

// GetFilter get defined action
func (res *Resource) GetFilter(name string) *Filter {
	for _, action := range res.filters {
		if action.Name == name {
			return action
		}
	}
	return nil
}
