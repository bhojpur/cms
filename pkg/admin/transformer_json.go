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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/bhojpur/application/pkg/roles"
)

// JSONTransformer json transformer
type JSONTransformer struct{}

// CouldEncode check if encodable
func (JSONTransformer) CouldEncode(encoder Encoder) bool {
	return true
}

// Encode encode encoder to writer as JSON
func (JSONTransformer) Encode(writer io.Writer, encoder Encoder) error {
	var (
		context = encoder.Context
		res     = encoder.Resource
	)

	js, err := json.MarshalIndent(convertObjectToJSONMap(res, context, encoder.Result, encoder.Action), "", "\t")
	if err != nil {
		result := make(map[string]string)
		result["error"] = err.Error()
		js, _ = json.Marshal(result)
	}

	if w, ok := writer.(http.ResponseWriter); ok {
		w.Header().Set("Content-Type", "application/json")
	}

	_, err = writer.Write(js)
	return err
}

func convertObjectToJSONMap(res *Resource, context *Context, value interface{}, kind string) interface{} {
	reflectValue := reflect.ValueOf(value)
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}

	switch reflectValue.Kind() {
	case reflect.Slice:
		values := []interface{}{}
		for i := 0; i < reflectValue.Len(); i++ {
			if reflect.Indirect(reflectValue.Index(i)).Kind() == reflect.Struct {
				if reflectValue.Index(i).Kind() == reflect.Ptr {
					values = append(values, convertObjectToJSONMap(res, context, reflectValue.Index(i).Interface(), kind))
				} else {
					values = append(values, convertObjectToJSONMap(res, context, reflectValue.Index(i).Addr().Interface(), kind))
				}
			} else {
				values = append(values, fmt.Sprint(reflectValue.Index(i).Interface()))
			}
		}
		return values
	case reflect.Struct:
		var metas []*Meta
		if kind == "index" {
			metas = res.ConvertSectionToMetas(res.allowedSections(res.IndexAttrs(), context, roles.Read))
		} else if kind == "edit" {
			metas = res.ConvertSectionToMetas(res.allowedSections(res.EditAttrs(), context, roles.Update))
		} else if kind == "show" {
			metas = res.ConvertSectionToMetas(res.allowedSections(res.ShowAttrs(), context, roles.Read))
		}

		values := map[string]interface{}{}
		for _, meta := range metas {
			if meta.HasPermission(roles.Read, context.Context) {
				// has_one, has_many checker to avoid dead loop
				if meta.Resource != nil && (meta.FieldStruct != nil && meta.FieldStruct.Relationship != nil && (meta.FieldStruct.Relationship.Kind == "has_one" || meta.FieldStruct.Relationship.Kind == "has_many" || meta.Type == "single_edit" || meta.Type == "collection_edit")) {
					values[meta.GetName()] = convertObjectToJSONMap(meta.Resource, context, context.RawValueOf(value, meta), kind)
				} else {
					values[meta.GetName()] = context.FormattedValueOf(value, meta)
				}
			}
		}
		return values
	case reflect.Map:
		for _, key := range reflectValue.MapKeys() {
			reflectValue.SetMapIndex(key, reflect.ValueOf(convertObjectToJSONMap(res, context, reflectValue.MapIndex(key).Interface(), kind)))
		}
		return reflectValue.Interface()
	default:
		return value
	}
}
