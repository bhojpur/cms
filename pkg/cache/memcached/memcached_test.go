package memcached

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
	"reflect"
	"testing"

	"github.com/bhojpur/cms/pkg/cache"
)

var client cache.CacheStoreInterface

func init() {
	client = New(&Config{Hosts: []string{"127.0.0.1:11211"}})
}

func TestPlainText(t *testing.T) {
	if err := client.Set("hello_world", "Hello World"); err != nil {
		t.Errorf("No error should happen when saving plain text into client")
	}

	if value, err := client.Get("hello_world"); err != nil || value != "Hello World" {
		t.Errorf("found value: %v", value)
	}

	if err := client.Set("hello_world", "Hello World2"); err != nil {
		t.Errorf("No error should happen when updating saved value")
	}

	if value, err := client.Get("hello_world"); err != nil || value != "Hello World2" {
		t.Errorf("value should been updated: %v", value)
	}

	if err := client.Delete("hello_world"); err != nil {
		t.Errorf("failed to delete value: %v", err)
	}

	if _, err := client.Get("hello_world"); err == nil {
		t.Errorf("the key should been deleted")
	}
}

func TestUnmarshal(t *testing.T) {
	type result struct {
		Name  string
		Value string
	}

	r1 := result{Name: "result_name_1", Value: "result_value_1"}
	if err := client.Set("unmarshal", r1); err != nil {
		t.Errorf("No error should happen when saving struct into client")
	}

	var r2 result
	if err := client.Unmarshal("unmarshal", &r2); err != nil || !reflect.DeepEqual(r1, r2) {
		t.Errorf("found value: %#v", r2)
	}

	if err := client.Delete("unmarshal"); err != nil {
		t.Errorf("failed to delete value: %v", err)
	}

	if err := client.Unmarshal("unmarshal", &r2); err == nil {
		t.Errorf("the key should been deleted")
	}
}

func TestFetch(t *testing.T) {
	var result int
	var fc = func() interface{} {
		result++
		return result
	}

	if value, err := client.Fetch("fetch", fc); err != nil || value != "1" {
		t.Errorf("Should get result from func if key not found")
	}

	if value, err := client.Fetch("fetch", fc); err != nil || value != "1" {
		t.Errorf("Should lookup result from cache store if key is existing")
	}
}
