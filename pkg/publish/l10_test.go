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
	"testing"

	"github.com/bhojpur/cms/pkg/l10n"
	"github.com/bhojpur/cms/pkg/publish"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type Book struct {
	orm.Model
	l10n.Locale
	publish.Status
	Name        string
	CategoryID  uint
	Category    Category
	PublisherID uint
	Publisher   Publisher
	Comments    []Comment
	Authors     []Author `orm:"many2many:author_books;ForeignKey:ID;AssociationForeignKey:ID"`
}

type Publisher struct {
	orm.Model
	publish.Status
	Name string
}

type Comment struct {
	orm.Model
	l10n.Locale
	publish.Status
	Content string
	BookID  uint
}

type Author struct {
	orm.Model
	l10n.Locale
	Name string
}

func generateBook(name string) *Book {
	book := Book{
		Name: name,
		Category: Category{
			Name: name + "_category",
		},
		Publisher: Publisher{
			Name: name + "_publisher",
		},
		Comments: []Comment{
			{Content: name + "_comment1"},
			{Content: name + "_comment2"},
		},
		Authors: []Author{
			{Name: name + "_author1"},
			{Name: name + "_author2"},
		},
	}
	return &book
}

func TestBelongsToForL10nResource(t *testing.T) {
	name := "belongs_to_for_l10n"
	book := generateBook(name)
	pbdraft.Save(book)

	pb.Publish(book)

	if pbprod.Where("id = ?", book.ID).First(&Book{}).RecordNotFound() {
		t.Errorf("should find book from production db")
	}

	if pbprod.Where("name LIKE ?", name+"%").First(&Publisher{}).RecordNotFound() {
		t.Errorf("should find publisher from production db")
	}

	if pbprod.Where("name LIKE ?", name+"%").First(&Category{}).RecordNotFound() {
		t.Errorf("should find category from production db")
	}
}

func TestMany2ManyForL10nResource(t *testing.T) {
	name := "many2many_for_l10n"
	book := generateBook(name)
	pbdraft.Save(book)

	if pbdraft.Model(book).Association("Authors").Count() != 2 {
		t.Errorf("should find two authors from draft db before publish")
	}

	if pbprod.Model(book).Association("Authors").Count() != 0 {
		t.Errorf("should find none author from production db before publish")
	}

	pb.Publish(book)

	if pbprod.Where("id = ?", book.ID).First(&Book{}).RecordNotFound() {
		t.Errorf("should find book from production db")
	}

	if pbdraft.Model(book).Association("Authors").Count() != 2 {
		t.Errorf("should find two authors from draft db after publish")
	}

	if pbprod.Model(book).Association("Authors").Count() != 2 {
		t.Errorf("should find two authors from draft db after publish")
	}
}
