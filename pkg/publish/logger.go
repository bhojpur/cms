package publish

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
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	orm "github.com/bhojpur/orm/pkg/engine"
)

// LoggerInterface logger interface used to print publish logs
type LoggerInterface interface {
	Print(...interface{})
}

// Logger default logger used to print publish logs
var Logger LoggerInterface

func init() {
	Logger = log.New(os.Stdout, "\r\n", 0)
}

func stringify(object interface{}) string {
	if obj, ok := object.(interface {
		Stringify() string
	}); ok {
		return obj.Stringify()
	}

	scope := orm.Scope{Value: object}
	for _, column := range []string{"Description", "Name", "Title", "Code"} {
		if field, ok := scope.FieldByName(column); ok {
			return fmt.Sprintf("%v", field.Field.Interface())
		}
	}

	if scope.PrimaryField() != nil {
		if scope.PrimaryKeyZero() {
			return ""
		}
		return fmt.Sprintf("%v#%v", scope.GetModelStruct().ModelType.Name(), scope.PrimaryKeyValue())
	}

	return fmt.Sprint(reflect.Indirect(reflect.ValueOf(object)).Interface())
}

func stringifyPrimaryValues(primaryValues [][][]interface{}, columns ...string) string {
	var values []string
	for _, primaryValue := range primaryValues {
		var primaryKeys []string
		for _, value := range primaryValue {
			if len(columns) == 0 {
				primaryKeys = append(primaryKeys, fmt.Sprint(value[1]))
			} else {
				for _, column := range columns {
					if column == fmt.Sprint(value[0]) {
						primaryKeys = append(primaryKeys, fmt.Sprint(value[1]))
					}
				}
			}
		}
		if len(primaryKeys) > 1 {
			values = append(values, fmt.Sprintf("[%v]", strings.Join(primaryKeys, ", ")))
		} else {
			values = append(values, primaryKeys...)
		}
	}
	return strings.Join(values, "; ")
}
