package sorting_test

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
	"testing"

	"github.com/bhojpur/cms/pkg/l10n"
	"github.com/bhojpur/cms/pkg/sorting"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type Brand struct {
	orm.Model
	l10n.Locale
	sorting.Sorting
	Name string
}

func prepareBrand() {
	db.Delete(&Brand{})
	globalDB := db.Set("l10n:locale", l10n.Global)
	zhDB := db.Set("l10n:locale", "zh-CN")

	for i := 1; i <= 5; i++ {
		brand := Brand{Name: fmt.Sprintf("brand%v", i)}
		globalDB.Save(&brand)
		if i > 3 {
			zhDB.Save(&brand)
		}
	}
}

func getBrand(db *orm.DB, name string) *Brand {
	var brand Brand
	db.First(&brand, "name = ?", name)
	return &brand
}

func checkBrandPosition(db *orm.DB, t *testing.T, description string) {
	var brands []Brand
	if err := db.Set("l10n:mode", "locale").Find(&brands).Error; err != nil {
		t.Errorf("no error should happen when find brands, but got %v", err)
	}
	for i, brand := range brands {
		if brand.Position != i+1 {
			t.Errorf("Brand %v(%v)'s position should be %v after %v, but got %v", brand.ID, brand.LanguageCode, i+1, description, brand.Position)
		}
	}
}

func TestBrandPosition(t *testing.T) {
	prepareBrand()
	globalDB := db.Set("l10n:locale", "en-US")
	zhDB := db.Set("l10n:locale", "zh-CN")

	checkBrandPosition(db, t, "initalize")
	checkBrandPosition(zhDB, t, "initalize")

	if err := globalDB.Delete(getBrand(globalDB, "brand1")).Error; err != nil {
		t.Errorf("no error should happen when delete an en-US brand, but got %v", err)
	}
	checkBrandPosition(globalDB, t, "delete an brand from global db")
	checkBrandPosition(zhDB, t, "delete an brand from global db")

	if err := zhDB.Delete(getBrand(zhDB, "brand4")).Error; err != nil {
		t.Errorf("no error should happen when delete an zh-CN brand, but got %v", err)
	}
	checkBrandPosition(globalDB, t, "delete an brand from zh db")
	checkBrandPosition(zhDB, t, "delete an brand from zh db")
}
