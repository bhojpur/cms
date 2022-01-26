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
	"encoding/xml"
	"errors"
	"testing"

	"github.com/bhojpur/cms/pkg/admin"
)

func TestXMLTransformerEncode(t *testing.T) {
	t.Skip()
	xmlResult := admin.XMLStruct{
		Result: map[string]interface{}{"error": errors.New("error message"), "status": map[string]int{"code": 200}},
	}
	result := "<response>\n\t<error>error message</error>\n\t<status>\n\t\t<code>200</code>\n\t</status>\n</response>"

	if xmlMarshalResult, err := xml.MarshalIndent(xmlResult, "", "\t"); err != nil {
		t.Errorf("no error should happen, but got %v", err)
	} else if string(xmlMarshalResult) != result {
		t.Errorf("Generated XML got %v, but should be %v", string(xmlMarshalResult), result)
	}
}
