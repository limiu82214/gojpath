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

	result, err := Get(jsonData, "$")
	assert.Nil(t, err)
	assert.Equal(t, jsonData, result)

	result, err = Get(jsonData, "$.store.book[0].title")
	assert.Nil(t, err)
	assert.Equal(t, "Harry Potter and the Philosopher's Stone", result)

	result, err = Get(jsonData, "$['store']['book'][0]['title']")
	assert.Nil(t, err)
	assert.Equal(t, "Harry Potter and the Philosopher's Stone", result)

	result, err = Get(jsonData, "$.store.bicycle.color")
	assert.Nil(t, err)
	assert.Equal(t, "red", result)

	result, err = Get(jsonData, "$.store.bicycle.price")
	assert.Nil(t, err)
	assert.Equal(t, 19.95, result)

	result, err = Get(jsonData, "$.store.book[1].author")
	assert.Nil(t, err)
	assert.Equal(t, "J.K. Rowling", result)

	// not support `*` because "https://learn.microsoft.com/en-us/azure/data-explorer/kusto/query/jsonpath"
	// not contain operation `*` and I use this tool for map golang struct
	// result = Get(jsonData, "$.store.bicycle.*")
	// assert.Equal(t, map[string]interface{}{"color": "red", "price": 19.95}, result)

	result, err = Get(jsonData, "$.store['book'][1].price")
	assert.Nil(t, err)
	assert.Equal(t, 9.99, result)

	result, err = Get(jsonData, "$.store.book[1]")
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"title":  "Harry Potter and the Chamber of Secrets",
		"author": "J.K. Rowling",
		"price":  9.99,
	}, result)

	_, err = Get(jsonData, "$.store.book[fail].title")
	assert.ErrorIs(t, err, ErrArrayIndexNotNumber)

	_, err = Get(jsonData, "$.store.book[8].title")
	assert.ErrorIs(t, err, ErrArrayIndexOutOfRange)

	_, err = Get(jsonData, "$.notexist")
	assert.ErrorIs(t, err, ErrObjectKeyNotFound)

	_, err = Get(jsonData, "$.notexist")
	assert.ErrorIs(t, err, ErrObjectKeyNotFound)

	_, err = Get(jsonData, ".store")
	assert.ErrorIs(t, err, ErrFirstCharMustBeDollar)
}

func TestNodeErr(t *testing.T) {
	jsonString := `2`

	var jsonData interface{}

	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Get(jsonData, "$.")
	assert.ErrorIs(t, err, ErrNodeIsNotObjectOrArray)
}
