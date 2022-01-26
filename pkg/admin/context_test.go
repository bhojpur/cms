package admin_test

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

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/cms/pkg/admin"
	. "github.com/bhojpur/cms/tests/dummy"
)

// Template helpers test

func TestUrlForAdmin(t *testing.T) {
	context := &admin.Context{Admin: Admin}

	rootLink := context.URLFor(Admin)

	if rootLink != "/admin" {
		t.Error("Admin link not generated by URL for")
	}
}

func TestUrlForResource(t *testing.T) {
	context := &admin.Context{Admin: Admin}
	user := Admin.GetResource("User")

	userLink := context.URLFor(user)

	if userLink != "/admin/users" {
		t.Error("resource link not generated by URL for")
	}
}

type Store struct {
	ID   string `gorm:"primary_key" sql:"type:varchar(20)"`
	Name string
}

func TestUrlForResourceWithSpecialPrimaryKey(t *testing.T) {
	db.AutoMigrate(&Store{})
	context := &admin.Context{Admin: Admin, Context: &appsvr.Context{}}
	context.SetDB(db)

	storeRes := Admin.AddResource(&Store{}, &admin.Config{Permission: roles.Allow(roles.CRUD, roles.Anyone)})
	s := Store{ID: "00022 %alert%29", Name: "test"}
	if err := db.Save(&s).Error; err != nil {
		t.Fatal(err)
	}

	storeLink := context.URLFor(s, storeRes)
	if storeLink != "/admin/stores/00022%20%25alert%2529" {
		t.Error("special primary key is not escaped properly")
	}
}

func TestUrlForResourceName(t *testing.T) {
	user := &User{Name: "test"}
	db.Create(&user)

	context := &admin.Context{Admin: Admin, Context: &appsvr.Context{}}
	context.SetDB(db)

	userLink := context.URLFor(user)

	if userLink != "/admin/users/"+fmt.Sprintf("%v", user.ID) {
		t.Error("resource link not generated by URL for")
	}
}

func TestPagination(t *testing.T) {
	context := &admin.Context{Admin: Admin}
	context.Resource = &admin.Resource{Config: &admin.Config{PageCount: 10}}
	context.Searcher = &admin.Searcher{Context: context}

	// Test no pagination if total result count is less than PageCount
	for _, count := range []int{8, 10} {
		context.Searcher.Pagination.Total = count
		if context.Pagination() != nil {
			t.Error(fmt.Sprintf("Don't display pagination if only has one page (%v)", count))
		}
	}

	context.Searcher.Pagination.CurrentPage = 3
	if context.Pagination() == nil {
		t.Error("Should show pagination for page without records")
	}

	// Test current page 1
	context.Searcher.Pagination.Total = 1000
	context.Searcher.Pagination.Pages = 10
	context.Searcher.Pagination.CurrentPage = 1
	pages := context.Pagination().Pages

	if !pages[0].Current {
		t.Error("first page not set as current page")
	}

	if !pages[len(pages)-2].IsNext && pages[len(pages)-2].Page != 2 {
		t.Error("Should have next page arrow")
	}

	// +1 for "Next page" link which is a "Page" too
	// +1 for "Last page"
	if len(pages) != 8+1+1 {
		t.Error("visible pages in current context beyond the bound of VISIBLE_PAGE_COUNT")
	}

	// Test current page 8 => the length between start and end less than MAX_VISIBLE_PAGES
	context.Searcher.Pagination.Pages = 10
	context.Searcher.Pagination.CurrentPage = 8
	pages = context.Pagination().Pages

	if !pages[7].Current {
		t.Error("visible previous pages count incorrect")
	}

	if !pages[1].IsPrevious && pages[1].Page != 7 {
		t.Error("Should have previous page arrow")
	}

	// +1 for "Prev"
	// +1 for "First page"
	if len(pages) != 8+1+1 {
		t.Error("visible pages in current context beyond the bound of VISIBLE_PAGE_COUNT")
	}

	// Test current page at last
	context.Searcher.Pagination.Pages = 10
	context.Searcher.Pagination.CurrentPage = 10
	pages = context.Pagination().Pages

	if !pages[len(pages)-1].Current {
		t.Error("last page is not the current page")
	}

	if len(pages) != 8+2 {
		t.Error("visible pages count is incorrect")
	}

	// Test current page at last but total page count less than VISIBLE_PAGE_COUNT
	context.Searcher.Pagination.Pages = 5
	context.Searcher.Pagination.CurrentPage = 5
	pages = context.Pagination().Pages

	if len(pages) != 5 {
		t.Error("incorrect pages count")
	}
}
