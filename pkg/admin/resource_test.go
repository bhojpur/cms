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
	"fmt"
	"strings"
	"testing"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/fatih/color"
)

type Product struct {
	Name        string
	Code        string
	URL         string
	Description string
}

// Test Edit Attrs
type EditAttrsTestCase struct {
	Params []string
	Result []string
}

func TestEditAttrs(t *testing.T) {
	var testCases []EditAttrsTestCase
	testCases = append(testCases,
		EditAttrsTestCase{Params: []string{"Name", "Code"}, Result: []string{"Name", "Code"}},
		EditAttrsTestCase{Params: []string{"-Name"}, Result: []string{"Code", "URL", "Description"}},
		EditAttrsTestCase{Params: []string{"Name", "-Code"}, Result: []string{"Name"}},
		EditAttrsTestCase{Params: []string{"Name", "-Code", "-Name"}, Result: []string{}},
		EditAttrsTestCase{Params: []string{"Name", "Code", "-Name"}, Result: []string{"Code"}},
		EditAttrsTestCase{Params: []string{"-Name", "Code", "Name"}, Result: []string{"Code", "Name"}},
		EditAttrsTestCase{Params: []string{"Section:Name+Code+Description", "-Name"}, Result: []string{"Code+Description"}},
	)

	admin := New(&appsvr.Config{DB: db})
	product := admin.AddResource(&Product{})
	i := 1
	for _, testCase := range testCases {
		var attrs []interface{}
		for _, param := range testCase.Params {
			if strings.HasPrefix(param, "Section:") {
				var rows [][]string
				param = strings.Replace(param, "Section:", "", 1)
				rows = append(rows, strings.Split(param, "+"))
				attrs = append(attrs, &Section{Rows: rows})
			} else {
				attrs = append(attrs, param)
			}
		}

		editSections := product.EditAttrs(attrs...)
		var results []string
		for _, section := range editSections {
			columnStr := strings.Join(section.Rows[0], "+")
			results = append(results, columnStr)
		}
		if compareStringSlice(results, testCase.Result) {
			color.Green(fmt.Sprintf("Edit Attrs TestCase #%d: Success\n", i))
		} else {
			t.Errorf(color.RedString(fmt.Sprintf("\nEdit Attrs TestCase #%d: Failure Result:%v\n", i, results)))
		}
		i++
	}
}

func compareStringSlice(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	i := 0
	for _, s := range slice1 {
		if s != slice2[i] {
			return false
		}
		i++
	}
	return true
}
