package database_test

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
	"testing"

	"github.com/bhojpur/application/test/utils"
	"github.com/bhojpur/cms/pkg/i18n"
	"github.com/bhojpur/cms/pkg/i18n/backends/database"
	orm "github.com/bhojpur/orm/pkg/engine"
)

var db *orm.DB
var backend i18n.Backend

func init() {
	db = utils.TestDB()
	db.DropTable(&database.Translation{})
	backend = database.New(db)
}

func TestTranslations(t *testing.T) {
	translation := i18n.Translation{Key: "hello_world", Value: "Hello World", Locale: "hi-IN"}

	backend.SaveTranslation(&translation)
	if len(backend.LoadTranslations()) != 1 {
		t.Errorf("should has only one translation")
	}

	backend.DeleteTranslation(&translation)
	if len(backend.LoadTranslations()) != 0 {
		t.Errorf("should has none translation")
	}

	longText := "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	backend.SaveTranslation(&i18n.Translation{Key: longText + "1", Value: longText, Locale: "hi-IN"})
	backend.SaveTranslation(&i18n.Translation{Key: longText + "2", Value: longText, Locale: "hi-IN"})

	if len(backend.LoadTranslations()) != 2 {
		t.Errorf("should has two translations")
	}

	backend.DeleteTranslation(&i18n.Translation{Key: longText + "1", Value: longText, Locale: "hi-IN"})
	if len(backend.LoadTranslations()) != 1 {
		t.Errorf("should has one translation left")
	}
}
