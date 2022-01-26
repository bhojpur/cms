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

import "testing"

func TestPublishManyToManyFromProduction(t *testing.T) {
	name := "create_product_with_multi_categories_from_production"
	pbprod.Create(&Product{
		Name:       name,
		Categories: []Category{{Name: "category1"}, {Name: "category2"}},
	})

	var product Product
	pbprod.First(&product, "name = ?", name)

	if pbprod.Model(&product).Association("Categories").Count() != 2 {
		t.Errorf("categories count should be 2 in production db")
	}

	if pbdraft.Model(&product).Association("Categories").Count() != 2 {
		t.Errorf("categories count should be 2 in draft db")
	}
}

func TestPublishManyToManyFromDraft(t *testing.T) {
	name := "create_product_with_multi_categories_from_draft"
	pbdraft.Create(&Product{
		Name:       name,
		Categories: []Category{{Name: "category1"}, {Name: "category2"}},
	})

	var product Product
	pbdraft.First(&product, "name = ?", name)

	if pbprod.Model(&product).Association("Categories").Count() != 0 {
		t.Errorf("categories count should be 0 in production db")
	}

	if pbdraft.Model(&product).Association("Categories").Count() != 2 {
		t.Errorf("categories count should be 2 in draft db")
	}

	pb.Publish(&product)
	var categories []Category
	pbdraft.Find(&categories)
	pb.Publish(&categories)

	if pbprod.Model(&product).Association("Categories").Count() != 2 {
		t.Errorf("categories count should be 2 in production db after publish")
	}

	if pbdraft.Model(&product).Association("Categories").Count() != 2 {
		t.Errorf("categories count should be 2 in draft db after publish")
	}
}

func TestDiscardManyToManyFromDraft(t *testing.T) {
	name := "discard_product_with_multi_categories_from_draft"
	pbdraft.Create(&Product{
		Name:       name,
		Categories: []Category{{Name: "category1"}, {Name: "category2"}},
	})

	var product Product
	pbdraft.First(&product, "name = ?", name)

	if pbprod.Model(&product).Association("Categories").Count() != 0 {
		t.Errorf("categories count should be 0 in production db")
	}

	if pbdraft.Model(&product).Association("Categories").Count() != 2 {
		t.Errorf("categories count should be 2 in draft db")
	}

	pb.Discard(&product)

	if pbprod.Model(&product).Association("Categories").Count() != 0 {
		t.Errorf("categories count should be 0 in production db after discard")
	}

	if pbdraft.Model(&product).Association("Categories").Count() != 0 {
		t.Errorf("categories count should be 0 in draft db after discard")
	}
}
