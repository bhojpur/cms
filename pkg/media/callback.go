package media

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
	"reflect"

	"github.com/bhojpur/cms/pkg/serializable_meta"
	orm "github.com/bhojpur/orm/pkg/engine"
)

var (
	// set MediaLibraryURL to change the default url /system/{{class}}/{{primary_key}}/{{column}}.{{extension}}
	MediaLibraryURL = ""
)

func cropField(field *orm.Field, scope *orm.Scope) (cropped bool) {
	if field.Field.CanAddr() {
		// TODO Handle scanner
		if media, ok := field.Field.Addr().Interface().(Media); ok && !media.Cropped() {
			option := parseTagOption(field.Tag.Get("media_library"))
			if MediaLibraryURL != "" {
				option.Set("url", MediaLibraryURL)
			}
			if media.GetFileHeader() != nil || media.NeedCrop() {
				var mediaFile FileInterface
				var err error
				if fileHeader := media.GetFileHeader(); fileHeader != nil {
					mediaFile, err = media.GetFileHeader().Open()
				} else {
					mediaFile, err = media.Retrieve(media.URL("original"))
				}

				if err != nil {
					scope.Err(err)
					return false
				}

				media.Cropped(true)

				if url := media.GetURL(option, scope, field, media); url == "" {
					scope.Err(errors.New("invalid URL"))
				} else {
					result, _ := json.Marshal(map[string]string{"Url": url})
					media.Scan(string(result))
				}

				if mediaFile != nil {
					defer mediaFile.Close()
					var handled = false
					for _, handler := range mediaHandlers {
						if handler.CouldHandle(media) {
							mediaFile.Seek(0, 0)
							if scope.Err(handler.Handle(media, mediaFile, option)) == nil {
								handled = true
							}
						}
					}

					// Save File
					if !handled {
						scope.Err(media.Store(media.URL(), option, mediaFile))
					}
				}
				return true
			}
		}
	}
	return false
}

func saveAndCropImage(isCreate bool) func(scope *orm.Scope) {
	return func(scope *orm.Scope) {
		if !scope.HasError() {
			var updateColumns = map[string]interface{}{}

			// Handle SerializableMeta
			if value, ok := scope.Value.(serializable_meta.SerializableMetaInterface); ok {
				var (
					isCropped        bool
					handleNestedCrop func(record interface{})
				)

				handleNestedCrop = func(record interface{}) {
					newScope := scope.New(record)
					for _, field := range newScope.Fields() {
						if cropField(field, scope) {
							isCropped = true
							continue
						}

						if reflect.Indirect(field.Field).Kind() == reflect.Struct {
							handleNestedCrop(field.Field.Addr().Interface())
						}

						if reflect.Indirect(field.Field).Kind() == reflect.Slice {
							for i := 0; i < reflect.Indirect(field.Field).Len(); i++ {
								handleNestedCrop(reflect.Indirect(field.Field).Index(i).Addr().Interface())
							}
						}
					}
				}

				record := value.GetSerializableArgument(value)
				handleNestedCrop(record)
				if isCreate && isCropped {
					updateColumns["value"], _ = json.Marshal(record)
				}
			}

			// Handle Normal Field
			for _, field := range scope.Fields() {
				if cropField(field, scope) && isCreate {
					updateColumns[field.DBName] = field.Field.Interface()
				}
			}

			if !scope.HasError() && len(updateColumns) != 0 {
				scope.Err(scope.NewDB().Model(scope.Value).UpdateColumns(updateColumns).Error)
			}
		}
	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *orm.DB) {
	if db.Callback().Create().Get("media:save_and_crop") == nil {
		db.Callback().Create().After("gorm:after_create").Register("media:save_and_crop", saveAndCropImage(true))
	}
	if db.Callback().Update().Get("media:save_and_crop") == nil {
		db.Callback().Update().Before("gorm:before_update").Register("media:save_and_crop", saveAndCropImage(false))
	}
}
