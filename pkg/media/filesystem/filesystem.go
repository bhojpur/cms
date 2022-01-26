package filesystem

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
	"io"
	"os"
	"path/filepath"

	"github.com/bhojpur/cms/pkg/media"
)

var _ media.Media = &FileSystem{}

// FileSystem defined a media library storage using file system
type FileSystem struct {
	media.Base
}

// GetFullPath return full file path from a relative file path
func (f FileSystem) GetFullPath(url string, option *media.Option) (path string, err error) {
	if option != nil && option.Get("path") != "" {
		path = filepath.Join(option.Get("path"), url)
	} else {
		path = filepath.Join("./public", url)
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}

	return
}

// Store save reader's context with name
func (f FileSystem) Store(name string, option *media.Option, reader io.Reader) (err error) {
	if fullpath, err := f.GetFullPath(name, option); err == nil {
		if dst, err := os.Create(fullpath); err == nil {
			_, err = io.Copy(dst, reader)
		}
	}
	return err
}

// Retrieve retrieve file content with url
func (f FileSystem) Retrieve(url string) (media.FileInterface, error) {
	if fullpath, err := f.GetFullPath(url, nil); err == nil {
		return os.Open(fullpath)
	}
	return nil, os.ErrNotExist
}
