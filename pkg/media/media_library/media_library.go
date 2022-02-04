package media_library

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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"path"
	"reflect"
	"strings"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/utils"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/media"
	"github.com/bhojpur/cms/pkg/media/oss"
	orm "github.com/bhojpur/orm/pkg/engine"
)

func init() {
	admin.RegisterViewPath("github.com/bhojpur/cms/pkg/media/media_library/views")
}

type MediaLibraryInterface interface {
	ScanMediaOptions(MediaOption) error
	SetSelectedType(string)
	GetSelectedType() string
	GetMediaOption() MediaOption
}

type MediaLibrary struct {
	orm.Model
	SelectedType string
	File         MediaLibraryStorage `sql:"size:4294967295;" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
}

type MediaOption struct {
	Video        string                       `json:",omitempty"`
	FileName     string                       `json:",omitempty"`
	URL          string                       `json:",omitempty"`
	OriginalURL  string                       `json:",omitempty"`
	CropOptions  map[string]*media.CropOption `json:",omitempty"`
	Sizes        map[string]*media.Size       `json:",omitempty"`
	SelectedType string                       `json:",omitempty"`
	Description  string                       `json:",omitempty"`
	Crop         bool
}

func (mediaLibrary *MediaLibrary) ScanMediaOptions(mediaOption MediaOption) error {
	bytes, err := json.Marshal(mediaOption)
	if err == nil {
		return mediaLibrary.File.Scan(bytes)
	}
	return err
}

func (mediaLibrary *MediaLibrary) GetMediaOption() MediaOption {
	return MediaOption{
		Video:        mediaLibrary.File.Video,
		FileName:     mediaLibrary.File.GetFileName(),
		URL:          mediaLibrary.File.URL(),
		OriginalURL:  mediaLibrary.File.URL("original"),
		CropOptions:  mediaLibrary.File.CropOptions,
		Sizes:        mediaLibrary.File.GetSizes(),
		SelectedType: mediaLibrary.File.SelectedType,
		Description:  mediaLibrary.File.Description,
	}
}

func (mediaLibrary *MediaLibrary) SetSelectedType(typ string) {
	mediaLibrary.SelectedType = typ
}

func (mediaLibrary *MediaLibrary) GetSelectedType() string {
	return mediaLibrary.SelectedType
}

func (MediaLibrary) ConfigureBhojpurResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		res.UseTheme("grid")
		res.UseTheme("media_library")
		res.IndexAttrs("File")
	}
}

type MediaLibraryStorage struct {
	oss.OSS
	Sizes        map[string]*media.Size `json:",omitempty"`
	Video        string
	SelectedType string
	Description  string
}

func (mediaLibraryStorage MediaLibraryStorage) GetSizes() map[string]*media.Size {
	if len(mediaLibraryStorage.Sizes) == 0 && !(mediaLibraryStorage.GetFileHeader() != nil || mediaLibraryStorage.Crop) {
		return map[string]*media.Size{}
	}

	var sizes = map[string]*media.Size{
		"@bhojpur_preview": &media.Size{Width: 200, Height: 200},
	}

	for key, value := range mediaLibraryStorage.Sizes {
		sizes[key] = value
	}
	return sizes
}

func (mediaLibraryStorage *MediaLibraryStorage) Scan(data interface{}) (err error) {
	switch values := data.(type) {
	case []byte:
		if mediaLibraryStorage.Sizes == nil {
			mediaLibraryStorage.Sizes = map[string]*media.Size{}
		}
		// cropOptions := mediaLibraryStorage.CropOptions
		sizeOptions := mediaLibraryStorage.Sizes

		if string(values) != "" {
			mediaLibraryStorage.Base.Scan(values)
			if err = json.Unmarshal(values, mediaLibraryStorage); err == nil {
				if mediaLibraryStorage.CropOptions == nil {
					mediaLibraryStorage.CropOptions = map[string]*media.CropOption{}
				}

				// for key, value := range cropOptions {
				// 	if _, ok := mediaLibraryStorage.CropOptions[key]; !ok {
				// 		mediaLibraryStorage.CropOptions[key] = value
				// 	}
				// }

				for key, value := range sizeOptions {
					if key != "original" {
						if _, ok := mediaLibraryStorage.Sizes[key]; !ok {
							mediaLibraryStorage.Sizes[key] = value
						}
					}
				}

				for key, value := range mediaLibraryStorage.CropOptions {
					if _, ok := mediaLibraryStorage.Sizes[key]; !ok {
						mediaLibraryStorage.Sizes[key] = &media.Size{Width: value.Width, Height: value.Height}
					}
				}
			}
		}
	case string:
		err = mediaLibraryStorage.Scan([]byte(values))
	case []string:
		for _, str := range values {
			if err = mediaLibraryStorage.Scan(str); err != nil {
				return err
			}
		}
	default:
		return mediaLibraryStorage.Base.Scan(data)
	}
	return nil
}

func (mediaLibraryStorage MediaLibraryStorage) Value() (driver.Value, error) {
	results, err := json.Marshal(mediaLibraryStorage)
	return string(results), err
}

func (mediaLibraryStorage MediaLibraryStorage) ConfigureBhojpurMeta(metaor resource.Metaor) {
	if meta, ok := metaor.(*admin.Meta); ok {
		meta.Type = "media_library"
		meta.SetFormattedValuer(func(record interface{}, context *appsvr.Context) interface{} {
			return meta.GetValuer()(record, context)
		})
	}
}

type MediaBox struct {
	Values string `json:"-" orm:"size:4294967295;"`
	Files  []File `json:",omitempty"`
}

func (mediaBox MediaBox) URL(styles ...string) string {
	for _, file := range mediaBox.Files {
		return file.URL(styles...)
	}
	return ""
}

func (mediaBox *MediaBox) Scan(data interface{}) (err error) {
	switch values := data.(type) {
	case []byte:
		if mediaBox.Values = string(values); mediaBox.Values != "" {
			return json.Unmarshal(values, &mediaBox.Files)
		}
	case string:
		return mediaBox.Scan([]byte(values))
	case []string:
		for _, str := range values {
			if err := mediaBox.Scan(str); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mediaBox MediaBox) Value() (driver.Value, error) {
	if len(mediaBox.Files) > 0 {
		return json.Marshal(mediaBox.Files)
	}
	return mediaBox.Values, nil
}

func (mediaBox MediaBox) ConfigureBhojpurMeta(metaor resource.Metaor) {
	if meta, ok := metaor.(*admin.Meta); ok {
		if meta.Config == nil {
			meta.Config = &MediaBoxConfig{}
		}

		if meta.FormattedValuer == nil {
			meta.FormattedValuer = func(record interface{}, context *appsvr.Context) interface{} {
				if mediaBox, ok := meta.GetValuer()(record, context).(interface {
					URL(styles ...string) string
				}); ok {
					return mediaBox.URL()
				}
				return ""
			}
			meta.SetFormattedValuer(meta.FormattedValuer)
		}

		if config, ok := meta.Config.(*MediaBoxConfig); ok {
			Admin := meta.GetBaseResource().(*admin.Resource).GetAdmin()
			if config.RemoteDataResource == nil {
				mediaLibraryResource := Admin.GetResource("MediaLibrary")
				if mediaLibraryResource == nil {
					mediaLibraryResource = Admin.NewResource(&MediaLibrary{})
				}
				config.RemoteDataResource = mediaLibraryResource
			}

			if _, ok := config.RemoteDataResource.Value.(MediaLibraryInterface); !ok {
				utils.ExitWithMsg("%v havn't implement MediaLibraryInterface, please fix that.", reflect.TypeOf(config.RemoteDataResource.Value))
			}

			config.RemoteDataResource.Meta(&admin.Meta{
				Name: "MediaOption",
				Type: "hidden",
				Setter: func(record interface{}, metaValue *resource.MetaValue, context *appsvr.Context) {
					if mediaLibrary, ok := record.(MediaLibraryInterface); ok {
						var mediaOption MediaOption
						if err := json.Unmarshal([]byte(utils.ToString(metaValue.Value)), &mediaOption); err == nil {
							mediaOption.FileName = ""
							mediaOption.URL = ""
							mediaOption.OriginalURL = ""
							mediaLibrary.ScanMediaOptions(mediaOption)
						}
					}
				},
				Valuer: func(record interface{}, context *appsvr.Context) interface{} {
					if mediaLibrary, ok := record.(MediaLibraryInterface); ok {
						if value, err := json.Marshal(mediaLibrary.GetMediaOption()); err == nil {
							return string(value)
						}
					}
					return ""
				},
			})

			config.RemoteDataResource.Meta(&admin.Meta{
				Name: "SelectedType",
				Type: "hidden",
				Valuer: func(record interface{}, context *appsvr.Context) interface{} {
					if mediaLibrary, ok := record.(MediaLibraryInterface); ok {
						return mediaLibrary.GetSelectedType()
					}
					return ""
				},
			})

			config.RemoteDataResource.AddProcessor(&resource.Processor{
				Name: "media-library-processor",
				Handler: func(record interface{}, metaValues *resource.MetaValues, context *appsvr.Context) error {
					if mediaLibrary, ok := record.(MediaLibraryInterface); ok {
						var filename string
						var mediaOption MediaOption

						for _, metaValue := range metaValues.Values {
							if fileHeaders, ok := metaValue.Value.([]*multipart.FileHeader); ok {
								for _, fileHeader := range fileHeaders {
									filename = fileHeader.Filename
								}
							}
						}

						if metaValue := metaValues.Get("MediaOption"); metaValue != nil {
							mediaOptionStr := utils.ToString(metaValue.Value)
							json.Unmarshal([]byte(mediaOptionStr), &mediaOption)
						}

						if mediaOption.SelectedType == "video_link" {
							mediaLibrary.SetSelectedType("video_link")
						} else if filename != "" {
							if media.IsImageFormat(filename) {
								mediaLibrary.SetSelectedType("image")
							} else if media.IsVideoFormat(filename) {
								mediaLibrary.SetSelectedType("video")
							} else {
								mediaLibrary.SetSelectedType("file")
							}
						}
					}
					return nil
				},
			})

			config.RemoteDataResource.UseTheme("grid")
			config.RemoteDataResource.UseTheme("media_library")
			if config.RemoteDataResource.Config.PageCount == 0 {
				config.RemoteDataResource.Config.PageCount = admin.PaginationPageCount / 2 * 3
			}

			config.RemoteDataResource.OverrideIndexAttrs(func() {
				config.RemoteDataResource.IndexAttrs(config.RemoteDataResource.IndexAttrs(), "-MediaOption")
			})

			config.RemoteDataResource.OverrideNewAttrs(func() {
				config.RemoteDataResource.NewAttrs(config.RemoteDataResource.NewAttrs(), "MediaOption")
			})

			config.RemoteDataResource.OverrideEditAttrs(func() {
				config.RemoteDataResource.EditAttrs(config.RemoteDataResource.EditAttrs(), "MediaOption")
			})

			config.SelectManyConfig.RemoteDataResource = config.RemoteDataResource
			config.SelectManyConfig.ConfigureBhojpurMeta(meta)
		}

		meta.Type = "media_box"
	}
}

type File struct {
	ID          json.Number
	Url         string
	VideoLink   string
	FileName    string
	Description string
}

// IsImage return if it is an image
func (f File) IsImage() bool {
	return media.IsImageFormat(f.Url)
}

func (f File) IsVideo() bool {
	return media.IsVideoFormat(f.Url)
}

func (f File) IsSVG() bool {
	return media.IsSVGFormat(f.Url)
}

func (file File) URL(styles ...string) string {
	if file.Url != "" && len(styles) > 0 {
		ext := path.Ext(file.Url)
		return fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file.Url, ext), styles[0], ext)
	}
	return file.Url
}

func (mediaBox MediaBox) Crop(res *admin.Resource, db *orm.DB, mediaOption MediaOption) (err error) {
	for _, file := range mediaBox.Files {
		context := &appsvr.Context{ResourceID: string(file.ID), DB: db}
		record := res.NewStruct()
		if err = res.CallFindOne(record, nil, context); err == nil {
			if mediaLibrary, ok := record.(MediaLibraryInterface); ok {
				mediaOption.Crop = true
				if err = mediaLibrary.ScanMediaOptions(mediaOption); err == nil {
					err = res.CallSave(record, context)
				}
			} else {
				err = errors.New("invalid media library resource")
			}
		}
		if err != nil {
			return
		}
	}
	return
}

const (
	ALLOW_TYPE_FILE  = "file"
	ALLOW_TYPE_IMAGE = "image"
	ALLOW_TYPE_VIDEO = "video"
)

// MediaBoxConfig configure MediaBox metas
type MediaBoxConfig struct {
	RemoteDataResource *admin.Resource
	Sizes              map[string]*media.Size
	Max                uint
	AllowType          string
	admin.SelectManyConfig
}

func (*MediaBoxConfig) ConfigureBhojpurMeta(resource.Metaor) {
}

func (*MediaBoxConfig) GetTemplate(context *admin.Context, metaType string) ([]byte, error) {
	return nil, errors.New("not implemented")
}
