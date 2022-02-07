# Bhojpur CMS - VIPs Handler

## Dev Environment

```brew install vips```

export CGO_CFLAGS_ALLOW="-Xpreprocessor"


## Build Dockerfile

```
FROM alpine:3.12

RUN apk add --update go gcc g++ git

RUN apk add --update build-base vips-dev
```
## Build Command

set CGO_ENABLED=1, eg:
```
GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -tags 'template' -a -o main main.go
```

## Deploy Dockerfile

```
FROM alpine:3.12

RUN apk --update upgrade && \

    apk add ca-certificates && \
    
    apk add tzdata && \
    
    apk add build-base vips-dev && \
    
    rm -rf /var/cache/apk/*
```
 
## Simple Usage

[Setup media library](https://github.com/bhojpur/cms/pkg/media#how-to-setup-a-media-library-and-use-media-box) and add below code, then it will compress jpg/png and generate webp for you.

```
import "github.com/bhojpur/cms/pkg/media/handlers/vips"

vips.UseVips(vips.Config{EnableGenerateWebp: true})
```

You can adjust the image quality by config if you want.
```
type Config struct {
	EnableGenerateWebp bool
	PNGtoWebpQuality   int
	JPEGtoWebpQuality  int
	JPEGQuality        int
	PNGQuality         int
	PNGCompression     int
}
  ```  
