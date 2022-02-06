# Bhojpur CMS - Internationalization

The `I18N` provides internationalization support for your application, it supports two kinds of storages(backends), the database and file system.

## Usage

Initialize `I18N` with the storage mode. You can use both storages together, the earlier one has higher priority. So, in the example, I18N will look up the translation in database first, then continue finding it in the YAML file if not found.

```go
import (
  orm "github.com/bhojpur/orm/pkg/engine"
  "github.com/bhojpur/cms/pkg/i18n"
  "github.com/bhojpur/cms/pkg/i18n/backends/database"
  "github.com/bhojpur/cms/pkg/i18n/backends/yaml"
)

func main() {
  db, _ := orm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")

  I18N := i18n.New(
    database.New(&db), // load translations from the database
    yaml.New(filepath.Join(config.Root, "config/locales")), // load translations from the YAML files in directory `config/locales`
  )

  I18N.T("en-US", "demo.greeting") // Not exist at first
  I18N.T("en-US", "demo.hello") // Exists in the yml file

  i18n.Default = "hi-IN" // change the default locale. the original value is "en-US"
}
```

Once a database has been set for `I18N`, all **untranslated** translations inside `I18N.T()` will be loaded into `translations` table in the database when compiling the application. For example, we have an untranslated `I18N.T("en-US", "demo.greeting")` in the example, so `I18N` will generate this record in the `translations` table after compiling.

| locale | key           | value  |
| ---    | ---           | ---    |
| en-US  | demo.greeting | &nbsp; |

The YAML file format is

```yaml
en-US:
  demo:
    hello: "Hello, world"
```

### Use built-in interface for translation management with [Bhojpur CMS - Admin](http://github.com/bhojpur/cms/pkg/admin)

The `I18N` has a built-in web interface for translations which is integrated with [Bhojpur CMS - Admin](http://github.com/bhojpur/cms/pkg/admin).

```go
Admin.AddResource(I18N)
```

To let users able to translate between locales in the admin interface, your "User" need to implement these interfaces.

```go
func (user User) EditableLocales() []string {
	return []string{"en-US", "hi-IN"}
}

func (user User) ViewableLocales() []string {
	return []string{"en-US", "in-IN"}
}
```

### Use with Go templates

An easy way to use `I18N` in a template is to define a `t` function and register it as `FuncMap`:

```go
func T(key string, value string, args ...interface{}) string {
  return I18N.Default(value).T("en-US", key, args...)
}

// then use it in the template
{{ t "demo.greet" "Hello, {{$1}}" "John" }} // -> Hello, John
```

### Built-in functions for Translations Management

The `I18N` has functions to manage translation directly.

```go
// Add Translation
I18N.AddTranslation(&i18n.Translation{Key: "hello-world", Locale: "en-US", Value: "hello world"})

// Update Translation
I18N.SaveTranslation(&i18n.Translation{Key: "hello-world", Locale: "en-US", Value: "Hello World"})

// Delete Translation
I18N.DeleteTranslation(&i18n.Translation{Key: "hello-world", Locale: "en-US", Value: "Hello World"})
```

### Scope and default value

Call Translation with `Scope` or set default value.

```go
// Read Translation with `Scope`
I18N.Scope("home-page").T("hi-IN", "hello-world") // read translation with translation key `home-page.hello-world`

// Read Translation with `Default Value`
I18N.Default("Default Value").T("hi-IN", "non-existing-key") // Will return default value `Default Value`
```

### Fallbacks

The `I18N` has a `Fallbacks` function to register fallbacks. For example, registering `en-GB` as a fallback to `hi-IN`:

```go
i18n := New(&backend{})
i18n.AddTranslation(&Translation{Key: "hello-world", Locale: "en-GB", Value: "Hello World"})

fmt.Print(i18n.Fallbacks("en-GB").T("hi-IN", "hello-world")) // "Hello World"
```

**To set fallback [*Locale*](https://en.wikipedia.org/wiki/Locale_(computer_software)) globally** you can use `I18N.FallbackLocales`. This function accepts a `map[string][]string` as parameter. The key is the fallback *Locale* and the `[]string` is the *Locales* that could fallback to the first *Locale*.

For example, setting `"fr-FR", "de-DE", "hi-IN"` fallback to `en-GB` globally:

```go
I18N.FallbackLocales = map[string][]string{"en-GB": []{"fr-FR", "de-DE", "hi-IN"}}
```

### Interpolation

I18N utilizes a Go template to parse translations with an interpolation variable.

```go
type User struct {
  Name string
}

I18N.AddTranslation(&i18n.Translation{Key: "hello", Locale: "en-US", Value: "Hello {{.Name}}"})

I18N.T("en-US", "hello", User{Name: "Sanjay"}) //=> Hello Sanjay
```

### Pluralization

I18N utilizes [cldr](https://github.com/theplant/cldr) to achieve pluralization, it provides the functions `p`, `zero`, `one`, `two`, `few`, `many`, `other` for this purpose. Please refer to [cldr documentation](https://github.com/theplant/cldr) for more information.

```go
I18N.AddTranslation(&i18n.Translation{Key: "count", Locale: "en-US", Value: "{{p "Count" (one "{{.Count}} item") (other "{{.Count}} items")}}"})
I18N.T("en-US", "count", map[string]int{"Count": 1}) //=> 1 item
```

### Ordered Params

```go
I18N.AddTranslation(&i18n.Translation{Key: "ordered_params", Locale: "en-US", Value: "{{$1}} {{$2}} {{$1}}"})
I18N.T("en-US", "ordered_params", "string1", "string2") //=> string1 string2 string1
```

### Inline Edit

You could manage translations' data with [Bhojpur CMS - Admin](http://github.com/bhojpur/cms/pkg/admin) interface (UI) after registering it into [Bhojpur CMS -Admin](http://github.com/bhojpur/cms/pkg/admin), however we warn you that it is usually quite hard (and error prone!) to *translate a translation* without knowing its context...Fortunately, the *Inline Edit* feature of [Bhojpur CMS - Admin](http://github.com/bhojpur/cms/pkg/admin) was developed to resolve this problem!

*Inline Edit* allows administrators to manage translations from the frontend. Similarly to [integrating with Go Templates](#integrate-with-golang-templates), you need to register a func map for Go templates to render *inline editable* translations.

The good thing is we have created a package for you to do this easily, it will generate a `FuncMap`, you just need to use it when parsing your templates:

```go
// `I18N` hold translations backends
// `en-US` current locale
// `true` enable inline edit mode or not, if inline edit not enabled, it works just like the funcmap in section "Integrate with Go Templates"
inline_edit.FuncMap(I18N, "en-US", true) // => map[string]interface{}{
                                         //     "t": func(string, ...interface{}) template.HTML {
                                         //        // ...
                                         //      },
                                         //    }
```



## License

Released under the [MIT License](http://opensource.org/licenses/MIT).