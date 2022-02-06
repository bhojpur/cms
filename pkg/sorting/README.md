# Bhojpur CMS - Sorting

Sorting adds reordering abilities to [Bhojpur ORM](https://github.com/bhojpur/orm) models and sorts collections.

## Register Bhojpur ORM Callbacks

Sorting utilises [Bhojpur ORM](https://github.com/bhojpur/orm) callbacks to log data, so you will need to register callbacks first:

```go
import (
  orm "github.com/bhojpur/orm/pkg/engine"
  "github.com/bhojpur/cms/pkg/sorting"
)

func main() {
  db, err := orm.Open("sqlite3", "demo_db")
  sorting.RegisterCallbacks(db)
}
```

### Sort Modes

Sorting two modes which can be applied as anonymous fields in a model.

- Ascending mode:smallest first (`sorting.Sorting`)
- Descending mode: smallest last (`sorting.SortingDESC`)

They could be used as follows:

```go
// Ascending mode
type Category struct {
  orm.Model
  sorting.Sorting // this will register a `position` column to model Category, used to save record's order
}

db.Find(&categories)
// SELECT * FROM categories ORDER BY position;

// Descending mode
type Product struct {
  orm.Model
  sorting.SortingDESC // this will register a `position` column to model Product, used to save record's order
}

db.Find(&products)
// SELECT * FROM products ORDER BY position DESC;
```

### Reordering

```go
// Move Up
sorting.MoveUp(&db, &product, 1)
// If a record is in positon 5, it will be brought to 4

// Move Down
sorting.MoveDown(&db, &product, 1)
// If a record is in positon 5, it will be brought to 6

// Move To
sorting.MoveTo(&db, &product, 1)
// If a record is in positon 5, it will be brought to 1
```

## Sorting Collections

Sorts a slice of data:

```go
sorter := sorting.SortableCollection{
  PrimaryKeys: []string{"5", "3", "1", "2"}
}

products := []Product{
  {Model: orm.Model{ID: 1}, Code: "1"},
  {Model: orm.Model{ID: 2}, Code: "2"},
  {Model: orm.Model{ID: 3}, Code: "3"},
  {Model: orm.Model{ID: 3}, Code: "4"},
  {Model: orm.Model{ID: 3}, Code: "5"},
}

sorter.Sort(products)

products // => []Product{
         //      {Model: orm.Model{ID: 3}, Code: "5"},
         //      {Model: orm.Model{ID: 3}, Code: "3"},
         //      {Model: orm.Model{ID: 1}, Code: "1"},
         //      {Model: orm.Model{ID: 2}, Code: "2"},
         //      {Model: orm.Model{ID: 3}, Code: "4"},
         //    }
```

### Sorting ORM-backend Models

After enabling sorting modes for [Bhojpur ORM](https://github.com/bhojpur/orm) models, the [Admin](https://github.com/bhojpur/cms/pkg/admin) will automatically enable the sorting feature for the resource.


### Sorting Collections

If you want to make a sortable [select_many](https://docs.bhojpur.net/admin/metas/select-many.html), [collection_edit](http://docs.bhojpur.net/admin/metas/collection-edit.html) field, You could add a `sorting.SortableCollection` field with name: Field's name + 'Sorter'; which is used to save above field's data order. That Field will also be identified as sortable in [Admin](https://github.com/bhojpur/cms/pkg/admin).

```go
// For model relations
type Product struct {
  orm.Model
  l10n.Locale
  Collections           []Collection
  CollectionsSorter     sorting.SortableCollection
  ColorVariations       []ColorVariation `l10n:"sync"`
  ColorVariationsSorter sorting.SortableCollection
}

// For virtual arguments
type selectedProductsArgument struct {
  Products       []string
  ProductsSorter sorting.SortableCollection
}

selectedProductsResource := Admin.NewResource(&selectedProductsArgument{})
selectedProductsResource.Meta(&admin.Meta{Name: "Products", Type: "select_many", Collection: func(value interface{}, context *appsvr.Context) [][]string {
  var collectionValues [][]string
  var products []*models.Product
  db.DB.Find(&products)
  for _, product := range products {
    collectionValues = append(collectionValues, []string{fmt.Sprint(product.ID), product.Name})
  }
  return collectionValues
}})

Widgets.RegisterWidget(&widget.Widget{
  Name:      "Products",
  Templates: []string{"products"},
  Setting:   selectedProductsResource,
}
```

### About record with composite primary key.

It does support sorting records with composite primary key. However, there is a exception, the `version_name` is a "reserved" primary key for the bhojpur/cms/pkg/publish2 support. Please DO NOT use `version_name` as a part of the composite primary key unless you are using bhojpur/cms/pkg/publish2.

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).