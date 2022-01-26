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
	"testing"
)

func TestGetDepVersionFromMod(t *testing.T) {
	goModDeps = []string{
		"github.com/bhojpur/orm v0.0.1",
		"github.com/bhojpur/application v0.0.2",
	}
	cases := []struct {
		view string
		want string
	}{
		{view: "github.com/bhojpur/cms/pkg/l10n/views", want: "pkg/mod/github.com/bhojpur/cms/pkg/l10n@v0.0.1/views"},
		{view: "github.com/bhojpur/cms/pkg/admin/views", want: "pkg/mod/github.com/bhojpur/cms/pkg/admin@v0.0.1/views"},
		{view: "github.com/bhojpur/cms/pkg/publish2/views", want: "pkg/mod/github.com/bhojpur/cms/pkg/publish2@v0.0.1/views"},
		{view: "github.com/bhojpur/cms/pkg/media/media_library/views", want: "pkg/mod/github.com/bhojpur/cms/pkg/media@v0.0.1/media_library/views"},
		{view: "github.com/bhojpur/cms/pkg/i18n/exchange_actions/views", want: "pkg/mod/github.com/bhojpur/cms/pkg/i18n@v0.0.1/exchange_actions/views"},
		{view: "no/unknown/nonexistent", want: "no/unknown/nonexistent"},
	}
	for _, v := range cases {
		if got := getDepVersionFromMod(v.view); v.want != got {
			t.Errorf("GetDepVersionFromMod-viewpath: %v, want: %v, got: %v", v.view, v.want, got)
		} else {
			fmt.Println(got)
		}
	}
}
