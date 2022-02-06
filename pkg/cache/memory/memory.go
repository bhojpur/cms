package memory

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
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type Memory struct {
	values map[string][]byte
	mutex  *sync.RWMutex
}

func New() *Memory {
	return &Memory{values: map[string][]byte{}, mutex: &sync.RWMutex{}}
}

func (memory *Memory) Get(key string) (string, error) {
	memory.mutex.RLock()
	defer memory.mutex.RUnlock()

	if value, ok := memory.values[key]; ok {
		return string(value), nil
	}
	return "", ErrNotFound
}

func (memory *Memory) Unmarshal(key string, object interface{}) error {
	memory.mutex.RLock()
	defer memory.mutex.RUnlock()

	if value, ok := memory.values[key]; ok {
		return json.Unmarshal(value, object)
	}
	return ErrNotFound
}

func convertToBytes(value interface{}) []byte {
	switch result := value.(type) {
	case string:
		return []byte(result)
	case []byte:
		return result
	default:
		bytes, _ := json.Marshal(value)
		return bytes
	}
}

func (memory *Memory) Set(key string, value interface{}) error {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()

	memory.values[key] = convertToBytes(value)
	return nil
}

func (memory *Memory) Fetch(key string, fc func() interface{}) (string, error) {
	if str, err := memory.Get(key); err == nil {
		return str, nil
	}
	results := convertToBytes(fc())
	return string(results), memory.Set(key, results)
}

func (memory *Memory) Delete(key string) error {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()

	delete(memory.values, key)
	return nil
}
