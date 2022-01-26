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
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/resource"
	bhojpurTestUtils "github.com/bhojpur/application/test/utils"
	. "github.com/bhojpur/cms/tests/dummy"
)

func TestUpdateRecord(t *testing.T) {
	bhojpurTestUtils.ResetDBTables(db, &Language{}, &User{})
	createLoggedInUser()

	user := User{Name: "update_record", Role: Role_system_administrator}
	db.Save(&user)

	form := url.Values{
		"BhojpurResource.Name": {user.Name + "_new"},
		"BhojpurResource.Role": {Role_system_administrator},
	}

	if req, err := http.PostForm(server.URL+"/admin/users/"+fmt.Sprint(user.ID), form); err == nil {
		if req.StatusCode != 200 {
			t.Errorf("Create request should be processed successfully")
		}

		if db.First(&User{}, "name = ?", user.Name+"_new").RecordNotFound() {
			t.Errorf("User should be updated successfully")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func TestUpdateRecordWithRollback(t *testing.T) {
	bhojpurTestUtils.ResetDBTables(db, &Language{}, &User{})
	createLoggedInUser()

	db.Model(&User{}).AddUniqueIndex("uix_user_name", "name")

	userR := Admin.GetResource("User")
	userR.AddProcessor(&resource.Processor{
		Name: "product-admin-prroduct-res-processor",
		Handler: func(v interface{}, meta *resource.MetaValues, c *appsvr.Context) error {
			user := v.(*User)
			c.DB.Model(user).Association("Languages").Replace([]Language{{Name: "CN"}})
			return nil
		},
	})

	anotherUsersName := "Katin"
	db.Save(&User{Name: anotherUsersName, Role: Role_system_administrator})

	user := User{Name: "update_record", Role: Role_system_administrator, Languages: []Language{{Name: "CN"}, {Name: "JP"}}}
	db.Save(&user)

	form := url.Values{
		"BhojpurResource.Name": {anotherUsersName},
		"BhojpurResource.Role": {Role_system_administrator},
	}

	if req, err := http.PostForm(server.URL+"/admin/users/"+fmt.Sprint(user.ID), form); err == nil {
		if req.StatusCode == 200 {
			t.Errorf("Should update user failure when name already be token by other user.")
		}

		u := User{}
		if err := db.Where("name = 'update_record'").Preload("Languages").First(&u).Error; err != nil {
			t.Fatal(err)
		}

		languages := []string{}
		for _, language := range u.Languages {
			languages = append(languages, language.Name)
		}

		if strings.Join(languages, ",") != "CN,JP" {
			t.Errorf("Should keep origin value for languages, but got %v", strings.Join(languages, ","))
		}
	} else {
		t.Errorf(err.Error())
	}
}

func TestUpdateHasOneRecord(t *testing.T) {
	bhojpurTestUtils.ResetDBTables(db, &CreditCard{}, &User{})
	createLoggedInUser()

	user := User{Name: "update_record_and_has_one", Role: Role_system_administrator, CreditCard: CreditCard{Number: "1234567890", Issuer: "JCB"}}
	db.Save(&user)

	form := url.Values{
		"BhojpurResource.Name":              {user.Name + "_new"},
		"BhojpurResource.Role":              {Role_system_administrator},
		"BhojpurResource.CreditCard.ID":     {fmt.Sprint(user.CreditCard.ID)},
		"BhojpurResource.CreditCard.Number": {"1234567890"},
		"BhojpurResource.CreditCard.Issuer": {"UnionPay"},
	}

	if req, err := http.PostForm(server.URL+"/admin/users/"+fmt.Sprint(user.ID), form); err == nil {
		if req.StatusCode != 200 {
			t.Errorf("User request should be processed successfully")
		}

		if db.First(&User{}, "name = ?", user.Name+"_new").RecordNotFound() {
			t.Errorf("User should be updated successfully")
		}

		var creditCard CreditCard
		if db.Model(&user).Related(&creditCard).RecordNotFound() ||
			creditCard.Issuer != "UnionPay" || creditCard.ID != user.CreditCard.ID {
			t.Errorf("Embedded struct should be updated successfully")
		}

		if !db.First(&CreditCard{}, "number = ? and issuer = ?", "1234567890", "JCB").RecordNotFound() {
			t.Errorf("Old embedded struct should be updated")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func TestUpdateHasManyRecord(t *testing.T) {
	bhojpurTestUtils.ResetDBTables(db, &Address{}, &User{})
	createLoggedInUser()

	user := User{Name: "update_record_and_has_many", Role: Role_system_administrator, Addresses: []Address{{Address1: "address 1.1", Address2: "address 1.2"}, {Address1: "address 2.1"}, {Address1: "address 3.1"}}}
	db.Save(&user)

	form := url.Values{
		"BhojpurResource.Name":                  {user.Name},
		"BhojpurResource.Role":                  {Role_system_administrator},
		"BhojpurResource.Addresses[0].ID":       {fmt.Sprint(user.Addresses[0].ID)},
		"BhojpurResource.Addresses[0].Address1": {"address 1.1 new"},
		"BhojpurResource.Addresses[1].ID":       {fmt.Sprint(user.Addresses[1].ID)},
		"BhojpurResource.Addresses[1].Address1": {"address 2.1 new"},
		"BhojpurResource.Addresses[2].ID":       {fmt.Sprint(user.Addresses[2].ID)},
		"BhojpurResource.Addresses[2]._destroy": {"1"},
		"BhojpurResource.Addresses[2].Address1": {"address 3.1"},
		"BhojpurResource.Addresses[3].Address1": {"address 4.1"},
	}

	if req, err := http.PostForm(server.URL+"/admin/users/"+fmt.Sprint(user.ID), form); err == nil {
		if req.StatusCode != 200 {
			t.Errorf("Create request should be processed successfully")
		}

		var address1 Address
		if db.First(&address1, "user_id = ? and address1 = ?", user.ID, "address 1.1 new").RecordNotFound() {
			t.Errorf("Address 1 should be updated successfully")
		} else if address1.Address2 != "address 1.2" {
			t.Errorf("Address 1's Address 2 should not be updated")
		}

		if db.First(&Address{}, "user_id = ? and address1 = ?", user.ID, "address 2.1 new").RecordNotFound() {
			t.Errorf("Address 2 should be updated successfully")
		}

		if !db.First(&Address{}, "user_id = ? and address1 = ?", user.ID, "address 3.1").RecordNotFound() {
			t.Errorf("Address 3 should be destroyed successfully")
		}

		if db.First(&Address{}, "user_id = ? and address1 = ?", user.ID, "address 4.1").RecordNotFound() {
			t.Errorf("Address 4 should be created successfully")
		}

		var addresses []Address
		if db.Find(&addresses, "user_id = ?", user.ID); len(addresses) != 3 {
			t.Errorf("Addresses's count should be updated after update")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func TestDestroyEmbeddedHasOneRecord(t *testing.T) {
	bhojpurTestUtils.ResetDBTables(db, &CreditCard{}, &User{})
	createLoggedInUser()

	user := User{Name: "destroy_embedded_has_one_record", Role: Role_system_administrator, CreditCard: CreditCard{Number: "1234567890", Issuer: "JCB"}}
	db.Save(&user)

	form := url.Values{
		"BhojpurResource.Name":                {user.Name + "_new"},
		"BhojpurResource.Role":                {Role_system_administrator},
		"BhojpurResource.CreditCard.ID":       {fmt.Sprint(user.CreditCard.ID)},
		"BhojpurResource.CreditCard._destroy": {"1"},
		"BhojpurResource.CreditCard.Number":   {"1234567890"},
		"BhojpurResource.CreditCard.Issuer":   {"UnionPay"},
	}

	if req, err := http.PostForm(server.URL+"/admin/users/"+fmt.Sprint(user.ID), form); err == nil {
		if req.StatusCode != 200 {
			t.Errorf("User request should be processed successfully")
		}

		var newUser User
		if db.First(&newUser, "name = ?", user.Name+"_new").RecordNotFound() {
			t.Errorf("User should be updated successfully")
		}

		if !db.Model(&newUser).Related(&CreditCard{}).RecordNotFound() {
			t.Errorf("Embedded struct should be destroyed successfully")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func TestUpdateManyToManyRecord(t *testing.T) {
	bhojpurTestUtils.ResetDBTables(db, &Language{}, &User{})
	createLoggedInUser()

	name := "update_record_many_to_many"
	var languageCN Language
	var languageEN Language
	db.FirstOrCreate(&languageCN, Language{Name: "CN"})
	db.FirstOrCreate(&languageEN, Language{Name: "EN"})
	user := User{Name: name, Role: Role_system_administrator, Languages: []Language{languageCN, languageEN}}
	db.Save(&user)

	form := url.Values{
		"BhojpurResource.Name":      {name + "_new"},
		"BhojpurResource.Role":      {Role_system_administrator},
		"BhojpurResource.Languages": {fmt.Sprint(languageCN.ID)},
	}

	if req, err := http.PostForm(server.URL+"/admin/users/"+fmt.Sprint(user.ID), form); err == nil {
		if req.StatusCode != 200 {
			t.Errorf("Update request should be processed successfully")
		}

		var user User
		if db.First(&user, "name = ?", name+"_new").RecordNotFound() {
			t.Errorf("User should be updated successfully")
		}

		var languages []Language
		db.Model(&user).Related(&languages, "Languages")

		if len(languages) != 1 {
			t.Errorf("User should have one languages after update")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func TestUpdateSelectOne(t *testing.T) {
	bhojpurTestUtils.ResetDBTables(db, &Company{}, &User{})
	createLoggedInUser()

	name := "update_record_select_one"
	var company1, company2 Company
	if err := db.FirstOrCreate(&company1, &Company{Name: "Company 1"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.FirstOrCreate(&company2, &Company{Name: "Company 2"}).Error; err != nil {
		t.Fatal(err)
	}
	user := User{Name: name, Role: Role_system_administrator, Company: &company1}
	db.Save(&user)

	form := url.Values{
		"BhojpurResource.Name":    {name + "_new"},
		"BhojpurResource.Role":    {Role_system_administrator},
		"BhojpurResource.Company": {fmt.Sprint(company2.ID)},
	}

	if req, err := http.PostForm(server.URL+"/admin/users/"+fmt.Sprint(user.ID), form); err == nil {
		if req.StatusCode != 200 {
			t.Errorf("Update request should be processed successfully")
		}

		var user User
		if db.Preload("Company").First(&user, "name = ?", name+"_new").RecordNotFound() {
			t.Errorf("User should be updated successfully")
		}

		if user.Company.ID != company2.ID {
			t.Errorf("user's company should be updated")
		}
	} else {
		t.Errorf(err.Error())
	}
}

func TestUpdateAttachment(t *testing.T) {
	name := "update_record_attachment"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if attachment, err := filepath.Abs("tests/bhojpur.png"); err == nil {
		if part, err := writer.CreateFormFile("BhojpurResource.Avatar", filepath.Base(attachment)); err == nil {
			if file, err := os.Open(attachment); err == nil {
				io.Copy(part, file)
			}
		}
		form := url.Values{
			"BhojpurResource.Name": {name},
			"BhojpurResource.Role": {Role_system_administrator},
		}
		for key, val := range form {
			_ = writer.WriteField(key, val[0])
		}
		writer.Close()

		var user User
		if req, err := http.Post(server.URL+"/admin/users", writer.FormDataContentType(), body); err == nil {
			if req.StatusCode != 200 {
				t.Errorf("Create request should be processed successfully")
			}

			if db.First(&user, "name = ?", name).RecordNotFound() {
				t.Errorf("User should be created successfully")
			}

			if !regexp.MustCompile("bhojpur").MatchString(user.Avatar.URL()) {
				t.Errorf("Avatar should be saved, but its URL is %v", user.Avatar.URL())
			}
		}

		attachment, err := filepath.Abs("tests/logo.png")
		if err != nil {
			panic(err)
		}
		if part, err := writer.CreateFormFile("BhojpurResource.Avatar", filepath.Base(attachment)); err == nil {
			if file, err := os.Open(attachment); err == nil {
				io.Copy(part, file)
			}
		}
		for key, val := range form {
			_ = writer.WriteField(key, val[0])
		}
		writer.Close()

		if req, err := http.Post(fmt.Sprintf("%v/admin/users/%v", server.URL, user.ID), writer.FormDataContentType(), body); err == nil {
			if req.StatusCode != 200 {
				t.Errorf("Create request should be processed successfully")
			}

			if db.First(&user, "name = ?", name).RecordNotFound() {
				t.Errorf("User should be created successfully")
			}

			if !regexp.MustCompile("logo").MatchString(user.Avatar.URL()) {
				t.Errorf("Avatar should be updated, but its URL is %v", user.Avatar.URL())
			}
		}
	}
}
