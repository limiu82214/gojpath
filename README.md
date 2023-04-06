# gojpath

[![Go Reference](https://pkg.go.dev/badge/github.com/limiu82214/gojpath.svg)](https://pkg.go.dev/github.com/limiu82214/gojpath) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![codecov](https://codecov.io/gh/limiu82214/gojpath/branch/master/graph/badge.svg?token=0XAK9BB5WL)](https://codecov.io/gh/limiu82214/gojpath) [![Go Report Card](https://goreportcard.com/badge/github.com/limiu82214/gojpath)](https://goreportcard.com/report/github.com/limiu82214/gojpath) ![github actions workflow](https://github.com/limiu82214/gojpath/actions/workflows/go.yml/badge.svg)

`gojpath` is a language for querying JSON data that is similar to XPath. In Golang, you can use the `Get` function to query JSON data.  
The extent of support for this function can be found in this [link](https://learn.microsoft.com/en-us/azure/data-explorer/kusto/query/jsonpath).  

## Usage

Here is an example code that uses the `Get` function for JSON Path queries:

```go
package main

import (
    "encoding/json"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
    jsonString := `{
        "store": {
            "book": [
                {
                    "title": "Harry Potter and the Philosopher's Stone",
                    "author": "J.K. Rowling",
                    "price": 7.99
                },
                {
                    "title": "Harry Potter and the Chamber of Secrets",
                    "author": "J.K. Rowling",
                    "price": 9.99
                }
            ],
            "bicycle": {
                "color": "red",
                "price": 19.95
            }
        }
    }`

    var jsonData interface{}

    err := json.Unmarshal([]byte(jsonString), &jsonData)
    if err != nil {
        t.Fatal(err)
    }

    // Query with JSON Path
    result, _ := Get(jsonData, "$.store.book[0].title")
    assert.Equal(t, "Harry Potter and the Philosopher's Stone", result)

    result, _ = Get(jsonData, "$['store']['book'][0]['title']")
    assert.Equal(t, "Harry Potter and the Philosopher's Stone", result)

    result, _ = Get(jsonData, "$.store.bicycle.color")
    assert.Equal(t, "red", result)

    result, _ = Get(jsonData, "$.store.book[1].author")
    assert.Equal(t, "J.K. Rowling", result)

    result, _ = Get(jsonData, "$.store['book'][1].price")
    assert.Equal(t, 9.99, result)
}
```

In this example, the TestGet function uses the Get function for multiple JSON Path queries.  
By comparing the returned values with the expected values, you can confirm the accuracy of the queries.

## Considerations

As the extent of support for the Get function is limited and some operations, such as the * operator, are not supported, it is recommended to check if the operator and syntax used in JSON Path queries are supported by the Get function to avoid generating incorrect results.

## Other

If you encounter any issues during use, please feel free to raise an issue on the GitHub project or contact me via email. If you find this project helpful, please consider giving it a star.

## LICENSE

[MIT License](./LICENSE)
