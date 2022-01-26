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
	"errors"
	"fmt"
	"html/template"
	"path"
	"reflect"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/utils"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// SelectOneConfig meta configuration used for select one
type SelectOneConfig struct {
	Collection               interface{} // []string, [][]string, func(interface{}, *appsvr.Context) [][]string, func(interface{}, *admin.Context) [][]string
	Placeholder              string
	AllowBlank               bool
	DefaultCreating          bool
	SelectionTemplate        string
	SelectMode               string // select, select_async, bottom_sheet
	Select2ResultTemplate    template.JS
	Select2SelectionTemplate template.JS
	RemoteDataResource       *Resource
	RemoteDataHasImage       bool
	ForSerializedObject      bool
	PrimaryField             string
	metaConfig
	getCollection func(interface{}, *Context) [][]string
}

// GetPlaceholder get placeholder
func (selectOneConfig SelectOneConfig) GetPlaceholder(*Context) (template.HTML, bool) {
	return template.HTML(selectOneConfig.Placeholder), selectOneConfig.Placeholder != ""
}

// GetTemplate get template for selection template
func (selectOneConfig SelectOneConfig) GetTemplate(context *Context, metaType string) ([]byte, error) {
	if metaType == "form" && selectOneConfig.SelectionTemplate != "" {
		return context.Asset(selectOneConfig.SelectionTemplate)
	}
	return nil, errors.New("not implemented")
}

// GetCollection get collections from select one meta
func (selectOneConfig *SelectOneConfig) GetCollection(value interface{}, context *Context) [][]string {
	if selectOneConfig.getCollection == nil {
		selectOneConfig.prepareDataSource(nil, nil, "!remote_data_selector")
	}

	if selectOneConfig.getCollection != nil {
		return selectOneConfig.getCollection(value, context)
	}
	return [][]string{}
}

// ConfigureBhojpurMeta configure select one meta
func (selectOneConfig *SelectOneConfig) ConfigureBhojpurMeta(metaor resource.Metaor) {
	if meta, ok := metaor.(*Meta); ok {
		// Set FormattedValuer
		if meta.FormattedValuer == nil {
			meta.SetFormattedValuer(func(record interface{}, context *appsvr.Context) interface{} {
				return utils.Stringify(meta.GetValuer()(record, context))
			})
		}

		selectOneConfig.prepareDataSource(meta.FieldStruct, meta.baseResource, "!remote_data_selector")

		meta.Type = "select_one"
	}
}

// ConfigureBhojpurAdminFilter configure admin filter
func (selectOneConfig *SelectOneConfig) ConfigureBhojpurAdminFilter(filter *Filter) {
	var structField *orm.StructField
	if field, ok := filter.Resource.GetAdmin().DB.NewScope(filter.Resource.Value).FieldByName(filter.Name); ok {
		structField = field.StructField
	}

	selectOneConfig.prepareDataSource(structField, filter.Resource, "!remote_data_filter")

	if len(filter.Operations) == 0 {
		filter.Operations = []string{"equal"}
	}
	filter.Type = "select_one"
}

// FilterValue filter value
func (selectOneConfig *SelectOneConfig) FilterValue(filter *Filter, context *Context) interface{} {
	var (
		prefix  = fmt.Sprintf("filters[%v].", filter.Name)
		keyword string
	)

	if metaValues, err := resource.ConvertFormToMetaValues(context.Request, []resource.Metaor{}, prefix); err == nil {
		if metaValue := metaValues.Get("Value"); metaValue != nil {
			keyword = utils.ToString(metaValue.Value)
		}
	}

	if keyword != "" && selectOneConfig.RemoteDataResource != nil {
		result := selectOneConfig.RemoteDataResource.NewStruct()
		clone := context.Clone()
		clone.ResourceID = keyword
		if selectOneConfig.RemoteDataResource.CallFindOne(result, nil, clone) == nil {
			return result
		}
	}

	return keyword
}

