package redis

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

	redis "gopkg.in/redis.v3"
)

// Redis provides a cache backed by a Redis server.
type Redis struct {
	Config *redis.Options
	Client *redis.Client
}

// New returns an initialized Redis cache object.
func New(config *redis.Options) *Redis {
	client := redis.NewClient(config)
	return &Redis{Config: config, Client: client}
}

// Get returns the value saved under a given key.
func (r *Redis) Get(key string) (string, error) {
	return r.Client.Get(key).Result()
}

// Unmarshal retrieves a value from the Redis server and unmarshals
// it into the passed object.
func (r *Redis) Unmarshal(key string, object interface{}) error {
	value, err := r.Get(key)
	if err == nil {
		err = json.Unmarshal([]byte(value), object)
	}
	return err
}

// Set saves an arbitrary value under a specific key.
func (r *Redis) Set(key string, value interface{}) error {
	return r.Client.Set(key, convertToBytes(value), 0).Err()
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

// Fetch returns the value for the key if it exists or sets and returns the value via the passed function.
func (r *Redis) Fetch(key string, fc func() interface{}) (string, error) {
	if str, err := r.Get(key); err == nil {
		return str, nil
	}
	results := convertToBytes(fc())
	return string(results), r.Set(key, results)
}

// Delete removes a specific key and its value from the Redis server.
func (r *Redis) Delete(key string) error {
	return r.Client.Del(key).Err()
}
