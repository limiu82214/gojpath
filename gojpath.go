// Package gojpath is a simple JSON path selector
// Package gojpath 是一個簡單的 JSON path 選擇器
package gojpath

import (
	"strconv"
	"strings"
)

func selectJSONNode(node interface{}, parts []string) interface{} {
	if len(parts) == 0 {
		return node
	}

	target := parts[0]

	switch node := node.(type) {
	case map[string]interface{}:
		if len(parts[1:]) != 0 {
			return selectJSONNode(node[target], parts[1:])
		}

		return node[target]
	case []interface{}:
		index, err := strconv.Atoi(target)
		if err != nil {
			return nil
		}

		if len(parts[1:]) == 0 {
			return node[index]
		}

		return selectJSONNode(node[index], parts[1:])
	default:
		return nil
	}
}

// Get return value which locate by JSON path from decoded JSON data
//
// Get 會從解碼後的 JSON 資料中，取得 JSON path 所指定的值
func Get(jsonData interface{}, path string) (interface{}, error) {
	if !strings.HasPrefix(path, "$") {
		panic("first char must be $")
	}

	path = path[1:]
	path = strings.ReplaceAll(path, "\"", "")
	path = strings.ReplaceAll(path, "'", "")
	path = strings.ReplaceAll(path, "[", ".")
	path = strings.ReplaceAll(path, "]", "")
	parts := strings.Split(path, ".")
	parts = parts[1:] // remove first empty string

	return selectJSONNode(jsonData, parts), nil
}
