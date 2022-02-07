# Bhojpur CMS - Administrator Dashboard

You can instantly create a beautiful, cross platform, configurable System Administrator's web user interface and Backend APIs for managing data in minutes.

## Core Features

- Generate [Bhojpur CMS - Admin](https://github.com/bhojpur/cms/pkg/admin) interface for managing data
- RESTful JSON APIs
- Association handling
- Search and filtering
- Actions/Batch Actions
- Authentication and Authorization
- Extendability

## Quick Start

```go
package main

import (
  "fmt"
  "net/http"
  orm "github.com/bhojpur/orm/pkg/engine"
  _ "github.com/mattn/go-sqlite3"
  "github.com/bhojpur/cms/pkg/admin"
)

// Create a simple Bhojpur ORM backend model
type User struct {
  orm.Model
  Name string
}

// Create another Bhojpur ORM backend model
type Product struct {
  orm.Model
  Name        string
  Description string
}

func main() {
  demoDB, _ := orm.Open("sqlite3", "internal/demo.db")
  demoDB.AutoMigrate(&User{}, &Product{})

  // Initialize
  Admin := admin.New(&admin.AdminConfig{DB: demoDB})

  // Allow to use Admin module to manage User, Product
  Admin.AddResource(&User{})
  Admin.AddResource(&Product{})

  // initialize an HTTP service request multiplexer
  mux := http.NewServeMux()

  // Mount Administrator's web user interface to mux
  Admin.MountTo("/admin", mux)

  fmt.Println("Listening on: 3000")
  http.ListenAndServe(":3000", mux)
}
```

`go run main.go` and visit `localhost:3000/admin` to see the result!

## How to use remoteSelector with publish2.version integrated record

### **For many relationship**

Suppose we have two models Factory and Item. Factory **has many** Items.

In the struct, you need add a field `resource.CompositePrimaryKeyField` to the "many" side, which is `Item` here.

```go
type Factory struct {
	orm.Model
	Name string

	publish2.Version
	Items       []Item `orm:"many2many:factory_items;association_autoupdate:false"`
	ItemsSorter sorting.SortableCollection
}

type Item struct {
	orm.Model
	Name string
	publish2.Version

	// github.com/bhojpur/application/pkg/resource
	resource.CompositePrimaryKeyField // Required
}
```

then, define a remote resource selector. You need configure the `ID` meta like below to make it support composite primary key, this is mandatory.

```go
func generateRemoteItemSelector(adm *admin.Admin) (res *admin.Resource) {
	res = adm.AddResource(&Item{}, &admin.Config{Name: "ItemSelector"})
	res.IndexAttrs("ID", "Name")

	// Required. Convert single ID into composite primary key
	res.Meta(&admin.Meta{
	Name: "ID",
	Valuer: func(value interface{}, ctx *appsvr.Context) interface{} {
		if r, ok := value.(*Item); ok {
			// github.com/bhojpur/application/pkg/resource
			return resource.GenCompositePrimaryKey(r.ID, r.GetVersionName())
		}
		return ""
	},
	})

	return res
}
```

Last, use it in the Factory resource.

```go
itemSelector := generateRemoteItemSelector(adm)
factoryRes.Meta(&admin.Meta{
	Name: "Items",
	Config: &admin.SelectManyConfig{
	RemoteDataResource: itemSelector,
	},
})
```

### **For single relationship**

Suppose we have two models. Factory and Manager. Factory **has one** Manager.

Firstly, in the struct, you need add a field `resource.CompositePrimaryKeyField` to the "one" side, which is `Manager` here.

```go
type Factory struct {
	orm.Model
	Name string
	publish2.Version

	ManagerID          uint
	ManagerVersionName string // Required. in "xxxVersionName" format.
	Manager            Manager
}

type Manager struct {
	orm.Model
	Name string
	publish2.Version

	// github.com/bhojpur/application/pkg/resource
	resource.CompositePrimaryKeyField // Required
}
```

then, define a remote resource selector. You need configure the `ID` meta like below to make it support composite primary key, this is mandatory.

```go
func generateRemoteManagerSelector(adm *admin.Admin) (res *admin.Resource) {
	res = adm.AddResource(&Manager{}, &admin.Config{Name: "ManagerSelector"})
	res.IndexAttrs("ID", "Name")

	// Required. Convert single ID into composite primary key
	res.Meta(&admin.Meta{
		Name: "ID",
		Valuer: func(value interface{}, ctx *appsvr.Context) interface{} {
			if r, ok := value.(*Manager); ok {
				// github.com/bhojpur/application/pkg/resource
				return resource.GenCompositePrimaryKey(r.ID, r.GetVersionName())
			}
			return ""
		},
	})

	return res
}

Lastly, use it in the Factory resource.

```go
managerSelector := generateRemoteManagerSelector(adm)
factoryRes.Meta(&admin.Meta{
	Name: "Manager",
	Config: &admin.SelectOneConfig{
		RemoteDataResource: managerSelector,
	},
})
```

If you need to overwrite Collection. you have to pass composite primary key as the first element of the returning array instead of ID.

```go
factoryRes.Meta(&admin.Meta{
  Name: "Items",
  Config: &admin.SelectManyConfig{
	Collection: func(value interface{}, ctx *appsvr.Context) (results [][]string) {
		if c, ok := value.(*Factory); ok {
		var items []Item
		ctx.GetDB().Model(c).Related(&items, "Items")

		for _, p := range items {
		// The first element must be the composite primary key instead of ID
		results = append(results, []string{resource.GenCompositePrimaryKey(p.ID, p.GetVersionName()), p.Name})
		}
		}
		return
	},
	RemoteDataResource: itemSelector,
  },
})
```

## To support assign associations while creating a new version

If you want to assign associations while creating a new version of object immediately. You need to define a function called `AssignVersionName` to the versioned struct with **pointer** receiver which should contains the generating new version name's logic and assign the new version name to the object. For example

```go
func (fac *Factory) AssignVersionName(db *orm.DB) {
	var count int
	name := time.Now().Format("2006-01-02")
	if err := db.Model(&CollectionWithVersion{}).Where("id = ? AND version_name like ?", fac.ID, name+"%").Count(&count).Error; err != nil {
    panic(err)
  }
	fac.VersionName = fmt.Sprintf("%s-v%v", name, count+1)
}
```

## Documentation

To print all registered routes

```go
// adm is a Bhojpur CMS admin instance
adm.GetRouter().PrintRoutes()
```

## ViewPath Note

Bhojpur CMS still support Go path while finding its template files. The priority is

1. check vendor, if not found
2. check $GOPATH/pkg/mod/github.com/bhojpur/cms@v0.x/pkg/admin/views. The version would be detected automatically by your go.mod file, if still not found
3. load view path from $GOPATH/src/github.com/bhojpur/cms/pkg/admin/views


So, if you want to use the template under the pkg/mod, make sure $GOPATH/src/github.com/bhojpur/cms/pkg/admin is absent.

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).