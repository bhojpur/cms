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
	"fmt"
	"time"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/utils"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// DatetimeConfig meta configuration used for datetime
type DatetimeConfig struct {
	MinTime  *time.Time
	MaxTime  *time.Time
	ShowTime bool
}

// ConfigureBhojpurMeta configure datetime meta
func (datetimeConfig *DatetimeConfig) ConfigureBhojpurMeta(metaor resource.Metaor) {
	if meta, ok := metaor.(*Meta); ok {
		timeFormat := "2006-01-02"
		if meta.Type == "datetime" {
			datetimeConfig.ShowTime = true
		}

		if meta.Type == "" {
			meta.Type = "datetime"
		}

		if datetimeConfig.ShowTime {
			timeFormat = "2006-01-02 15:04"
		}

		if meta.FormattedValuer == nil {
			meta.SetFormattedValuer(func(value interface{}, context *appsvr.Context) interface{} {
				switch date := meta.GetValuer()(value, context).(type) {
				case *time.Time:
					if date == nil {
						return ""
					}
					if date.IsZero() {
						return ""
					}
					return utils.FormatTime(*date, timeFormat, context)
				case time.Time:
					if date.IsZero() {
						return ""
					}
					return utils.FormatTime(date, timeFormat, context)
				default:
					return date
				}
			})
		}
	}
}

// ConfigureBhojpurAdminFilter configure admin filter for datetime
func (datetimeConfig *DatetimeConfig) ConfigureBhojpurAdminFilter(filter *Filter) {
	if filter.Handler == nil {
		if dbName := filter.Resource.GetMeta(filter.Name).DBName(); dbName != "" {
			filter.Handler = func(tx *orm.DB, filterArgument *FilterArgument) *orm.DB {
				if metaValue := filterArgument.Value.Get("Start"); metaValue != nil {
					if start := utils.ToString(metaValue.Value); start != "" {
						tx = tx.Where(fmt.Sprintf("%v > ?", dbName), start)
					}
				}

				if metaValue := filterArgument.Value.Get("End"); metaValue != nil {
					if end := utils.ToString(metaValue.Value); end != "" {
						tx = tx.Where(fmt.Sprintf("%v < ?", dbName), end)
					}
				}

				return tx
			}
		}
	}
}
