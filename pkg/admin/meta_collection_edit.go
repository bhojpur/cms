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
	"errors"

	"github.com/bhojpur/application/pkg/resource"
)

// CollectionEditConfig meta configuration used for collection edit
type CollectionEditConfig struct {
	Template string
	Max      uint
	metaConfig
}

// GetTemplate get template for collection edit
func (collectionEditConfig CollectionEditConfig) GetTemplate(context *Context, metaType string) ([]byte, error) {
	if metaType == "form" && collectionEditConfig.Template != "" {
		return context.Asset(collectionEditConfig.Template)
	}
	return nil, errors.New("not implemented")
}

// ConfigureBhojpurMeta configure collection edit meta
func (collectionEditConfig *CollectionEditConfig) ConfigureBhojpurMeta(metaor resource.Metaor) {
	if meta, ok := metaor.(*Meta); ok {
		if meta.Permission != nil || meta.Resource.Permission != nil {
			meta.Permission = meta.Permission.Concat(meta.Resource.Permission)
		}
	}
}
