package assetfs

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
	"os"
	"path/filepath"
	"strings"
)

// AssetFileSystem AssetFS based on FileSystem
type AssetFileSystem struct {
	paths        []string
	nameSpacedFS map[string]Interface
}

// RegisterPath register view paths
func (fs *AssetFileSystem) RegisterPath(pth string) error {
	if _, err := os.Stat(pth); !os.IsNotExist(err) {
		var existing bool
		for _, p := range fs.paths {
			if p == pth {
				existing = true
				break
			}
		}
		if !existing {
			fs.paths = append(fs.paths, pth)
		}
		return nil
	}
	return errors.New("not found")
}

// PrependPath prepend path to view paths
func (fs *AssetFileSystem) PrependPath(pth string) error {
	if _, err := os.Stat(pth); !os.IsNotExist(err) {
		var existing bool
		for _, p := range fs.paths {
			if p == pth {
				existing = true
				break
			}
		}
		if !existing {
			fs.paths = append([]string{pth}, fs.paths...)
		}
		return nil
	}
	return errors.New("not found")
}

// Asset get content with name from assetfs
func (fs *AssetFileSystem) Asset(name string) ([]byte, error) {
	for _, pth := range fs.paths {
		if _, err := os.Stat(filepath.Join(pth, name)); err == nil {
			return ioutil.ReadFile(filepath.Join(pth, name))
		}
	}
	return []byte{}, fmt.Errorf("%v not found", name)
}

// Glob list matched files from assetfs
func (fs *AssetFileSystem) Glob(pattern string) (matches []string, err error) {
	for _, pth := range fs.paths {
		if results, err := filepath.Glob(filepath.Join(pth, pattern)); err == nil {
			for _, result := range results {
				matches = append(matches, strings.TrimPrefix(result, pth))
			}
		}
	}
	return
}

// Compile compile assetfs
func (fs *AssetFileSystem) Compile() error {
	return nil
}

// NameSpace return namespaced filesystem
func (fs *AssetFileSystem) NameSpace(nameSpace string) Interface {
	if fs.nameSpacedFS == nil {
		fs.nameSpacedFS = map[string]Interface{}
	}
	fs.nameSpacedFS[nameSpace] = &AssetFileSystem{}
	return fs.nameSpacedFS[nameSpace]
}
