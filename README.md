[![wercker status](https://app.wercker.com/status/1438db57ab5937e1fd92d45eca974257/s/master "wercker status")](https://app.wercker.com/project/bykey/1438db57ab5937e1fd92d45eca974257)
[![Coverage Status](https://coveralls.io/repos/jadengore/go-json-rest-middleware-force-ssl/badge.svg?branch=HEAD&service=github)](https://coveralls.io/github/jadengore/go-json-rest-middleware-force-ssl?branch=HEAD)
[![GoDoc](https://godoc.org/github.com/jadengore/go-json-rest-middleware-force-ssl?status.svg)](https://godoc.org/github.com/jadengore/go-json-rest-middleware-force-ssl)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/jadengore/go-json-rest-middleware-force-ssl/master/LICENSE)

# Force SSL Middleware for go-json-rest
Middleware to force SSL on requests to a `go-json-rest` API.

## Installation

```sh
go get github.com/jadengore/go-json-rest-middleware-force-ssl
```

## Example Usage

```go
package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/jadengore/go-json-rest-middleware-force-ssl"
    "log"
    "net/http"
)

func main() {
    api := rest.NewApi()
    api.Use(forceSSL.Middleware{}) // struct with options
    api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request) {
        w.WriteJson(map[string]string{"body": "Hello World!"})
    }))
    log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
```

## Options

| Option                 | Type   | Description | Defaults to
|------------------------|--------|-------------|
| **TrustXFPHeader**     | bool   | Trust `X-Forwarded-Proto` headers (this can allow a client to spoof whether they were using HTTPS) | false
| **Enable301Redirects** | bool   | Enables 301 redirects to the HTTPS version of the request. | false
| **Message**            | String | Allows a custom response message when forcing SSL without redirect. | `SSL Required.`

## Middleware Options Example

```go
api.Use(forceSSL.Middleware{
  TrustXFPHeader: true,
  Enable301Redirects: true,
  Message: "We are unable to process your request over HTTP."
})
```
