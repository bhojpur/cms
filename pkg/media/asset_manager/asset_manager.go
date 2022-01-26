package asset_manager

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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/media/oss"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// AssetManager defined a asset manager that could be used to manage assets in Bhojpur CMS admin
type AssetManager struct {
	orm.Model
	File oss.OSS `media_library:"URL:/system/assets/{{primary_key}}/{{filename_with_hash}}"`
}

// ConfigureBhojpurResource configure locale for Bhojpur CMS Admin
func (*AssetManager) ConfigureBhojpurResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		router := res.GetAdmin().GetRouter()
		router.Post(fmt.Sprintf("/%v/upload", res.ToParam()), func(context *admin.Context) {
			result := AssetManager{}
			result.File.Scan(context.Request.MultipartForm.File["file[]"])
			context.GetDB().Save(&result)
			bytes, _ := json.Marshal(map[string]interface{}{"file-0": map[string]string{"url": result.File.URL(), "id": result.File.GetFileName()}})
			context.Writer.Write(bytes)
		})

		assetURL := regexp.MustCompile(`^/system/assets/(\d+)/`)
		router.Post(fmt.Sprintf("/%v/crop", res.ToParam()), func(context *admin.Context) {
			defer context.Request.Body.Close()
			var (
				err error
				url struct{ URL string }
				buf bytes.Buffer
			)

			io.Copy(&buf, context.Request.Body)
			if err = json.Unmarshal(buf.Bytes(), &url); err == nil {
				if matches := assetURL.FindStringSubmatch(url.URL); len(matches) > 1 {
					result := &AssetManager{}
					if err = context.GetDB().Find(result, matches[1]).Error; err == nil {
						if err = result.File.Scan(buf.Bytes()); err == nil {
							if err = context.GetDB().Save(result).Error; err == nil {
								bytes, _ := json.Marshal(map[string]interface{}{"file-0": map[string]string{"url": result.File.URL(), "id": result.File.GetFileName()}})
								context.Writer.Write(bytes)
								return
							}
						}
					}
				}
			}

			bytes, _ := json.Marshal(map[string]string{"err": err.Error()})
			context.Writer.Write(bytes)
		})
	}
}
