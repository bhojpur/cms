package exchange_actions

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
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/bhojpur/cms/pkg/admin"
	"github.com/bhojpur/cms/pkg/i18n"
	"github.com/bhojpur/cms/pkg/media/oss"
	"github.com/bhojpur/cms/pkg/worker"
)

type ExportTranslationArgument struct {
	Scope string
}

type ImportTranslationArgument struct {
	TranslationsFile oss.OSS
}

// RegisterExchangeJobs register i18n jobs into worker
func RegisterExchangeJobs(I18n *i18n.I18N, Worker *worker.Worker) {
	if I18n.Resource == nil {
		debug.PrintStack()
		fmt.Println("I18N should be registered into `Admin` before register jobs")
		return
	}

	Admin := I18n.Resource.GetAdmin()
	Admin.RegisterViewPath("github.com/bhojpur/cms/pkg/i18n/exchange_actions/views")

	// Export Translations
	exportTranslationResource := Admin.NewResource(&ExportTranslationArgument{})
	exportTranslationResource.Meta(&admin.Meta{Name: "Scope", Type: "select_one", Collection: []string{"All", "Backend", "Frontend"}})

	Worker.RegisterJob(&worker.Job{
		Name:     "Export Translations",
		Group:    "Export/Import Translations From CSV file",
		Resource: exportTranslationResource,
		Handler: func(arg interface{}, bhojpurJob worker.BhojpurJobInterface) (err error) {
			var (
				locales          []string
				translationKeys  []string
				translationsMap  = map[string]bool{}
				filename         = fmt.Sprintf("/downloads/translations.%v.csv", time.Now().UnixNano())
				fullFilename     = filepath.Join("public", filename)
				i18nTranslations = I18n.LoadTranslations()
				scope            = arg.(*ExportTranslationArgument).Scope
			)
			bhojpurJob.AddLog("Exporting translations...")

			// Sort locales
			for locale := range i18nTranslations {
				locales = append(locales, locale)
			}
			sort.Strings(locales)

			// Create download file
			if _, err = os.Stat(filepath.Dir(fullFilename)); os.IsNotExist(err) {
				err = os.MkdirAll(filepath.Dir(fullFilename), os.ModePerm)
			}
			csvfile, err := os.OpenFile(fullFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			defer csvfile.Close()
			if err != nil {
				return err
			}

			writer := csv.NewWriter(csvfile)

			// Append Headers
			writer.Write(append([]string{"Translation Keys"}, locales...))

			// Sort translation keys
			for _, locale := range locales {
				for key := range i18nTranslations[locale] {
					translationsMap[key] = true
				}
			}

			for key := range translationsMap {
				translationKeys = append(translationKeys, key)
			}
			sort.Strings(translationKeys)

			// Write CSV file
			var (
				recordCount         = len(translationKeys)
				perCount            = recordCount/20 + 1
				processedRecordLogs = []string{}
				index               = 0
				progressCount       = 0
			)
			for _, translationKey := range translationKeys {
				// Filter out translation by scope
				index++
				if scope == "Backend" && !strings.HasPrefix(translationKey, "bhojpur_") {
					continue
				}
				if scope == "Frontend" && strings.HasPrefix(translationKey, "bhojpur_") {
					continue
				}
				var translations = []string{translationKey}
				for _, locale := range locales {
					var value string
					if translation := i18nTranslations[locale][translationKey]; translation != nil {
						value = translation.Value
					}
					translations = append(translations, value)
				}
				writer.Write(translations)
				processedRecordLogs = append(processedRecordLogs, fmt.Sprintf("Exported %v\n", strings.Join(translations, ",")))
				if index == perCount {
					bhojpurJob.AddLog(strings.Join(processedRecordLogs, ""))
					processedRecordLogs = []string{}
					progressCount++
					bhojpurJob.SetProgress(uint(float32(progressCount) / float32(20) * 100))
					index = 0
				}
			}
			writer.Flush()

			bhojpurJob.SetProgressText(fmt.Sprintf("<a href='%v'>Download exported translations</a>", filename))
			return
		},
	})

	// Import Translations

	Worker.RegisterJob(&worker.Job{
		Name:     "Import Translations",
		Group:    "Export/Import Translations From CSV file",
		Resource: Admin.NewResource(&ImportTranslationArgument{}),
		Handler: func(arg interface{}, bhojpurJob worker.BhojpurJobInterface) (err error) {
			importTranslationArgument := arg.(*ImportTranslationArgument)
			bhojpurJob.AddLog("Importing translations...")
			if csvfile, err := os.Open(filepath.Join("public", importTranslationArgument.TranslationsFile.URL())); err == nil {
				reader := csv.NewReader(csvfile)
				reader.TrimLeadingSpace = true
				if records, err := reader.ReadAll(); err == nil {
					if len(records) > 1 && len(records[0]) > 1 {
						var (
							recordCount         = len(records) - 1
							perCount            = recordCount/20 + 1
							processedRecordLogs = []string{}
							locales             = records[0][1:]
							index               = 1
						)
						for _, values := range records[1:] {
							logMsg := ""
							for idx, value := range values[1:] {
								if value == "" {
									if values[0] != "" && locales[idx] != "" {
										I18n.DeleteTranslation(&i18n.Translation{
											Key:    values[0],
											Locale: locales[idx],
										})
										logMsg += fmt.Sprintf("%v/%v Deleted %v,%v\n", index, recordCount, locales[idx], values[0])
									}
								} else {
									I18n.SaveTranslation(&i18n.Translation{
										Key:    values[0],
										Locale: locales[idx],
										Value:  value,
									})
									logMsg += fmt.Sprintf("%v/%v Imported %v,%v,%v\n", index, recordCount, locales[idx], values[0], value)
								}
							}
							processedRecordLogs = append(processedRecordLogs, logMsg)
							if len(processedRecordLogs) == perCount {
								bhojpurJob.AddLog(strings.Join(processedRecordLogs, ""))
								processedRecordLogs = []string{}
								bhojpurJob.SetProgress(uint(float32(index) / float32(recordCount+1) * 100))
							}
							index++
						}
						bhojpurJob.AddLog(strings.Join(processedRecordLogs, ""))
					}
				}
				bhojpurJob.AddLog("Imported translations")
			}
			return
		},
	})
}
