package exchange_actions_test

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
	"io/ioutil"
	"os"
	"testing"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/test/utils"
	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/i18n"
	"github.com/bhojpur/cms/pkg/i18n/backends/database"
	"github.com/bhojpur/cms/pkg/i18n/exchange_actions"
	"github.com/bhojpur/cms/pkg/media"
	"github.com/bhojpur/cms/pkg/media/oss"
	"github.com/bhojpur/cms/pkg/worker"
	orm "github.com/bhojpur/orm/pkg/engine"
	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
)

var db *orm.DB
var Worker *worker.Worker
var I18N *i18n.I18N

func init() {
	db = utils.TestDB()
	reset()
}

func reset() {
	db.DropTable(&database.Translation{})
	database.New(db)
	Admin := admin.New(&appsvr.Config{DB: db})
	Worker = worker.New()
	Admin.AddResource(Worker)
	I18N = i18n.New(database.New(db))
	I18N.SaveTranslation(&i18n.Translation{Key: "bhojpur_admin.title", Value: "title", Locale: "en-US"})
	I18N.SaveTranslation(&i18n.Translation{Key: "bhojpur_admin.subtitle", Value: "subtitle", Locale: "en-US"})
	I18N.SaveTranslation(&i18n.Translation{Key: "bhojpur_admin.description", Value: "description", Locale: "en-US"})
	I18N.SaveTranslation(&i18n.Translation{Key: "header.title", Value: "Header Title", Locale: "en-US"})
	exchange_actions.RegisterExchangeJobs(I18N, Worker)
}

// Test export translations with scope
type testExportWithScopedCase struct {
	Scope            string
	ExpectExportFile string
}

func TestExportTranslations(t *testing.T) {
	reset()
	I18N.SaveTranslation(&i18n.Translation{Key: "header.title", Value: "标题", Locale: "zh-CN"})

	testCases := []*testExportWithScopedCase{
		&testExportWithScopedCase{Scope: "", ExpectExportFile: "export_all.csv"},
		&testExportWithScopedCase{Scope: "All", ExpectExportFile: "export_all.csv"},
		&testExportWithScopedCase{Scope: "Backend", ExpectExportFile: "export_backend.csv"},
		&testExportWithScopedCase{Scope: "Frontend", ExpectExportFile: "export_frontend.csv"},
	}

	for i, testcase := range testCases {
		clearDownloadDir()
		for _, job := range Worker.Jobs {
			if job.Name == "Export Translations" {
				job.Handler(&exchange_actions.ExportTranslationArgument{Scope: testcase.Scope}, job.NewStruct().(worker.BhojpurJobInterface))
				if downloadedFileContent() != loadFixture(testcase.ExpectExportFile) {
					t.Errorf(color.RedString(fmt.Sprintf("\nExchange TestCase #%d: Failure (%s)\n", i+1, "export results are incorrect")))
				} else {
					color.Green(fmt.Sprintf("Export with scope TestCase #%d: Success\n", i+1))
				}
			}
		}
	}
}

// Test import translations
type testImportTranslationsCase struct {
	ImportFileDesc string
	ImportFile     string
	ExpectZhValues map[string]string
}

func TestImportTranslations(t *testing.T) {
	reset()
	testCases := []*testImportTranslationsCase{
		&testImportTranslationsCase{
			ImportFileDesc: "Normal tranlsation file",
			ImportFile:     "import_1.csv",
			ExpectZhValues: map[string]string{"bhojpur_admin.title": "标题", "bhojpur_admin.subtitle": "小标题", "bhojpur_admin.description": "描述", "header.title": "标题"},
		},
		&testImportTranslationsCase{
			ImportFileDesc: "Translation file with missing header.title",
			ImportFile:     "import_2.csv",
			ExpectZhValues: map[string]string{"bhojpur_admin.title": "标题", "bhojpur_admin.subtitle": "小标题", "bhojpur_admin.description": "描述"},
		},
		&testImportTranslationsCase{
			ImportFileDesc: "Translation file with empty column",
			ImportFile:     "import_3.csv",
			ExpectZhValues: map[string]string{"bhojpur_admin.title": "标题", "bhojpur_admin.subtitle": "小标题", "bhojpur_admin.description": "描述", "header.title": "标题"},
		},
	}

	for i, testCase := range testCases {
		for _, job := range Worker.Jobs {
			if job.Name == "Import Translations" {
				job.Handler(&exchange_actions.ImportTranslationArgument{TranslationsFile: oss.OSS{media.Base{Url: "imports/" + testCase.ImportFile}}}, job.NewStruct().(worker.BhojpurJobInterface))
				translations := I18N.LoadTranslations()["zh-CN"]
				if len(translations) == 0 {
					t.Errorf(color.RedString(fmt.Sprintf("\nImport TestCase #%d: Failure (%s)\n", i+1, "Doesn't have Zh translations")))
				}
				for key, translation := range translations {
					if testCase.ExpectZhValues[key] != translation.Value {
						t.Errorf(color.RedString(fmt.Sprintf("\nImport TestCase #%d: Failure (%s)\n", i+1, "Zh translations not match")))
					}
				}
				color.Green(fmt.Sprintf("Import TestCase #%d: Success\n", i+1))
			}
		}
	}
}

// Helper functions
func clearDownloadDir() {
	files, _ := ioutil.ReadDir("./public/downloads")
	for _, f := range files {
		os.Remove("./public/downloads/" + f.Name())
	}
}

func downloadedFileContent() string {
	files, _ := ioutil.ReadDir("./public/downloads")
	for _, f := range files {
		if content, err := ioutil.ReadFile("./public/downloads/" + f.Name()); err == nil {
			return string(content)
		}
	}
	return ""
}

func loadFixture(fileName string) string {
	if content, err := ioutil.ReadFile("./fixtures/" + fileName); err == nil {
		return string(content)
	}
	return ""
}
