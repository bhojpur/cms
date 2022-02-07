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
	"log"
	"path/filepath"
	"reflect"

	appsvr "github.com/bhojpur/application/pkg/engine"
	"github.com/bhojpur/application/pkg/resource"
	"github.com/bhojpur/application/pkg/utils"
	"github.com/bhojpur/cms/pkg/render/assetfs"
	"github.com/bhojpur/cms/pkg/session"
	"github.com/bhojpur/cms/pkg/session/manager"
	orm "github.com/bhojpur/orm/pkg/engine"
	"github.com/bhojpur/orm/pkg/inflection"
	"github.com/theplant/cldr"
)

// AdminConfig admin config struct
type AdminConfig struct {
	// SiteName set site's name, the name will be used as admin HTML title and admin interface will auto load javascripts, stylesheets files based on its value
	SiteName        string
	DB              *orm.DB
	Auth            Auth
	AssetFS         assetfs.Interface
	SessionManager  session.ManagerInterface
	SettingsStorage SettingsStorageInterface
	I18N            I18N
	*Transformer
}

// Admin is a struct that used to generate admin/api interface
type Admin struct {
	*AdminConfig
	menus             []*Menu
	resources         []*Resource
	searchResources   []*Resource
	router            *Router
	funcMaps          template.FuncMap
	metaConfigureMaps map[string]func(*Meta)
}

// New new admin with configuration
func New(config interface{}) *Admin {
	admin := Admin{
		funcMaps:          make(template.FuncMap),
		router:            newRouter(),
		metaConfigureMaps: defaultMetaConfigureMaps,
	}

	if c, ok := config.(*appsvr.Config); ok {
		admin.AdminConfig = &AdminConfig{DB: c.DB}
	} else if c, ok := config.(*AdminConfig); ok {
		admin.AdminConfig = c
	} else {
		admin.AdminConfig = &AdminConfig{}
	}

	if admin.SessionManager == nil {
		admin.SessionManager = manager.SessionManager
	}

	if admin.Transformer == nil {
		admin.Transformer = DefaultTransformer
	}

	if admin.AssetFS == nil {
		admin.AssetFS = assetfs.AssetFS().NameSpace("admin")
	}

	if admin.SettingsStorage == nil {
		admin.SettingsStorage = newSettings(admin.AdminConfig.DB)
	}

	admin.SetAssetFS(admin.AssetFS)

	if admin.AdminConfig.DB != nil {
		admin.AdminConfig.DB.AutoMigrate(&BhojpurAdminSetting{})
	}

	admin.registerCompositePrimaryKeyCallback()
	return &admin
}

// SetSiteName set site's name, the name will be used as admin HTML title and admin interface will auto load javascripts, stylesheets files based on its value
// For example, if you named it as `Bhojpur CMS Demo`, admin will look up `bhojpur_demo.js`, `bhojpur_demo.css` in Bhojpur CMS view paths, and load them if found
func (admin *Admin) SetSiteName(siteName string) {
	admin.SiteName = siteName
}

// SetAuth set admin's authorization gateway
func (admin *Admin) SetAuth(auth Auth) {
	admin.Auth = auth
}

// SetAssetFS set AssetFS for admin
func (admin *Admin) SetAssetFS(assetFS assetfs.Interface) {
	admin.AssetFS = assetFS
	globalAssetFSes = append(globalAssetFSes, assetFS)

	admin.AssetFS.RegisterPath(filepath.Join(utils.AppRoot, "pkg/admin/views"))
	admin.RegisterViewPath("pkg/admin/views")

	for _, viewPath := range globalViewPaths {
		admin.RegisterViewPath(viewPath)
	}
}

// RegisterViewPath register view path for admin
func (admin *Admin) RegisterViewPath(pth string) {
	var err error
	if err = admin.AssetFS.RegisterPath(filepath.Join(utils.AppRoot, "vendor", pth)); err != nil {
		for _, gopath := range utils.GOPATH() {
			if err = admin.AssetFS.RegisterPath(filepath.Join(gopath, getDepVersionFromMod(pth))); err == nil {
				break
			}

			if err = admin.AssetFS.RegisterPath(filepath.Join(gopath, "src", pth)); err == nil {
				break
			}
		}
	}
	if err != nil {
		log.Printf("RegisterViewPathError: %s %s!", pth, err.Error())
	}
}

