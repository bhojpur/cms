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
	"encoding/json"
	"errors"
	"testing"

	bhojpurTestUtils "github.com/bhojpur/application/test/utils"
	"github.com/bhojpur/cms/pkg/admin"
	. "github.com/bhojpur/cms/pkg/admin/tests/dummy"
	orm "github.com/bhojpur/orm/pkg/engine"
	"github.com/bhojpur/orm/pkg/now"
	"github.com/theplant/testingutils"
)

func TestJSONTransformerEncode(t *testing.T) {
	bhojpurTestUtils.ResetDBTables(db, &Language{}, &Profile{}, &CreditCard{}, &User{}, &Address{})

	ctx := Admin.NewContext(nil, nil)
	ctx.Context.Roles = []string{Role_system_administrator}

	var (
		buffer          bytes.Buffer
		registeredAt    = now.MustParse("2017-01-01")
		jsonTransformer = &admin.JSONTransformer{}
		encoder         = admin.Encoder{
			Action:   "show",
			Resource: Admin.GetResource("User"),
			Context:  ctx,
			Result: &User{
				Active:       true,
				Model:        orm.Model{ID: 1},
				Name:         "pramila",
				Role:         "admin",
				RegisteredAt: &registeredAt,
				CreditCard: CreditCard{
					Number: "411111111111",
					Issuer: "visa",
				},
				Profile: Profile{
					Name: "pramila",
					Phone: Phone{
						Num: "110",
					},
					Sex: "male",
				},
			},
		}
	)

	if err := jsonTransformer.Encode(&buffer, encoder); err != nil {
		t.Errorf("no error should returned when encode object to JSON")
	}

	var response, expect json.RawMessage
	json.Unmarshal(buffer.Bytes(), &response)

	jsonResponse := `{
        "Active": true,
        "Addresses": [],
        "Age": 0,
        "Avatar": "",
        "Company": "",
        "CreditCard": {
                "ID": 0,
                "Issuer": "visa",
                "Number": "411111111111"
        },
        "ID": 1,
        "Languages": null,
        "Name": "pramila",
        "Profile": {
                "ID": 0,
                "Name": "pramila",
                "Phone": {
                        "ID": 0,
                        "Num": "110"
                },
                "Sex": "male"
        },
        "RegisteredAt": "2017-01-01 00:00",
        "Role": "admin"
}`

	json.Unmarshal([]byte(jsonResponse), &expect)

	diff := testingutils.PrettyJsonDiff(expect, response)
	if len(diff) > 0 {
		t.Errorf("Got %v\n\n\n\n%v", string(buffer.Bytes()), diff)
	}
}

func TestJSONTransformerEncodeMap(t *testing.T) {
	var (
		buffer          bytes.Buffer
		jsonTransformer = &admin.JSONTransformer{}
		encoder         = admin.Encoder{
			Result: map[string]interface{}{"error": []error{errors.New("error1"), errors.New("error2")}},
		}
	)

	jsonTransformer.Encode(&buffer, encoder)

	except := "{\n\t\"error\": [\n\t\t\"error1\",\n\t\t\"error2\"\n\t]\n}"
	if except != buffer.String() {
		t.Errorf("Failed to decode errors map to JSON, except: %v, but got %v", except, buffer.String())
	}
}
