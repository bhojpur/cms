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

func TestCreateStructFromDraft(t *testing.T) {
	name := "create_product_from_draft"
	pbdraft.Create(&Product{Name: name, Color: Color{Name: name}})

	if !pbprod.First(&Product{}, "name = ?", name).RecordNotFound() {
		t.Errorf("record should not be found in production db")
	}

	if pbdraft.First(&Product{}, "name = ?", name).RecordNotFound() {
		t.Errorf("record should be found in draft db")
	}

	if pbprod.Table("colors").First(&Color{}, "name = ?", name).Error != nil {
		t.Errorf("color should be saved")
	}

	if pbprod.Table("colors_draft").First(&Color{}, "name = ?", name).Error == nil {
		t.Errorf("no colors_draft table")
	}

	var product Product
	pbdraft.First(&product, "name = ?", name)

	if !product.PublishStatus {
		t.Errorf("Product's publish status should be DIRTY when created from draft db")
	}

	if pbdraft.Model(&product).Related(&product.Color); product.Color.Name != name {
		t.Errorf("should be able to find related struct")
	}
}

func TestCreateStructFromProduction(t *testing.T) {
	name := "create_product_from_production"
	pbprod.Create(&Product{Name: name, Color: Color{Name: name}})

	if pbprod.First(&Product{}, "name = ?", name).RecordNotFound() {
		t.Errorf("record should not be found in production db")
	}

	if pbdraft.First(&Product{}, "name = ?", name).RecordNotFound() {
		t.Errorf("record should be found in draft db")
	}

	if pbprod.Table("colors").First(&Color{}, "name = ?", name).Error != nil {
		t.Errorf("color should be saved")
	}

	var product Product
	pbprod.First(&product, "name = ?", name)

	if product.PublishStatus {
		t.Errorf("Product's publish status should be PUBLISHED when created from production db")
	}

	if pbprod.Model(&product).Related(&product.Color); product.Color.Name != name {
		t.Errorf("should be able to find related struct")
	}
}
