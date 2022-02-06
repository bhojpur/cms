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
	"net/http"
	"net/http/httptest"
	"testing"

	bhojpurTestUtils "github.com/bhojpur/application/test/utils"
	. "github.com/bhojpur/cms/pkg/admin/tests/dummy"
	"github.com/theplant/htmltestingutils"
)

func TestEditPage(t *testing.T) {
	bhojpurTestUtils.ResetDBTables(db, &Language{}, &User{})
	user := createLoggedInUser()

	h := adminHandler

	var req *http.Request
	req, err := http.NewRequest("GET", fmt.Sprintf("/admin/users/%d/edit", user.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	// assert only form so that other sections won't break the tests
	diff := htmltestingutils.PrettyHtmlDiff(rr.Body, "form.bhojpur-form", expectedBody)
	if len(diff) > 0 {
		t.Error(diff)
	}
}

var expectedBody = `
<form class="bhojpur-form" action="/admin/users/1" method="POST" enctype="multipart/form-data">
<input name="_method" value="PUT" type="hidden">
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <input id="user_1_id" class="bhojpur-hidden__primary_key" name="BhojpurResource.ID" value="1" type="hidden">
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="user_1_name">
Name
</label>
<div class="bhojpur-field__show">Bhojpur CMS</div>
<div class="bhojpur-field__edit">
<input class="mdl-textfield__input" type="text" id="user_1_name" name="BhojpurResource.Name" value="BHOJPUR" >
</div>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="user_1_age">
Age
</label>
<div class="bhojpur-field__show">
0
</div>
<div class="bhojpur-field__edit">
<input class="mdl-textfield__input" type="number" id="user_1_age" name="BhojpurResource.Age" value="0" >
</div>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="user_1_role">
Role
</label>
<div class="bhojpur-field__show">admin</div>
<div class="bhojpur-field__edit">
<input class="mdl-textfield__input" type="text" id="user_1_role" name="BhojpurResource.Role" value="admin" >
</div>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<label class="mdl-checkbox mdl-js-checkbox mdl-js-ripple-effect" for="user_1_active">
<span class="bhojpur-field__label mdl-checkbox__label">Active</span>
<span class="bhojpur-field__edit">
<input type="checkbox" id="user_1_active" name="BhojpurResource.Active" class="mdl-checkbox__input" value="true" type="checkbox" >
<input type="hidden" name="BhojpurResource.Active" value="false">
</span>
</label>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="user_1_registered_at">
Registered At
</label>
<div class="bhojpur-field__show">
</div>
<div class="bhojpur-field__edit bhojpur-field__datetimepicker" data-picker-type="datetime">
<input class="mdl-textfield__input bhojpur-datetimepicker__input" placeholder=" YYYY-MM-DD HH:MM " type="text" id="user_1_registered_at" name="BhojpurResource.RegisteredAt" value="" >
<div>
	<button data-toggle="bhojpur.datepicker" class="mdl-button mdl-js-button mdl-button--icon bhojpur-action__datepicker" type="button">
	<i class="material-icons">date_range</i>
  </button>
  <button data-toggle="bhojpur.timepicker" class="mdl-button mdl-js-button mdl-button--icon bhojpur-action__timepicker" type="button">
	<i class="material-icons">access_time</i>
  </button>
</div>
</div>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
<div class="bhojpur-field">
<label class="bhojpur-field__label" for="user_1_avatar">
Avatar
</label>
<div class="bhojpur-field__block bhojpur-file ">
<div class="bhojpur-fieldset">
<textarea class="bhojpur-file__options hidden" data-cropper-title="Crop image" data-cropper-cancel="Cancel" data-cropper-ok="OK" name="BhojpurResource.Avatar" aria-hidden="true">{&#34;FileName&#34;:&#34;&#34;,&#34;Url&#34;:&#34;&#34;}</textarea>
<div class="bhojpur-file__list">
</div>
<label class="mdl-button mdl-button--primary bhojpur-button__icon-add" title="Choose File" >
  <input class="visuallyhidden bhojpur-file__input" id="user_1_avatar" name="BhojpurResource.Avatar" type="file">
  Add Avatar
</label>
</div>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
<div class="signle-edit bhojpur-field">
<label class="bhojpur-field__label" for="user_1_profile">
Profile
</label>
<div class="bhojpur-field__block">
<fieldset id="user_1_profile" class="bhojpur-fieldset">
  <input id="" class="bhojpur-hidden__primary_key" name="BhojpurResource.Profile.ID" value="0" type="hidden">
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <input id="" class="bhojpur-hidden__primary_key" name="BhojpurResource.Profile.ID" value="0" type="hidden">
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="">
Name
</label>
<div class="bhojpur-field__show"></div>
<div class="bhojpur-field__edit">
<input class="mdl-textfield__input" type="text" id="" name="BhojpurResource.Profile.Name" value="" >
</div>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="">
Sex
</label>
<div class="bhojpur-field__show"></div>
<div class="bhojpur-field__edit">
<input class="mdl-textfield__input" type="text" id="" name="BhojpurResource.Profile.Sex" value="" >
</div>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
<div class="signle-edit bhojpur-field">
<label class="bhojpur-field__label" for="">
Phone
</label>
<div class="bhojpur-field__block">
<fieldset id="" class="bhojpur-fieldset">
  <input id="" class="bhojpur-hidden__primary_key" name="BhojpurResource.Profile.Phone.ID" value="0" type="hidden">
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <input id="" class="bhojpur-hidden__primary_key" name="BhojpurResource.Profile.Phone.ID" value="0" type="hidden">
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="">
Num
</label>
<div class="bhojpur-field__show"></div>
<div class="bhojpur-field__edit">
<input class="mdl-textfield__input" type="text" id="" name="BhojpurResource.Profile.Phone.Num" value="" >
</div>
</div>
</div>
</div>
</div>
</div>
</fieldset>
</div>
</div>
</div>
</div>
</div>
</fieldset>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
<div class="signle-edit bhojpur-field">
<label class="bhojpur-field__label" for="user_1_credit_card">
Credit Card
</label>
<div class="bhojpur-field__block">
<fieldset id="user_1_credit_card" class="bhojpur-fieldset">
  <input id="" class="bhojpur-hidden__primary_key" name="BhojpurResource.CreditCard.ID" value="0" type="hidden">
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <input id="" class="bhojpur-hidden__primary_key" name="BhojpurResource.CreditCard.ID" value="0" type="hidden">
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="">
Number
</label>
<div class="bhojpur-field__show"></div>
<div class="bhojpur-field__edit">
<input class="mdl-textfield__input" type="text" id="" name="BhojpurResource.CreditCard.Number" value="" >
</div>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="">
Issuer
</label>
<div class="bhojpur-field__show"></div>
<div class="bhojpur-field__edit">
<input class="mdl-textfield__input" type="text" id="" name="BhojpurResource.CreditCard.Issuer" value="" >
</div>
</div>
</div>
</div>
</div>
</div>
</fieldset>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
<div class="bhojpur-field collection-edit bhojpur-fieldset-container" >
<label class="bhojpur-field__label" for="user_1_addresses">
Addresses
</label>
<div class="bhojpur-field__block">
<fieldset class="bhojpur-fieldset bhojpur-fieldset--new">
  <button data-confirm="Are you sure?" class="mdl-button bhojpur-button--muted mdl-button--icon mdl-js-button bhojpur-fieldset__delete" type="button">
	<i class="material-icons md-18">delete</i>
  </button>
	<input id="" class="bhojpur-hidden__primary_key" name="BhojpurResource.Addresses[0].ID" value="0" type="hidden">
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <input id="" class="bhojpur-hidden__primary_key" name="BhojpurResource.Addresses[0].ID" value="0" type="hidden">
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="">
Address1
</label>
<div class="bhojpur-field__show"></div>
<div class="bhojpur-field__edit">
<input class="mdl-textfield__input" type="text" id="" name="BhojpurResource.Addresses[0].Address1" value="" >
</div>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
<label class="bhojpur-field__label mdl-textfield__label" for="">
Address2
</label>
<div class="bhojpur-field__show"></div>
<div class="bhojpur-field__edit">
<input class="mdl-textfield__input" type="text" id="" name="BhojpurResource.Addresses[0].Address2" value="" >
</div>
</div>
</div>
</div>
</div>
</div>
</fieldset>
<button class="mdl-button mdl-button--primary bhojpur-fieldset__add" type="button">
  Add Address
</button>
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
<div class="bhojpur-field">
<label class="bhojpur-field__label" for="user_1_company">
Company
</label>
<div class="bhojpur-field__show"></div>
<div class="bhojpur-field__block bhojpur-field__edit  bhojpur-field__selectone" >
<select id="user_1_company" class="bhojpur-field__input hidden"  data-toggle="bhojpur.chooser" data-placeholder="Select an Option" name="BhojpurResource.Company"   data-remote-url="/admin/companies" data-remote-data="true" data-remote-data-primary-key="ID">
</select>
<input type="hidden" name="BhojpurResource.Company" value="">
</div>
</div>
</div>
</div>
</div>
<div class="bhojpur-form-section clearfix" data-section-title="">
<div >
<div class="bhojpur-form-section-rows bhojpur-section-columns-1 clearfix">
  <div class="bhojpur-field">
<label class="bhojpur-field__label" for="user_1_languages">
Languages
</label>
<div class="bhojpur-field__show bhojpur-field__selectmany-show">
</div>
<div class="bhojpur-field__edit bhojpur-field__block bhojpur-field__selectmany"  >
<select class="bhojpur-field__input hidden" id="user_1_languages"  data-toggle="bhojpur.chooser"  data-placeholder="Select some Options" name="BhojpurResource.Languages" multiple  >
</select>
<input type="hidden" name="BhojpurResource.Languages" value="">
</div>
</div>
</div>
</div>
</div>
  <div class="bhojpur-form__actions">
	<button class="mdl-button mdl-button--colored mdl-button--raised mdl-js-button mdl-js-ripple-effect bhojpur-button--save" type="submit">Save Changes</button>
	<a class="mdl-button mdl-button--primary mdl-js-button mdl-js-ripple-effect bhojpur-button--cancel" href="javascript:history.back();">Cancel Edit</a>
  </div>
</form>
`
