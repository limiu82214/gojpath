package gojpath

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

	result := Get(jsonData, "$.store.book[0].title")
	assert.Equal(t, "Harry Potter and the Philosopher's Stone", result)

	result = Get(jsonData, "$['store']['book'][0]['title']")
	assert.Equal(t, "Harry Potter and the Philosopher's Stone", result)

	result = Get(jsonData, "$.store.bicycle.color")
	assert.Equal(t, "red", result)

	result = Get(jsonData, "$.store.bicycle.price")
	assert.Equal(t, 19.95, result)

	result = Get(jsonData, "$.store.book[1].author")
	assert.Equal(t, "J.K. Rowling", result)

	// not support `*` because "https://learn.microsoft.com/en-us/azure/data-explorer/kusto/query/jsonpath"
	// not contain operation `*` and I use this tool for map golang struct
	// result = Get(jsonData, "$.store.bicycle.*")
	// assert.Equal(t, map[string]interface{}{"color": "red", "price": 19.95}, result)

	result = Get(jsonData, "$.store['book'][1].price")
	assert.Equal(t, 9.99, result)
}
