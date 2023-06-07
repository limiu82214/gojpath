// Package gojpath is a simple JSON path selector
//
// Package gojpath 是一個簡單的 JSON path 選擇器
package gojpath

import (
	"errors"
	"strconv"
	"strings"
)

// ErrArrayIndexOutOfRange is error when array index out of range
//
// ErrArrayIndexOutOfRange 是當陣列索引超出範圍時的錯誤
var ErrArrayIndexOutOfRange = errors.New("array index out of range")

// ErrArrayIndexNotNumber is error when array index is not number
//
// ErrArrayIndexNotNumber 是當陣列索引不是數字時的錯誤
var ErrArrayIndexNotNumber = errors.New("array index not number")

// ErrObjectKeyNotFound is error when object key not found
//
// ErrObjectKeyNotFound 是當物件鍵找不到時的錯誤
var ErrObjectKeyNotFound = errors.New("object key not found")

// ErrFirstCharMustBeDollar is error when first char of JSON path is not $
//
// ErrFirstCharMustBeDollar 是當 JSON path 的第一個字元不是 $ 時的錯誤
var ErrFirstCharMustBeDollar = errors.New("first char must be $")

// ErrNodeIsNotObjectOrArray is error when node is not object or array
//
// ErrNodeIsNotObjectOrArray 是當節點不是物件或陣列時的錯誤
var ErrNodeIsNotObjectOrArray = errors.New("node is not object or array")

func selectJSONNode(node interface{}, parts []string) (interface{}, error) {
	if len(parts) == 0 {
		return node, nil
	}

	target := parts[0]

	switch node := node.(type) {
	case map[string]interface{}:
		if _, ok := node[target]; !ok {
			return nil, ErrObjectKeyNotFound
		}

		if len(parts[1:]) != 0 {
			return selectJSONNode(node[target], parts[1:])
		}

		return node[target], nil
	case []interface{}:
		index, err := strconv.Atoi(target)
		if err != nil {
			return nil, ErrArrayIndexNotNumber
		}

		if index >= len(node) || index < 0 {
			return nil, ErrArrayIndexOutOfRange
		}

		if len(parts[1:]) == 0 {
			return node[index], nil
		}

		return selectJSONNode(node[index], parts[1:])
	default:
		return nil, ErrNodeIsNotObjectOrArray
	}
}

// Get return value which locate by JSON path from decoded JSON data
//
// Get 會從解碼後的 JSON 資料中，取得 JSON path 所指定的值
func Get(jsonData interface{}, path string) (interface{}, error) {
	if !strings.HasPrefix(path, "$") {
		return nil, ErrFirstCharMustBeDollar
	}

	path = path[1:]
	path = strings.ReplaceAll(path, "\"", "")
	path = strings.ReplaceAll(path, "'", "")
	path = strings.ReplaceAll(path, "[", ".")
	path = strings.ReplaceAll(path, "]", "")
	parts := strings.Split(path, ".")
	parts = parts[1:] // remove first empty string

	return selectJSONNode(jsonData, parts)
}

// IsNil return true if value which locate by JSON path is nil
//
// IsNil 會回傳 JSON path 所指定的值是否為 nil
func IsNil(jsonData interface{}, path string) (bool, error) {
	v, err := Get(jsonData, path)
	if err == nil && v == nil {
		return true, nil
	} else if err != nil {
		return false, err
	}

	return false, nil
}

// IsExist return true if value which locate by JSON path is exist
//
// IsExist 會回傳 JSON path 所指定的值是否存在
func IsExist(jsonData interface{}, path string) (bool, error) {
	_, err := Get(jsonData, path)
	if errors.Is(err, ErrObjectKeyNotFound) || errors.Is(err, ErrArrayIndexOutOfRange) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// IsBindNil return true if value which locate by JSON path is nil or not exist.
// It mean the value of struct will be fill with zero value with json package.
//
// IsBindNil 會回傳 JSON path 所指定的值是否為 nil 或不存在
// 這意味著該 struct 的值將會被填入零值
func IsBindNil(jsonData interface{}, path string) (bool, error) {
	isExist, err := IsExist(jsonData, path)
	if err != nil {
		return false, err
	}

	if !isExist {
		return true, nil
	}

	isNil, err := IsNil(jsonData, path)
	if err != nil {
		return false, err
	}

	if isNil {
		return true, nil
	}

	return false, nil
}
