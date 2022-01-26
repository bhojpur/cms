package admin

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
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/roles"
	"github.com/bhojpur/application/pkg/utils"
	"github.com/bhojpur/cms/pkg/render/assetfs"
)

var (
	globalViewPaths []string
	globalAssetFSes []assetfs.Interface
	goModDeps       []string
	goModPrefix     = "pkg/mod"
)

// HasPermissioner has permission interface
type HasPermissioner interface {
	HasPermission(roles.PermissionMode, *appsvr.Context) bool
}

// ResourceNamer is an interface for models that defined method `ResourceName`
type ResourceNamer interface {
	ResourceName() string
}

// I18N define admin's i18n interface
type I18N interface {
	Scope(scope string) I18N
	Default(value string) I18N
	T(locale string, key string, args ...interface{}) template.HTML
}

// RegisterViewPath register view path for all assetfs
func RegisterViewPath(pth string) {
	globalViewPaths = append(globalViewPaths, pth)
	var err error
	for _, assetFS := range globalAssetFSes {
		if err = assetFS.RegisterPath(filepath.Join(utils.AppRoot, "vendor", pth)); err != nil {
			for _, gopath := range utils.GOPATH() {
				if err = assetFS.RegisterPath(filepath.Join(gopath, getDepVersionFromMod(pth))); err == nil {
					break
				}

				if err = assetFS.RegisterPath(filepath.Join(gopath, "src", pth)); err == nil {
					break
				}
			}
		}
	}
	if err != nil {
		log.Printf("RegisterViewPathError: %s %s!", pth, err.Error())
	}
}

func equal(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func getDepVersionFromMod(pth string) string {
	if len(goModDeps) == 0 {
		if cont, err := ioutil.ReadFile("go.mod"); err == nil {
			goModDeps = strings.Split(string(cont), "\n")
		}
	}

	for _, val := range goModDeps {
		if txt := strings.Trim(val, "\t\r"); strings.HasPrefix(txt, pth) {
			return filepath.Join(goModPrefix, pth+"@"+strings.Split(txt, " ")[1])
		}
	}

	if strings.LastIndex(pth, "/") == -1 {
		return pth
	}

	return getDepVersionFromMod(pth[:strings.LastIndex(pth, "/")]) + pth[strings.LastIndex(pth, "/"):]
}
