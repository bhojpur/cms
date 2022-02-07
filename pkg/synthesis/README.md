# Bhojpur CMS - Synthesis Framework

The web application `synthesis` framework could be used to compile template files into a binary.

## Install [Synthesis](https://github.com/bhojpur/cms/pkg/synthesis)

```sh
$ go get -u -f github.com/bhojpur/cms/pkg/synthesis/...
```

Firstly, initialize your web application `synthesis` library for the project, set the PATH you want to store related template files. For example, `internal/demo`:

```sh
$ cmsctl template internal/demo
```

## Usage

```go
import "<your_project>/internal/demo"

func main() {
  demoFS := demo.AssetFS

  // Register view paths into demoFS
  demoFS.RegisterPath("<your_project>/app/views")
  demoFS.RegisterPath("<your_project>/vendor/plugin/views")

  // Compile application templates under registered view paths into binary
  demoFS.Compile()

  // Get file content with registered name
  fileContent, ok := demoFS.Asset("home/index.tmpl")
}
```

You need to compile application templates into a Go file with method `Compile` before run `go build`, and if any application templates changed, you need to recompile it again.

If you started the web application with a tag `demofs`, then demoFS will access file from generated Go source files or current binary

```sh
go run -tags 'demofs' main.go
```

Else, it will access the content from registered view paths of your filesystem, which is easier for your local development

```sh
go run main.go
```

### Using NameSpace

Although you could initalize several `application` packages to hold template files from different view paths (templates name might have conflicts) for different use cases, therefore the web application `synthesis` library provides an easier solution to you.

```go
func main() {
  // Generate a sub demoFS with a namespace
  adminAssetFS := demoFS.NameSpace("admin_related_files")

  // Register view paths into this name space
  adminAssetFS.RegisterPath("<your_project>/app/admin_views")

  // Access files that registered under views paths of current name space
  adminAssetFS.Asset("admin_view.tmpl")
}
```

### Use Synthesis with [Administrator Dashboard](https://github.com/bhojpur/cms/pkg/admin)

```go
import "<your_project>/internal/demo"

func main() {
  Admin = admin.New(&appsvr.Config{DB: db.Publish.DraftDB()})
  Admin.SetAssetFS(demo.AssetFS.NameSpace("admin"))
}
```

### Use Synthesis with [Render Framework](https://github.com/bhojpur/cms/pkg/render)

```go
import  "github.com/bhojpur/cms/pkg/render"

func main() {
  View := render.New()
  View.SetAssetFS(demoFS.NameSpace("views"))
}
```

### Use Synthesis with [Widget Framework](https://github.com/bhojpur/cms/pkg/widget)

```go
import  "github.com/bhojpur/cms/pkg/widget"

func main() {
  Widgets := widget.New(&widget.Config{DB: db.DB})
	Widgets.SetAssetFS(demoFS.NameSpace("widgets"))
}
```

### Use Synthesis with static files

```go
func main() {
  mux := http.NewServeMux()

  // this will add all files under public into a generated Go source file, which will be included into the binary
  demoFS := demoFS.FileServer(http.Dir("public"))

  // If you only want to include specified paths, you could use it like this
  demoFS := demoFS.FileServer(http.Dir("public"), "javascripts", "stylesheets", "images")

  for _, path := range []string{"javascripts", "stylesheets", "images"} {
    mux.Handle(fmt.Sprintf("/%s/", path), demoFS)
  }
}
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).