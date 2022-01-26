package l10n_test

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
	"time"

	"github.com/bhojpur/application/test/utils"
	"github.com/bhojpur/cms/pkg/l10n"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type Product struct {
	ID              int    `gorm:"primary_key"`
	Code            string `l10n:"sync"`
	Quantity        uint   `l10n:"sync"`
	Name            string
	DeletedAt       *time.Time
	ColorVariations []ColorVariation
	BrandID         uint `l10n:"sync"`
	Brand           Brand
	Tags            []Tag      `gorm:"many2many:product_tags"`
	Categories      []Category `gorm:"many2many:product_categories;ForeignKey:id;AssociationForeignKey:id"`
	l10n.Locale
}

// func (Product) LocaleCreatable() {}

type ColorVariation struct {
	ID       int `gorm:"primary_key"`
	Quantity int
	Color    Color
}

type Color struct {
	ID   int `gorm:"primary_key"`
	Code string
	Name string
	l10n.Locale
}

type Brand struct {
	ID   int `gorm:"primary_key"`
	Name string
	l10n.Locale
}

type Tag struct {
	ID   int `gorm:"primary_key"`
	Name string
	l10n.Locale
}

type Category struct {
	ID   int `gorm:"primary_key"`
	Name string
	l10n.Locale
}

var dbGlobal, dbCN, dbEN *orm.DB

func init() {
	db := utils.TestDB()
	l10n.RegisterCallbacks(db)

	db.DropTableIfExists(&Product{})
	db.DropTableIfExists(&Brand{})
	db.DropTableIfExists(&Tag{})
	db.DropTableIfExists(&Category{})
	db.Exec("drop table product_tags;")
	db.Exec("drop table product_categories;")
	db.AutoMigrate(&Product{}, &Brand{}, &Tag{}, &Category{})

	dbGlobal = db
	dbCN = dbGlobal.Set("l10n:locale", "zh")
	dbEN = dbGlobal.Set("l10n:locale", "en")
}