// RegisterMetaConfigure register configure for a kind, it will be called when register those kind of metas
func (admin *Admin) RegisterMetaConfigure(kind string, fc func(*Meta)) {
	admin.metaConfigureMaps[kind] = fc
}

// RegisterFuncMap register view funcs, it could be used in view templates
func (admin *Admin) RegisterFuncMap(name string, fc interface{}) {
	admin.funcMaps[name] = fc
}

// GetRouter get router from admin
func (admin *Admin) GetRouter() *Router {
	return admin.router
}

func (admin *Admin) newResource(value interface{}, config ...*Config) *Resource {
	var configuration *Config
	if len(config) > 0 {
		configuration = config[0]
	}

	if configuration == nil {
		configuration = &Config{}
	}

	res := &Resource{
		Resource: resource.New(value),
		Config:   configuration,
		admin:    admin,
	}

	res.Permission = configuration.Permission

	if configuration.Name != "" {
		res.Name = configuration.Name
	} else if namer, ok := value.(ResourceNamer); ok {
		res.Name = namer.ResourceName()
	}

	// Configure resource when initializing
	modelType := utils.ModelType(res.Value)
	for i := 0; i < modelType.NumField(); i++ {
		if fieldStruct := modelType.Field(i); fieldStruct.Anonymous {
			if injector, ok := reflect.New(fieldStruct.Type).Interface().(resource.ConfigureResourceBeforeInitializeInterface); ok {
				injector.ConfigureBhojpurResourceBeforeInitialize(res)
			}
		}
	}

	if injector, ok := res.Value.(resource.ConfigureResourceBeforeInitializeInterface); ok {
		injector.ConfigureBhojpurResourceBeforeInitialize(res)
	}

	findOneHandler := res.FindOneHandler
	res.FindOneHandler = func(result interface{}, metaValues *resource.MetaValues, context *appsvr.Context) error {
		if context.ResourceID == "" {
			context.ResourceID = res.GetPrimaryValue(context.Request)
		}
		return findOneHandler(result, metaValues, context)
	}

	res.UseTheme("slideout")
	return res
}

// NewResource initialize a new Bhojpur CMS resource, won't add it to admin, just initialize it
func (admin *Admin) NewResource(value interface{}, config ...*Config) *Resource {
	res := admin.newResource(value, config...)
	res.Config.Invisible = true
	res.configure()
	return res
}

// AddResource make a model manageable from admin interface
func (admin *Admin) AddResource(value interface{}, config ...*Config) *Resource {
	res := admin.newResource(value, config...)
	admin.resources = append(admin.resources, res)

	res.configure()

	if !res.Config.Invisible {
		res.Action(&Action{
			Name:   "Delete",
			Method: "DELETE",
			URL: func(record interface{}, context *Context) string {
				return context.URLFor(record, res)
			},
			Permission: res.Config.Permission,
			Modes:      []string{"menu_item"},
		})

		menuName := res.Name
		if !res.Config.Singleton {
			menuName = inflection.Plural(res.Name)
		}
		admin.AddMenu(&Menu{Name: menuName, IconName: res.Config.IconName, Permissioner: res, Priority: res.Config.Priority, Ancestors: res.Config.Menu, RelativePath: res.ToParam()})

		admin.RegisterResourceRouters(res, "create", "update", "read", "delete")
	}

	return res
}

// GetResources get defined resources from admin
func (admin *Admin) GetResources() []*Resource {
	return admin.resources
}

// GetResource get resource with name
func (admin *Admin) GetResource(name string) (resource *Resource) {
	for _, res := range admin.resources {
		modelType := utils.ModelType(res.Value)
		// find with defined name first
		if res.ToParam() == name || res.Name == name || modelType.String() == name {
			return res
		}

		// if failed to find, use its model name
		if modelType.Name() == name {
			resource = res
		}
	}

	return
}

// AddSearchResource make a resource searchable from search center
func (admin *Admin) AddSearchResource(resources ...*Resource) {
	admin.searchResources = append(admin.searchResources, resources...)
}

// T call i18n backend to translate
func (admin *Admin) T(context *appsvr.Context, key string, value string, values ...interface{}) template.HTML {
	locale := utils.GetLocale(context)

	if admin.I18N == nil {
		if result, err := cldr.Parse(locale, value, values...); err == nil {
			return template.HTML(result)
		}
		return template.HTML(key)
	}

	return admin.I18N.Default(value).T(locale, key, values...)
}
