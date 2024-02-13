

# Goxios

[![Go Reference](https://pkg.go.dev/badge/github.com/SmikeForYou/goaxios.svg)](https://pkg.go.dev/github.com/SmikeForYou/goaxios)
[![Go Report Card](https://goreportcard.com/badge/github.com/SmikeForYou/goaxios)](https://goreportcard.com/report/github.com/SmikeForYou/goaxios)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
# Goxios

Goaxios is a Golang implementation of the Axios library for JavaScript. It provides a simple and efficient way to make HTTP requests in Go.

## Installation

To install Goaxios, use the following command:

```bash
go get "github.com/SmikeForYou/goxios"
```

## Usage

Here are some examples of how to use Goxios:

### GET Request

```go
package main

import "github.com/SmikeForYou/goxios"

func main() {
	goxios.Get("https://jsonplaceholder.typicode.com/todos/1", nil)
}
```

### GET Request with Params

```go
package main

import "github.com/SmikeForYou/goxios"

func main() {
	goxios.Get("https://jsonplaceholder.typicode.com/todos", &goxios.RequestConfig{
        Params: map[string]interface{}{
            "id": 1,
        },
    })
}
```

### POST Request with Instance

```go
package main

import "github.com/SmikeForYou/goxios"

func main() {
	config := goxios.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
    }
    instance := goxios.NewGoxiosInstance(config)
    jsonRequest := goxios.NewJsonRequest(map[string]interface{}{
        "title": "foo",
        "body": "bar",
        "userId": 1,
    })
    instance.Post("/posts", jsonRequest, nil)
}
```

### POST Request with FormData

```go
package main

import (
	"github.com/SmikeForYou/goxios"
	"os"
)

func main() {
    instance := goxios.NewGoxiosInstance(goxios.Config{
        BaseURL: "https://jsonplaceholder.typicode.com",
    })
    formData := goxios.NewFormDataRequest()
    formData.AddValueStr("title", "foo")
    f, _ := os.Open("avatar.png")
    defer f.Close()
    formData.AddFile("file", f)
    instance.Post("/profile", formData, nil)
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
```
