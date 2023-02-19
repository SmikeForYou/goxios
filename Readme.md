# Goaxios

Goaxios is a Golang implementation of the Axios library for JavaScript.

## Installation

```bash
go get "github.com/SmikeForYou/goaxios"
```

## Usage

```go
package main

import "github.com/SmikeForYou/goaxios"

// GET
goaxios.Get("https://jsonplaceholder.typicode.com/todos/1")

// GET with params
goaxios.Get("https://jsonplaceholder.typicode.com/todos", &goaxios.RequestConfig{
    Params: map[string]interface{}{
        "id": 1,
    },
})
}) // https://jsonplaceholder.typicode.com/todos?id=1

// POST with instance
instance := goaxios.New(goaxios.Config{
    BaseURL: "https://jsonplaceholder.typicode.com",
})
jsonRequest := NewJsonRequest(map[string]interface{}{
    "title": "foo",
    "body": "bar",
    "userId": 1,
})
instance.Post("/posts", jsonRequest) // POST https://jsonplaceholder.typicode.com/posts with json body

// POST FormData
formData := NewFormDataRequest()
formData.Add("title", "foo")
f, _ := os.Open("avatar.png")
defer f.Close()
formData.Add("file", file)
instance.Post("/profile", formData, nil) // POST https://jsonplaceholder.typicode.com/posts with form data
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)