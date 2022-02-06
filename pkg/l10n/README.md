# Bhojpur CMS - Localization Methods

The L10N provides your [Bhojpur ORM](https://github.com/bhojpur/orm) models an ability to localize for different natural languages (i.e. Locales). It could be a catalyst for the adaptation of a product, service, application, or document content to meet the human language, cultural, and other requirements of a specific target market.

## Usage

The L10N package utilizes [Bhojpur ORM](https://github.com/bhojpur/orm) callbacks to handle localization, so you will need to register the callbacks first:

```go
import (
  orm "github.com/bhojpur/orm/pkg/engine"
  l10n "github.com/bhojpur/cms/pkg/l10n"
)

func main() {
  db, err := orm.Open("sqlite3", "demo_db")
  l10n.RegisterCallbacks(&db)
}
```

### Making a Model localizable

Embed `l10n.Locale` into your model as an anonymous field to enable localization. For example, in a hypothetical project which has a focus on Product Management:

```go
type Product struct {
  orm.Model
  Name string
  Code string
  l10n.Locale
}
```

`l10n.Locale` will add a `language_code` column as a composite primary key with existing primary keys, using [Bhojpur ORM](https://github.com/bhojpur/orm)'s AutoMigrate to create the field.

The `language_code` column will be used to save a localized model's Locale. If no Locale is set, then the global default Locale (`en-US`) will be used. You can override the global default Locale by setting `l10n.Global`. For example:

```go
l10n.Global = 'hi-IN'
```

### Create localized resources from global product

```go
// Create global product
product := Product{Name: "Global product", Description: "Global product description"}
DB.Create(&product)
product.LanguageCode   // "en-US"

// Create hi-IN product
product.Name = "ठेकुआ"
DB.Set("l10n:locale", "hi-IN").Create(&product)

// Query hi-IN product with primary key 111
DB.Set("l10n:locale", "hi-IN").First(&productIN, 111)
productIN.Name         // "ठेकुआ"
productIN.LanguageCode // "hi"
```

#### Create localized resource directly

By default, only global data is allowed to be created, local data have to localized from global one.

If you want to allow user create localized data directly, you can embeded `l10n.LocaleCreatable` for your model/struct, e.g:

```go
type Product struct {
  orm.Model
  Name string
  Code string
  l10n.LocaleCreatable
}
```

### Keeping localized resources' fields in sync

Add the tag `l10n:"sync"` to the fields that you wish to always sync with the *global* record:

```go
type Product struct {
  orm.Model
  Name  string
  Code  string `l10n:"sync"`
  l10n.Locale
}
```

Now, the localized product's `Code` will be the same as the global product's `Code`. The `Code` is not affected by localized resources, and when the global record changes its `Code` the localized records' `Code` will be synced automatically.

### Query Modes

The L10N provides five modes for querying.

* global   - find all global records,
* locale   - find localized records,
* reverse  - find global records that haven't been localized,
* unscoped - raw query, won't auto add `locale` conditions when querying,
* default  - find localized record, if not found, return the global one.

You can specify the mode in this way:

```go
dbIN := db.Set("l10n:locale", "hi-IN")

mode := "global"
dbIN.Set("l10n:mode", mode).First(&product, 111)
// SELECT * FROM products WHERE id = 111 AND language_code = 'en-US';

mode := "locale"
db.Set("l10n:mode", mode).First(&product, 111)
// SELECT * FROM products WHERE id = 111 AND language_code = 'hi-IN';
```

## Application Integration

Although the L10N could be used alone, it integrates nicely with [Bhojpur Application](https://github.com/bhojpur/application).

By default, [Bhojpur Application](https://github.com/bhojpur/application) will only allow you to manage the global language. If you have configured Authentication, [Bhojpur CMS - Administrator](https://github.com/bhojpur/cms/pkg/admin) will try to obtain the allowed Locales from the current user.

* Viewable Locales - Locales for which the current user has read permission:

```go
func (user User) ViewableLocales() []string {
  return []string{l10n.Global, "hi-IN", "JP", "EN", "DE"}
}
```

* <a name='editable-locales'></a> Editable Locales - Locales for which the current user has manage (create/update/delete) permission:

```go
func (user User) EditableLocales() []string {
  if user.role == "global_admin" {
    return []string{l10n.Global, "hi-IN", "EN"}
  } else {
    return []string{"hi-IN", "EN"}
  }
}
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).