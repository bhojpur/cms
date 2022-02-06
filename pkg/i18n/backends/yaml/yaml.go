package yaml

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
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bhojpur/cms/pkg/i18n"
	"gopkg.in/yaml.v2"
)

var _ i18n.Backend = &Backend{}

// New new YAML backend for I18n
func New(paths ...string) *Backend {
	backend := &Backend{}

	var files []string
	for _, p := range paths {
		if file, err := os.Open(p); err == nil {
			defer file.Close()
			if fileInfo, err := file.Stat(); err == nil {
				if fileInfo.IsDir() {
					yamlFiles, _ := filepath.Glob(filepath.Join(p, "*.yaml"))
					files = append(files, yamlFiles...)

					ymlFiles, _ := filepath.Glob(filepath.Join(p, "*.yml"))
					files = append(files, ymlFiles...)
				} else if fileInfo.Mode().IsRegular() {
					files = append(files, p)
				}
			}
		}
	}
	for _, file := range files {
		if content, err := ioutil.ReadFile(file); err == nil {
			backend.contents = append(backend.contents, content)
		}
	}
	return backend
}

// NewWithWalk has the same functionality as New but uses filepath.Walk to find all the translation files recursively.
func NewWithWalk(paths ...string) i18n.Backend {
	backend := &Backend{}

	var files []string
	for _, p := range paths {
		filepath.Walk(p, func(path string, fileInfo os.FileInfo, err error) error {
			if isYamlFile(fileInfo) {
				files = append(files, path)
			}
			return nil
		})
	}
	for _, file := range files {
		if content, err := ioutil.ReadFile(file); err == nil {
			backend.contents = append(backend.contents, content)
		}
	}

	return backend
}

func isYamlFile(fileInfo os.FileInfo) bool {
	if fileInfo == nil {
		return false
	}
	return fileInfo.Mode().IsRegular() && (strings.HasSuffix(fileInfo.Name(), ".yml") || strings.HasSuffix(fileInfo.Name(), ".yaml"))
}

func walkFilesystem(fs http.FileSystem, entry http.File, prefix string) [][]byte {
	var (
		contents [][]byte
		err      error
		isRoot   bool
	)
	if entry == nil {
		if entry, err = fs.Open("/"); err != nil {
			return nil
		}
		isRoot = true
		defer entry.Close()
	}
	fileInfo, err := entry.Stat()
	if err != nil {
		return nil
	}
	if !isRoot {
		prefix = prefix + fileInfo.Name() + "/"
	}
	if fileInfo.IsDir() {
		if entries, err := entry.Readdir(-1); err == nil {
			for _, e := range entries {
				if file, err := fs.Open(prefix + e.Name()); err == nil {
					defer file.Close()
					contents = append(contents, walkFilesystem(fs, file, prefix)...)
				}
			}
		}
	} else if isYamlFile(fileInfo) {
		if content, err := ioutil.ReadAll(entry); err == nil {
			contents = append(contents, content)
		}
	}
	return contents
}

// NewWithFilesystem initializes a backend that reads translation files from an http.FileSystem.
func NewWithFilesystem(fss ...http.FileSystem) i18n.Backend {
	backend := &Backend{}

	for _, fs := range fss {
		backend.contents = append(backend.contents, walkFilesystem(fs, nil, "/")...)
	}
	return backend
}

// Backend YAML backend
type Backend struct {
	contents [][]byte
}

func loadTranslationsFromYaml(locale string, value interface{}, scopes []string) (translations []*i18n.Translation) {
	switch v := value.(type) {
	case yaml.MapSlice:
		for _, s := range v {
			results := loadTranslationsFromYaml(locale, s.Value, append(scopes, fmt.Sprint(s.Key)))
			translations = append(translations, results...)
		}
	default:
		var translation = &i18n.Translation{
			Locale: locale,
			Key:    strings.Join(scopes, "."),
			Value:  fmt.Sprint(v),
		}
		translations = append(translations, translation)
	}
	return
}

// LoadYAMLContent load YAML content
func (backend *Backend) LoadYAMLContent(content []byte) (translations []*i18n.Translation, err error) {
	var slice yaml.MapSlice

	if err = yaml.Unmarshal(content, &slice); err == nil {
		for _, item := range slice {
			translations = append(translations, loadTranslationsFromYaml(item.Key.(string) /* locale */, item.Value, []string{})...)
		}
	}

	return translations, err
}

// LoadTranslations load translations from YAML backend
func (backend *Backend) LoadTranslations() (translations []*i18n.Translation) {
	for _, content := range backend.contents {
		if results, err := backend.LoadYAMLContent(content); err == nil {
			translations = append(translations, results...)
		} else {
			panic(err)
		}
	}
	return translations
}

// SaveTranslation save translation into YAML backend, not implemented
func (backend *Backend) SaveTranslation(t *i18n.Translation) error {
	return errors.New("not implemented")
}

// FindTranslation find translation from  backend
func (backend *Backend) FindTranslation(t *i18n.Translation) (translation i18n.Translation) {
	return translation //not implemented
}

// DeleteTranslation delete translation into YAML backend, not implemented
func (backend *Backend) DeleteTranslation(t *i18n.Translation) error {
	return errors.New("not implemented")
}
