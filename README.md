# gojpath

![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/limiu82214/gojpath?label=version) [![Go Reference](https://pkg.go.dev/badge/github.com/limiu82214/gojpath.svg)](https://pkg.go.dev/github.com/limiu82214/gojpath) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![codecov](https://codecov.io/gh/limiu82214/gojpath/branch/master/graph/badge.svg?token=0XAK9BB5WL)](https://codecov.io/gh/limiu82214/gojpath) [![Go Report Card](https://goreportcard.com/badge/github.com/limiu82214/gojpath)](https://goreportcard.com/report/github.com/limiu82214/gojpath) ![github actions workflow](https://github.com/limiu82214/gojpath/actions/workflows/go.yml/badge.svg)

[中文版文檔](./README_ZH.md)

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

### Zero Value, Compare Nil & Exist

You can use this package to check for zero values.

Sometimes, we map a request JSON to a struct using the Bind method of the Gin framework, like this:

```go
    req := RegisterUserReq{}
    err := c.BindJSON(&req)
    if err != nil {
        return err
    }
```

Usually, we verify the data using the built-in validator in Gin. However, the validator may not be able to detect the following situation:

* `zero value` or just client did not input it.
* `zero value` or client explicitly inputted `zero value`.

We want to determine whether a value is a zero value or simply not present in the request. Here is the solution for this:

```go
jsonString := c.Request.body
var jsonData interface{}

err := json.Unmarshal([]byte(jsonString), &jsonData)
if err != nil {
    t.Fatal(err)
}

// client not input it
isExist, err := IsExist("$.store['book'][1].price")
if !isExist {
    // ...
}

// client input zero value
isNil, err := IsNil("$.store['book'][1].price")
if isNil {
    // ...
}

// IsNilOrUnset return true if value which locate by JSON path is nil or not exist
// It mean the value of struct will be fill with zero value with json package.
isBindNil, err := IsNilOrUnset("$.store['book'][1].price")
if isBindNil {
    // ...
}
```
Notice: We use this because we dislike using pointers in the struct, as they make our code more complex.

## Function Behavior Example
```
// jsonString: {"Field2": true}
// json unmarshal Value: true
// IsNil: false
// IsExist: true
// IsNilOrUnset: false

// jsonString: {"Field2": false}
// json unmarshal Value: false
// IsNil: false
// IsExist: true
// IsNilOrUnset: false

// jsonString: {"Field2": null}
// json unmarshal Value: false
// IsNil: true
// IsExist: true
// IsNilOrUnset: true

// jsonString: {}
// json unmarshal Value: false
// IsNil error: object key not found
// IsExist: false
// IsNilOrUnset: true
 
// this is not a valid json, so parse will return error
// jsonString: {"Field2": undefined}
```

## Considerations

As the extent of support for the Get function is limited and some operations, such as the * operator, are not supported, it is recommended to check if the operator and syntax used in JSON Path queries are supported by the Get function to avoid generating incorrect results.

## Other

If you encounter any issues during use, please feel free to raise an issue on the GitHub project or contact me via email. If you find this project helpful, please consider giving it a star.

## LICENSE

[MIT License](./LICENSE)
