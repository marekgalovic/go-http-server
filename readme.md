# Go http server [![Build Status](https://travis-ci.org/marekgalovic/go-http-server.svg?branch=master)](https://travis-ci.org/marekgalovic/go-http-server)

This package provides wrapper around GO's native `net/http` library. It allows you to define parametrized routes and authentication providers while keeping minimal footprint.

## Getting started
Minimal setup is easy. You need to create a `Config` object that allows you to modify properties of your server and then pass the config to a `NewServer` function.
```go
import "github.com/marekgalovic/go-http-server"

func main() {
    config := server.NewConfig()
    s := server.NewServer(config)

    s.Get("/articles", ListArticles, nil)

    s.Listen()
}
```

## Handler functions
As you can see above, every route defnition needs to be associated with a handler function. The handler function accepts only one parameter of type `Request`.

```go
// example handler function
func ListArticles(request *server.Request) *server.Response {
    articles := []*Article{
        &Article{Id: 1, Title: "First"},
        &Article{Id: 2, Title: "Second"},
    }
    return request.Response().Json(articles)
}
```

## Parametrized routes
There are 4 methods that allow you to define your routes. A code example below shows how to define restful endpoinds for a resource. First parameter is a route of the resource, second parameter is a handler function and third parameter is an authentication provider.

```go
// example route definitions
s.Get("/articles", ListArticles, nil)
s.Get("/article/:id", ShowArticle, nil)
s.Post("/article", CreateArticle, nil)
s.Put("/article/:id", UpdateArticle, nil)
s.Delete("/article/:id", UpdateArticle, nil)
```

## Request
Request object exposes attributes about the request as well as defines methods to create responses. Available attributes are `Method (string)`, `Path (string)`, `Params (url.Values)`, `Header (http.Header)`, `RemoteAddr (string)`, `Body (io.ReadCloser)`. There is also a bunch of methods that you can use in your handlers.

`Get` method to access request parameters.
```go
// example path definition /article/:id
// request path /article/5?foo=bar

id := request.Get(":id")
foo := request.Get("foo")
```

`Json` method to read JSON encoded body to a struct.
```go
var article Article
request.Json(&article)
```

`SetCookie` & `GetCookie` helpers to modify stored cookies.
```go
request.GetCookie("name")
request.SetCookie("name", "value", 60 * time.Minute)
```

`Response` method returns a response`(*Response)` object associated with this request.
```go
request.Response()
```

## Response

`Plain` method to respond with plain-text body.
```go
request.Response().Plain("Article id %d", 1)
```

`Json` method to write JSON responses.
```go
request.Response().Json(map[string]string{"foo": "bar"})
```

`Error` to write error responses with specific code.
```go
request.Response().Error(404, "Resource id: %d not found", 2)
```

`ErrorJson` to write error responses with JSON body.
```go
request.Response().Error(500, map[string]string{"message": "Unable to connect to database"})
```

`File` to respond with a file. File path should be relative to `StaticRoot` defined in server config.
```go
request.Response().File("/path/to/my_file.pdf")
```

`Redirect`
```go
request.Response().Redirect(301, "http://google.com")
```

`SetCode`
```go
request.Response().SetCode(500)
```

`SetHeader` to set a response headers.
```go
request.Response().SetHeader("Keep-Alive", "timeout=5")
```
*Response methods support method chaining so you can use `request.Response().SetCode(404).Plain("Resource not found")`*

## Authentication providers
You can create authentication provider to create authentication strategy that best fits your needs. There is only one method defined by `AuthProvider` interface called `Verify`. This method accepts `Request` object as a parameter and returns `Response` object if authentication fails and `nil` if authentication was successful.

Define an authentication provider
```go
type MyAuthProvider struct {}

func (auth *MyAuthProvider) Verify(request *server.Request) *server.Response {
    if request.Get("auth_token") == "secret_token" {
        return nil
    }
    return request.Response().Error(401, "Authentication failed")
}
```
Use the provider to protect your routes
```go
auth := &MyAuthProvider{}

s.Get("/articles", ListArticles, nil) // public route
s.Post("/article", CreateArticle, auth) // protected route
```

## SSL support
To configure a secure server you need to provide a path to your certificate and key.
```go
config.CertFile = "server.crt"
config.KeyFile = "server.key"
```

## Contributing
Bug reports and pull requests are welcome on GitHub at https://github.com/marekgalovic/go-http-server. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.


## License
The package is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
