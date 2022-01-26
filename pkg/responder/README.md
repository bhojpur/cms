# Bhojpur CMS - Responder

The Responder provides a means to respond differently according to a request's accepted MIME type.

## Usage

### Register MIME type

```go
import "github.com/bhojpur/cms/pkg/responder"

responder.Register("text/html", "html")
responder.Register("application/json", "json")
responder.Register("application/xml", "xml")
```

The [Responder](https://github.com/bhojpur/cms/pkg/responder) has the above mentioned mime types registered by default. You can register more types with the `Register` function, which accepts two parameters:

1. The mime type, like `text/html`
2. The format of the mime type, like `html`

### Respond to registered mime types

```go
func handler(writer http.ResponseWriter, request *http.Request) {
  responder.With("html", func() {
    writer.Write([]byte("this is a html request"))
  }).With([]string{"json", "xml"}, func() {
    writer.Write([]byte("this is a json or xml request"))
  }).Respond(request)
})
```

The first `html` in the example will be the default response type if [Responder](https://github.com/bhojpur/cms/pkg/responder) cannot find a corresponding MIME type.

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).