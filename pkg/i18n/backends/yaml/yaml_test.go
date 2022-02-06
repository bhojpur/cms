package yaml_test

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
	"net/http"
	"testing"

	"github.com/bhojpur/cms/pkg/i18n"
	"github.com/bhojpur/cms/pkg/i18n/backends/yaml"
)

var values = map[string][][]string{
	"en": {
		{"hello", "Hello"},
		{"user.name", "User Name"},
		{"user.email", "Email"},
	},
	"de": {
		{"hello", "Hallo"},
		{"user.name", "Benutzername"},
		{"user.email", "E-Mail-Adresse"},
	},
	"zh-CN": {
		{"hello", "你好"},
		{"user.name", "用户名"},
		{"user.email", "邮箱"},
	},
}

func checkTranslations(translations []*i18n.Translation) error {
	for locale, results := range values {
		for _, result := range results {
			var found bool
			for _, translation := range translations {
				if (translation.Locale == locale) && (translation.Key == result[0]) && (translation.Value == result[1]) {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("failed to found translation %v for %v", result[0], locale)
			}
		}
	}
	return nil
}

func TestLoadTranslations(t *testing.T) {
	backend := yaml.New("tests", "tests/subdir")
	if err := checkTranslations(backend.LoadTranslations()); err != nil {
		t.Fatal(err)
	}
}

func TestLoadTranslationsFilesystem(t *testing.T) {
	backend := yaml.NewWithFilesystem(http.Dir("./tests"))
	if err := checkTranslations(backend.LoadTranslations()); err != nil {
		t.Fatal(err)
	}
}

func TestLoadTranslationsWalk(t *testing.T) {
	backend := yaml.NewWithWalk("tests")
	if err := checkTranslations(backend.LoadTranslations()); err != nil {
		t.Fatal(err)
	}
}

var benchmarkResult error

func BenchmarkLoadTranslations(b *testing.B) {
	var backend i18n.Backend
	var err error
	for i := 0; i < b.N; i++ {
		backend = yaml.New("tests", "tests/subdir")
		if err = checkTranslations(backend.LoadTranslations()); err != nil {
			b.Fatal(err)
		}
	}
	benchmarkResult = err
}

var benchmarkResult2 error

func BenchmarkLoadTranslationsWalk(b *testing.B) {
	var backend i18n.Backend
	var err error
	for i := 0; i < b.N; i++ {
		backend = yaml.NewWithWalk("tests")
		if err = checkTranslations(backend.LoadTranslations()); err != nil {
			b.Fatal(err)
		}
	}
	benchmarkResult2 = err
}

var benchmarkResult3 error

func BenchmarkLoadTranslationsFilesystem(b *testing.B) {
	var backend i18n.Backend
	var err error
	for i := 0; i < b.N; i++ {
		backend = yaml.NewWithFilesystem(http.Dir("./tests"))
		if err = checkTranslations(backend.LoadTranslations()); err != nil {
			b.Fatal(err)
		}
	}
	benchmarkResult3 = err
}
