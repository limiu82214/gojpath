# gojpath
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/limiu82214/gojpath?label=version) [![Go Reference](https://pkg.go.dev/badge/github.com/limiu82214/gojpath.svg)](https://pkg.go.dev/github.com/limiu82214/gojpath) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![codecov](https://codecov.io/gh/limiu82214/gojpath/branch/master/graph/badge.svg?token=0XAK9BB5WL)](https://codecov.io/gh/limiu82214/gojpath) [![Go Report Card](https://goreportcard.com/badge/github.com/limiu82214/gojpath)](https://goreportcard.com/report/github.com/limiu82214/gojpath) ![github actions workflow](https://github.com/limiu82214/gojpath/actions/workflows/go.yml/badge.svg)


`gojpath` 是一種用於在 JSON 數據中進行查詢的語言，它類似於 XPath。  
在 Golang 中，可以使用 `Get` 函數對 JSON 數據進行查詢。  
該函數支援的程度參考自 [link](https://learn.microsoft.com/en-us/azure/data-explorer/kusto/query/jsonpath)  

## 使用方法

以下是使用 `Get` 函數進行 JSON Path 查詢的示例代碼:

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

    // 使用 JSON Path 進行查詢
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

在這個示例中，TestGet 函數使用 Get 函數進行了多次 JSON Path 查詢。  
透過比較返回的值和預期的值，可以確認查詢的結果是否正確。

### 零值，比較 Nil & Exist
您可以使用此套件來檢查零值。

有時候，我們使用 Gin 框架的 Bind 方法將請求的 JSON 映射到結構體，像這樣：

```go
    req := RegisterUserReq{}
    err := c.BindJSON(&req)
    if err != nil {
        return err
    }
```

通常，我們使用 Gin 內建的驗證器來驗證資料。然而，驗證器可能無法檢測到以下的情況：

* 零值 或者客戶端沒有輸入。
* 零值 或者客戶端明確地輸入 零值。

我們希望判斷一個值是零值還是在請求中根本不存在。以下是解決這個問題的方法：

```go
jsonString := c.Request.body
var jsonData interface{}

err := json.Unmarshal([]byte(jsonString), &jsonData)
if err != nil {
    t.Fatal(err)
}

// 客戶端未輸入
isExist, err := IsExist("$.store['book'][1].price")
if !isExist {
    // ...
}

// 客戶端輸入零值
isNil, err := IsNil("$.store['book'][1].price")
if isNil {
    // ...
}

// 如果位於 JSON 路徑的值是 nil 或不存在，則 IsNilOrUnset 返回 true
// 這意味著結構體的值會被 json 套件填充為零值。
isBindNil, err := IsNilOrUnset("$.store['book'][1].price")
if isBindNil {
    // ...
}
```
注意：我們使用這個是因為我們不喜歡在結構體中使用指標，因為他們使我們的程式碼變得更複雜。

## 函式行為範例 
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

## 注意事項
由於 Get 函數的支援程度有限，不支援某些操作，例如 * 操作符。在進行 JSON Path 查詢時，應該先確認使用的操作符和語法是否被 Get 函數所支援，以免產生錯誤的結果。

## 其他
如果您在使用過程中有任何問題，歡迎在 GitHub 專案上發起一個 issue，或是透過 email 與我聯繫。如果您認為這個專案對您有所幫助，也請不吝給予一個 star。


## 授權

[MIT License](./LICENSE)
