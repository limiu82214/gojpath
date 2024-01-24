package gojpath

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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

func TestIsNil(t *testing.T) {
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
            },
			"null_data": null
        }
    }`

	var jsonData interface{}

	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		t.Fatal(err)
	}

	// function test
	isNil, err := IsNil(jsonData, "$.store.null_data")
	assert.Nil(t, err)
	assert.True(t, isNil)

	isNil, err = IsNil(jsonData, "$.store.book[0].title")
	assert.Nil(t, err)
	assert.False(t, isNil)

	// other test
	isNil, err = IsNil(jsonData, "$.notexist")
	assert.ErrorIs(t, err, ErrObjectKeyNotFound)
	assert.False(t, isNil)

	isNil, err = IsNil(jsonData, "$.store")
	assert.Nil(t, err)
	assert.False(t, isNil)

	isNil, err = IsNil(jsonData, "$.store.book[0].deep.notexist")
	assert.ErrorIs(t, err, ErrObjectKeyNotFound)
	assert.False(t, isNil)

	isNil, err = IsNil(jsonData, "$.store.book[fail].title")
	assert.ErrorIs(t, err, ErrArrayIndexNotNumber)
	assert.False(t, isNil)

	isNil, err = IsNil(jsonData, ".store")
	assert.ErrorIs(t, err, ErrFirstCharMustBeDollar)
	assert.False(t, isNil)
}

func TestIsExist(t *testing.T) {
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
            },
			"null_data": null
        }
    }`

	var jsonData interface{}

	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		t.Fatal(err)
	}

	// function test
	isExist, err := IsExist(jsonData, "$.notexist")
	assert.Nil(t, err)
	assert.False(t, isExist)

	isExist, err = IsExist(jsonData, "$.store.book[0].title")
	assert.Nil(t, err)
	assert.True(t, isExist)

	// other test
	isExist, err = IsExist(jsonData, "$.store.null_data")
	assert.Nil(t, err)
	assert.True(t, isExist)

	isExist, err = IsExist(jsonData, "$.store")
	assert.Nil(t, err)
	assert.True(t, isExist)

	isExist, err = IsExist(jsonData, "$.store.book[0].deep.notexist")
	assert.Nil(t, err)
	assert.False(t, isExist)

	isExist, err = IsExist(jsonData, "$.store.book[fail].title")
	assert.ErrorIs(t, err, ErrArrayIndexNotNumber)
	assert.False(t, isExist)

	isExist, err = IsExist(jsonData, ".store")
	assert.ErrorIs(t, err, ErrFirstCharMustBeDollar)
	assert.False(t, isExist)
}