func (selectOneConfig *SelectOneConfig) prepareDataSource(field *orm.StructField, res *Resource, routePrefix string) {
	// Set GetCollection
	if selectOneConfig.Collection != nil {
		selectOneConfig.SelectMode = "select"

		if values, ok := selectOneConfig.Collection.([]string); ok {
			selectOneConfig.getCollection = func(interface{}, *Context) (results [][]string) {
				for _, value := range values {
					results = append(results, []string{value, value})
				}
				return
			}
		} else if maps, ok := selectOneConfig.Collection.([][]string); ok {
			selectOneConfig.getCollection = func(interface{}, *Context) [][]string {
				return maps
			}
		} else if fc, ok := selectOneConfig.Collection.(func(interface{}, *appsvr.Context) [][]string); ok {
			selectOneConfig.getCollection = func(record interface{}, context *Context) [][]string {
				return fc(record, context.Context)
			}
		} else if fc, ok := selectOneConfig.Collection.(func(interface{}, *Context) [][]string); ok {
			selectOneConfig.getCollection = fc
		} else {
			utils.ExitWithMsg("Unsupported Collection format")
		}
	}

	// Set GetCollection if normal select mode
	if selectOneConfig.getCollection == nil {
		if selectOneConfig.RemoteDataResource == nil && field != nil {
			fieldType := field.Struct.Type
			for fieldType.Kind() == reflect.Ptr || fieldType.Kind() == reflect.Slice {
				fieldType = fieldType.Elem()
			}
			selectOneConfig.RemoteDataResource = res.GetAdmin().GetResource(fieldType.Name())
			if selectOneConfig.RemoteDataResource == nil {
				selectOneConfig.RemoteDataResource = res.GetAdmin().NewResource(reflect.New(fieldType).Interface())
			}
		}

		if selectOneConfig.PrimaryField == "" {
			for _, primaryField := range selectOneConfig.RemoteDataResource.PrimaryFields {
				selectOneConfig.PrimaryField = primaryField.Name
				break
			}
		}

		if selectOneConfig.SelectMode == "" {
			selectOneConfig.SelectMode = "select_async"
		}

		selectOneConfig.getCollection = func(_ interface{}, context *Context) (results [][]string) {
			cloneContext := context.clone()
			cloneContext.setResource(selectOneConfig.RemoteDataResource)
			searcher := &Searcher{Context: cloneContext}
			searcher.Pagination.CurrentPage = -1
			searchResults, _ := searcher.FindMany()

			reflectValues := reflect.Indirect(reflect.ValueOf(searchResults))
			for i := 0; i < reflectValues.Len(); i++ {
				value := reflectValues.Index(i).Interface()
				scope := context.GetDB().NewScope(value)

				obj := reflect.Indirect(reflect.ValueOf(value))
				idField := obj.FieldByName("ID")
				versionNameField := obj.FieldByName("VersionName")

				if idField.IsValid() && versionNameField.IsValid() {
					for i := 0; i < obj.Type().NumField(); i++ {
						// If given object has CompositePrimaryKey field, generate composite primary key and return it as the primary key.
						if obj.Type().Field(i).Name == resource.CompositePrimaryKeyFieldName {
							results = append(results, []string{resource.GenCompositePrimaryKey(idField.Uint(), versionNameField.String()), utils.Stringify(value)})
							continue
						}
					}

				}
				results = append(results, []string{fmt.Sprint(scope.PrimaryKeyValue()), utils.Stringify(value)})
			}
			return
		}
	}

	if res != nil && (selectOneConfig.SelectMode == "select_async" || selectOneConfig.SelectMode == "bottom_sheet") {
		if remoteDataResource := selectOneConfig.RemoteDataResource; remoteDataResource != nil {
			if !remoteDataResource.mounted {
				remoteDataResource.params = path.Join(routePrefix, res.ToParam(), field.Name, fmt.Sprintf("%p", remoteDataResource))
				res.GetAdmin().RegisterResourceRouters(remoteDataResource, "create", "update", "read", "delete")
			}
		} else {
			utils.ExitWithMsg("RemoteDataResource not configured")
		}
	}
}
