package database

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

	"github.com/bhojpur/cms/pkg/i18n"
	orm "github.com/bhojpur/orm/pkg/engine"
)

// Translation is a struct used to save translations into databae
type Translation struct {
	Locale string `sql:"size:12;"`
	Key    string `sql:"size:4294967295;"`
	Value  string `sql:"size:4294967295"`
}

// New new DB backend for I18n
func New(db *orm.DB) i18n.Backend {
	db.AutoMigrate(&Translation{})
	if err := db.Model(&Translation{}).AddUniqueIndex("idx_translations_key_with_locale", "locale", "key").Error; err != nil {
		fmt.Printf("Failed to create unique index for translations key & locale, got: %v\n", err.Error())
	}
	return &Backend{DB: db}
}

// Backend DB backend
type Backend struct {
	DB *orm.DB
}

// LoadTranslations load translations from DB backend
func (backend *Backend) LoadTranslations() (translations []*i18n.Translation) {
	backend.DB.Find(&translations)
	return translations
}

// SaveTranslation save translation into DB backend
func (backend *Backend) SaveTranslation(t *i18n.Translation) error {
	return backend.DB.Where(Translation{Key: t.Key, Locale: t.Locale}).
		Assign(Translation{Value: t.Value}).
		FirstOrCreate(&Translation{}).Error
}

// FindTranslation find translation from DB backend
func (backend *Backend) FindTranslation(t *i18n.Translation) (translation i18n.Translation) {
	backend.DB.Where(Translation{Key: t.Key, Locale: t.Locale}).Find(&translation)
	return translation
}

// DeleteTranslation delete translation into DB backend
func (backend *Backend) DeleteTranslation(t *i18n.Translation) error {
	return backend.DB.Where(Translation{Key: t.Key, Locale: t.Locale}).Delete(&Translation{}).Error
}