func TestIsBindNil(t *testing.T) {
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
            },
			"null_data": null
        }
    }`

	var jsonData interface{}

	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		t.Fatal(err)
	}

	// function test
	isBindNil, err := IsNilOrUnset(jsonData, "$.notexist")
	assert.Nil(t, err)
	assert.True(t, isBindNil)

	isExist, _ := IsExist(jsonData, "$.notexist")
	isNil, _ := IsNil(jsonData, "$.notexist")
	assert.Equal(t, !isExist || isNil, isBindNil)

	isBindNil, err = IsNilOrUnset(jsonData, "$.store.book[0].title")
	assert.Nil(t, err)
	assert.False(t, isBindNil)

	isExist, _ = IsExist(jsonData, "$.store.book[0].title")
	isNil, _ = IsNil(jsonData, "$.store.book[0].title")
	assert.Equal(t, !isExist || isNil, isBindNil)

	// other test
	isBindNil, err = IsNilOrUnset(jsonData, "$.store.null_data")
	assert.Nil(t, err)
	assert.True(t, isBindNil)

	isExist, _ = IsExist(jsonData, "$.store.null_data")
	isNil, _ = IsNil(jsonData, "$.store.null_data")
	assert.Equal(t, !isExist || isNil, isBindNil)

	isBindNil, err = IsNilOrUnset(jsonData, "$.store")
	assert.Nil(t, err)
	assert.False(t, isBindNil)

	isExist, _ = IsExist(jsonData, "$.store")
	isNil, _ = IsNil(jsonData, "$.store")
	assert.Equal(t, !isExist || isNil, isBindNil)

	isBindNil, err = IsNilOrUnset(jsonData, "$.store.book[0].deep.notexist")
	assert.Nil(t, err)
	assert.True(t, isBindNil)

	isExist, _ = IsExist(jsonData, "$.store.book[0].deep.notexist")
	isNil, _ = IsNil(jsonData, "$.store.book[0].deep.notexist")
	assert.Equal(t, !isExist || isNil, isBindNil)

	isBindNil, err = IsNilOrUnset(jsonData, "$.store.book[fail].title")
	assert.ErrorIs(t, err, ErrArrayIndexNotNumber)
	assert.False(t, isBindNil)

	isBindNil, err = IsNilOrUnset(jsonData, ".store")
	assert.ErrorIs(t, err, ErrFirstCharMustBeDollar)
	assert.False(t, isBindNil)
}

func TestBoolTrue(t *testing.T) {
	testIsNilRst(t, `{"Field2": true}`, `$.Field2`, false, nil)
	testIsExistRst(t, `{"Field2": true}`, `$.Field2`, true, nil)
	testIsBindNilRst(t, `{"Field2": true}`, `$.Field2`, false, nil)
}
func TestBoolFalse(t *testing.T) {
	testIsNilRst(t, `{"Field2": false}`, `$.Field2`, false, nil)
	testIsExistRst(t, `{"Field2": false}`, `$.Field2`, true, nil)
	testIsBindNilRst(t, `{"Field2": false}`, `$.Field2`, false, nil)
}

func TestBoolNil(t *testing.T) {
	testIsNilRst(t, `{"Field2": null}`, `$.Field2`, true, nil)
	testIsExistRst(t, `{"Field2": null}`, `$.Field2`, true, nil)
	testIsBindNilRst(t, `{"Field2": null}`, `$.Field2`, true, nil)
}
func TestBoolNotExist(t *testing.T) {
	testIsNilRst(t, `{}`, `$.Field2`, false, ErrObjectKeyNotFound)
	testIsExistRst(t, `{}`, `$.Field2`, false, nil)
	testIsBindNilRst(t, `{}`, `$.Field2`, true, nil)
}

func testIsNilRst(t *testing.T, jsonString string, path string, expectIsNil bool, expectErr error) {
	T := struct {
		Field2 bool
	}{}
	err := json.Unmarshal([]byte(jsonString), &T)
	assert.Nil(t, err)

	var jsonData interface{}
	err = json.Unmarshal([]byte(jsonString), &jsonData)
	assert.Nil(t, err)

	isNil, err := IsNil(jsonData, path)
	assert.Equal(t, expectIsNil, isNil)
	assert.ErrorIs(t, err, expectErr)
}

func testIsExistRst(t *testing.T, jsonString string, path string, expectIsExist bool, expectErr error) {
	T := struct {
		Field2 bool
	}{}
	err := json.Unmarshal([]byte(jsonString), &T)
	assert.Nil(t, err)

	var jsonData interface{}
	err = json.Unmarshal([]byte(jsonString), &jsonData)
	assert.Nil(t, err)

	isExist, err := IsExist(jsonData, path)
	assert.Equal(t, expectIsExist, isExist)
	assert.ErrorIs(t, err, expectErr)
}

func testIsBindNilRst(t *testing.T, jsonString string, path string, expectIsBindNil bool, expectErr error) {
	T := struct {
		Field2 bool
	}{}
	err := json.Unmarshal([]byte(jsonString), &T)
	assert.Nil(t, err)

	var jsonData interface{}
	err = json.Unmarshal([]byte(jsonString), &jsonData)
	assert.Nil(t, err)

	isBindNil, err := IsNilOrUnset(jsonData, path)
	assert.Equal(t, expectIsBindNil, isBindNil)
	assert.ErrorIs(t, err, expectErr)
}

func TestBehave(t *testing.T) {
	isPrintToStd := false
	showResult(`{"Field2": true}`, isPrintToStd)
	showResult(`{"Field2": false}`, isPrintToStd)
	showResult(`{"Field2": null}`, isPrintToStd)
	showResult(`{}`, isPrintToStd)
	//showResult(`{"Field2": undefined}`, isPrintToStd) // undefined is not a valid json
}

func showResult(jsonString string, printToStd bool) {
	var logger *log.Logger

	if printToStd {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	logger.Println("jsonString:", jsonString)

	T := struct {
		Field2 bool
	}{}
	json.Unmarshal([]byte(jsonString), &T)
	logger.Printf("json unmarshal Value: %v\n", T.Field2)

	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		logger.Printf("Unmarshal error: %v", err)
	}

	isNil, err := IsNil(jsonData, `$.Field2`)
	if err != nil {
		logger.Printf("IsNil error: %v\n", err)
	} else {
		logger.Printf("IsNil: %v\n", isNil)
	}

	isExist, err := IsExist(jsonData, `$.Field2`)
	if err != nil {
		logger.Printf("IsExist error: %v\n", err)
	} else {
		logger.Printf("IsExist: %v\n", isExist)
	}

	isBindNil, err := IsNilOrUnset(jsonData, `$.Field2`)
	if err != nil {
		logger.Printf("IsNilOrUnset error: %v\n", err)
	} else {
		logger.Printf("IsNilOrUnset: %v\n", isBindNil)
	}
	logger.Println()
}
