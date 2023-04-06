[中文版](./README_ZH.md)
# gojpath

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
    result := Get(jsonData, "$.store.book[0].title")
    assert.Equal(t, "Harry Potter and the Philosopher's Stone", result)

    result = Get(jsonData, "$['store']['book'][0]['title']")
    assert.Equal(t, "Harry Potter and the Philosopher's Stone", result)

    result = Get(jsonData, "$.store.bicycle.color")
    assert.Equal(t, "red", result)

    result = Get(jsonData, "$.store.book[1].author")
    assert.Equal(t, "J.K. Rowling", result)

    result = Get(jsonData, "$.store['book'][1].price")
    assert.Equal(t, 9.99, result)
}
```

在這個示例中，TestGet 函數使用 Get 函數進行了多次 JSON Path 查詢。  
透過比較返回的值和預期的值，可以確認查詢的結果是否正確。

## 注意事項
由於 Get 函數的支援程度有限，不支援某些操作，例如 * 操作符。在進行 JSON Path 查詢時，應該先確認使用的操作符和語法是否被 Get 函數所支援，以免產生錯誤的結果。

## 其他
如果您在使用過程中有任何問題，歡迎在 GitHub 專案上發起一個 issue，或是透過 email 與我聯繫。如果您認為這個專案對您有所幫助，也請不吝給予一個 star。
