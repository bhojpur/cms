package publish_test

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

	"github.com/bhojpur/application/test/utils"
	"github.com/bhojpur/cms/pkg/l10n"
	"github.com/bhojpur/cms/pkg/publish"
	orm "github.com/bhojpur/orm/pkg/engine"
	_ "github.com/mattn/go-sqlite3"
)

var pb *publish.Publish
var pbdraft *orm.DB
var pbprod *orm.DB
var db *orm.DB

func init() {
	db = utils.TestDB()
	l10n.RegisterCallbacks(db)

	pb = publish.New(db)
	pbdraft = pb.DraftDB()
	pbprod = pb.ProductionDB()

	for _, table := range []string{"product_categories", "product_categories_draft", "product_languages", "product_languages_draft", "author_books", "author_books_draft"} {
		pbprod.Exec(fmt.Sprintf("drop table %v", table))
	}

	for _, value := range []interface{}{&Product{}, &Color{}, &Category{}, &Language{}, &Book{}, &Publisher{}, &Comment{}, &Author{}} {
		pbprod.DropTable(value)
		pbdraft.DropTable(value)

		pbprod.AutoMigrate(value)
		pb.AutoMigrate(value)
	}
}

type Product struct {
	orm.Model
	Name       string
	Quantity   uint
	Color      Color
	ColorId    int
	Categories []Category `orm:"many2many:product_categories"`
	Languages  []Language `orm:"many2many:product_languages"`
	publish.Status
}

type Color struct {
	orm.Model
	Name string
}

type Language struct {
	orm.Model
	Name string
}

type Category struct {
	orm.Model
	Name string
	publish.Status
}
