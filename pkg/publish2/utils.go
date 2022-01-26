package publish2

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

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/utils"
	orm "github.com/bhojpur/orm/pkg/engine"
)

func getPublishScheduleTime(context *appsvr.Context) string {
	if values, ok := context.Request.URL.Query()["publish_scheduled_time"]; ok {
		if len(values) > 0 && values[0] != "" {
			return values[0]
		}
	} else if cookie, err := context.Request.Cookie("publish2_publish_scheduled_time"); err == nil {
		return cookie.Value
	}
	return ""
}

func requestingPublishDraftContent(context *appsvr.Context) bool {
	if values, ok := context.Request.URL.Query()["publish_draft_content"]; ok {
		if len(values) > 0 && values[0] != "" {
			return true
		}
	} else if cookie, err := context.Request.Cookie("pubilsh2_publish_draft_content"); err == nil && cookie.Value == "true" {
		return true
	}
	return false
}

func PreviewByDB(tx *orm.DB, context *appsvr.Context) *orm.DB {
	scheduledTime := getPublishScheduleTime(context)
	draftContent := requestingPublishDraftContent(context)

	utils.SetCookie(http.Cookie{Name: "publish2_publish_scheduled_time", Value: scheduledTime}, context)
	utils.SetCookie(http.Cookie{Name: "pubilsh2_publish_draft_content", Value: fmt.Sprint(draftContent)}, context)

	if scheduledTime != "" {
		if t, err := utils.ParseTime(scheduledTime, context); err == nil {
			tx = tx.Set(ScheduledTime, t)
		}
	}

	if draftContent {
		tx = tx.Set(VisibleMode, ModeOff)
	}

	return tx
}
