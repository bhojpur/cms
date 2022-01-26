package sorting_test

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
	"reflect"
	"testing"

	"github.com/bhojpur/cms/pkg/sorting"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type ColorVariation struct {
	orm.Model
	Code string
}

func checkOrder(results interface{}, order []string) error {
	values := reflect.Indirect(reflect.ValueOf(results))
	for idx, o := range order {
		value := values.Index(idx)
		primaryValue := fmt.Sprint(reflect.Indirect(value).FieldByName("ID").Interface())
		if primaryValue != o {
			return fmt.Errorf("#%v of values's primary key is %v, but should be %v", idx+1, primaryValue, o)
		}
	}
	return nil
}

func TestSortSlice(t *testing.T) {
	colorVariations := []ColorVariation{
		{Model: orm.Model{ID: 1}, Code: "1"},
		{Model: orm.Model{ID: 2}, Code: "2"},
		{Model: orm.Model{ID: 3}, Code: "3"},
	}

	collectionSorting := sorting.SortableCollection{PrimaryKeys: []string{"3", "1", "2"}}
	collectionSorting.Sort(colorVariations)

	if err := checkOrder(colorVariations, []string{"3", "1", "2"}); err != nil {
		t.Error(err)
	}
}

func TestSort(t *testing.T) {
	colorVariations := &[]ColorVariation{
		{Model: orm.Model{ID: 1}, Code: "1"},
		{Model: orm.Model{ID: 2}, Code: "2"},
		{Model: orm.Model{ID: 3}, Code: "3"},
	}

	collectionSorting := sorting.SortableCollection{PrimaryKeys: []string{"3", "1", "2"}}
	collectionSorting.Sort(colorVariations)

	if err := checkOrder(colorVariations, []string{"3", "1", "2"}); err != nil {
		t.Error(err)
	}
}

func TestSortPointer(t *testing.T) {
	colorVariations := &[]*ColorVariation{
		{Model: orm.Model{ID: 1}, Code: "1"},
		{Model: orm.Model{ID: 2}, Code: "2"},
		{Model: orm.Model{ID: 3}, Code: "3"},
	}

	collectionSorting := sorting.SortableCollection{PrimaryKeys: []string{"3", "1", "2"}}
	collectionSorting.Sort(colorVariations)

	if err := checkOrder(colorVariations, []string{"3", "1", "2"}); err != nil {
		t.Error(err)
	}
}

func TestSortWithSomePrimaryKeys(t *testing.T) {
	colorVariations := &[]ColorVariation{
		{Model: orm.Model{ID: 1}, Code: "1"},
		{Model: orm.Model{ID: 2}, Code: "2"},
		{Model: orm.Model{ID: 3}, Code: "3"},
		{Model: orm.Model{ID: 4}, Code: "4"},
	}

	collectionSorting := sorting.SortableCollection{PrimaryKeys: []string{"3", "1"}}
	collectionSorting.Sort(colorVariations)

	if err := checkOrder(colorVariations, []string{"3", "1", "2", "4"}); err != nil {
		t.Error(err)
	}
}

func TestSortPointerWithSomePrimaryKeys(t *testing.T) {
	colorVariations := &[]*ColorVariation{
		{Model: orm.Model{ID: 1}, Code: "1"},
		{Model: orm.Model{ID: 2}, Code: "2"},
		{Model: orm.Model{ID: 3}, Code: "3"},
		{Model: orm.Model{ID: 4}, Code: "4"},
	}

	collectionSorting := sorting.SortableCollection{PrimaryKeys: []string{"3", "1"}}
	collectionSorting.Sort(colorVariations)

	if err := checkOrder(colorVariations, []string{"3", "1", "2", "4"}); err != nil {
		t.Error(err)
	}
}
