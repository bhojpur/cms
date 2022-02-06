package i18n

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
	"testing"
)

type backend struct{}

func (b *backend) LoadTranslations() (translations []*Translation) { return translations }
func (b *backend) SaveTranslation(t *Translation) error            { return nil }
func (b *backend) DeleteTranslation(t *Translation) error          { return nil }

const BIGNUM = 10000

// run TestConcurrent* tests with -race flag would be better

func TestConcurrentReadWrite(t *testing.T) {
	i18n := New(&backend{})
	go func() {
		for i := 0; i < BIGNUM; i++ {
			i18n.AddTranslation(&Translation{Key: fmt.Sprintf("xx-%d", i), Locale: "xx", Value: fmt.Sprint(i)})
		}
	}()
	for i := 0; i < BIGNUM; i++ {
		i18n.T("xx", fmt.Sprintf("xx-%d", i))
	}
}

func TestConcurrentDeleteWrite(t *testing.T) {
	i18n := New(&backend{})
	go func() {
		for i := 0; i < BIGNUM; i++ {
			i18n.AddTranslation(&Translation{Key: fmt.Sprintf("xx-%d", i), Locale: "xx", Value: fmt.Sprint(i)})
		}
	}()
	for i := 0; i < BIGNUM; i++ {
		i18n.DeleteTranslation(&Translation{Key: fmt.Sprintf("xx-%d", i), Locale: "xx", Value: fmt.Sprint(i)})
	}
}

func TestFallbackLocale(t *testing.T) {
	i18n := New(&backend{})
	i18n.AddTranslation(&Translation{Key: "hello-world", Locale: "en-AU", Value: "Hello World"})

	if i18n.Fallbacks("en-AU").T("en-UK", "hello-world") != "Hello World" {
		t.Errorf("Should fallback en-UK to en-US")
	}

	if i18n.T("en-DE", "hello-world") != "hello-world" {
		t.Errorf("Haven't setup any fallback")
	}
}
