package dummy

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

	"github.com/bhojpur/cms/pkg/media/oss"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type CreditCard struct {
	orm.Model
	Number string
	Issuer string
}

type Company struct {
	orm.Model
	Name string
}

type Address struct {
	orm.Model
	UserID   uint
	Address1 string
	Address2 string
}

type Language struct {
	orm.Model
	Name string
}

type User struct {
	orm.Model
	Name         string `orm:"size:50"`
	Age          uint
	Role         string
	Active       bool
	RegisteredAt *time.Time
	Avatar       oss.OSS
	Profile      Profile // has one
	CreditCardID uint
	CreditCard   CreditCard // belongs to
	Addresses    []Address  // has many
	CompanyID    uint
	Company      *Company   // belongs to
	Languages    []Language `orm:"many2many:user_languages;"` // many 2 many
}

type Profile struct {
	orm.Model
	UserID uint
	Name   string
	Sex    string

	Phone Phone
}

type Phone struct {
	orm.Model

	ProfileID uint64
	Num       string
}

func (u User) DisplayName() string {
	return u.Name
}

func (u User) GetID() uint {
	return u.ID
}

func (u User) GetUsersByIDs(db *orm.DB, ids []string) interface{} {
	return u
}
