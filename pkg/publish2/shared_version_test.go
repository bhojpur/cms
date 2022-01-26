package publish2_test

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

	"github.com/bhojpur/cms/pkg/publish2"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type SharedVersionProduct struct {
	orm.Model
	Name            string
	ColorVariations []SharedVersionColorVariation
	publish2.Version
}

type SharedVersionColorVariation struct {
	orm.Model
	Name                   string
	SharedVersionProductID uint
	SizeVariations         []SharedVersionSizeVariation
	publish2.SharedVersion
}

type SharedVersionSizeVariation struct {
	orm.Model
	Name                          string
	SharedVersionColorVariationID uint
	publish2.SharedVersion
}

func prepareSharedVersionProduct() *SharedVersionProduct {
	product := SharedVersionProduct{
		Name: "shared product 1",
		ColorVariations: []SharedVersionColorVariation{
			{
				Name: "cv1",
			},
			{
				Name: "cv2",
			},
		},
	}
	DB.Create(&product)

	product.SetVersionName("v1")
	product.ColorVariations[0].SetSharedVersionName("v1")
	DB.Save(&product)

	product.SetVersionName("v2")
	product.ColorVariations[0].SetSharedVersionName("")
	colorVariation := SharedVersionColorVariation{
		Name: "cv3",
	}
	colorVariation.SetSharedVersionName("v2")
	product.ColorVariations = append(product.ColorVariations, colorVariation)
	DB.Save(&product)

	return &product
}

func TestSharedVersions(t *testing.T) {
	product1 := prepareSharedVersionProduct()
	product2 := prepareSharedVersionProduct()

	var product1V1 SharedVersionProduct
	DB.Set(publish2.VersionNameMode, "v1").Preload("ColorVariations").Find(&product1V1, "id = ?", product1.ID)

	if len(product1V1.ColorVariations) != 2 {
		t.Errorf("Preload: Should have 2 color variations for product v1, but got %v", len(product1V1.ColorVariations))
	}

	var colorVariations1V1 []SharedVersionColorVariation
	DB.Model(&product1V1).Related(&colorVariations1V1)
	if len(colorVariations1V1) != 2 {
		t.Errorf("Related: Should have 2 color variations for product v1, but got %v", len(colorVariations1V1))
	}

	var product1V2 SharedVersionProduct
	DB.Set(publish2.VersionNameMode, "v2").Preload("ColorVariations").Find(&product1V2, "id = ?", product1.ID)

	if len(product1V2.ColorVariations) != 3 {
		t.Errorf("Preload: Should have 3 color variations for product v2, but got %v", len(product1V2.ColorVariations))
	}

	var products []SharedVersionProduct
	DB.Preload("ColorVariations").Find(&products)

	var product2V2 SharedVersionProduct
	for _, p := range products {
		if p.ID == product2.ID && p.VersionName == "v2" {
			product2V2 = p
		}
	}

	if len(product2V2.ColorVariations) != 3 {
		t.Errorf("Preload: Should have 3 color variations for product v2, but got %v", len(product2V2.ColorVariations))
	}
}
