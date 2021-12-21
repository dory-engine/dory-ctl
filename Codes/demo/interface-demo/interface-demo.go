package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func removeItem(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		vv := reflect.ValueOf(v)
		switch v.(type) {
		case int:
			if vv.IsZero() {
				delete(m, k)
			}
		case bool:
			if vv.Bool() == false {
				delete(m, k)
			}
		case string:
			if vv.String() == "" {
				delete(m, k)
			}
		}
		if !vv.IsValid() {
			delete(m, k)
		}
		if vv.Kind() == reflect.Slice {
			if vv.Len() == 0 {
				delete(m, k)
			} else {
				x := []map[string]interface{}{}
				for i := 0; i < vv.Len(); i++ {
					vvv := reflect.ValueOf(vv.Index(i))
					if vvv.Kind() == reflect.Struct {
						v3 := removeItem(vv.Index(i).Interface().(map[string]interface{}))
						x = append(x, v3)
					}
				}
				m[k] = x
			}
		}
		if vv.Kind() == reflect.Map {
			v2 := removeItem(v.(map[string]interface{}))
			m[k] = v2
		}
	}
	m2 := m
	return m2
}

func main() {
	m := map[string]interface{}{
		"string_foo":   "foo",
		"string_empty": "",
		"int_zero":     0,
		"int":          1,
		"bool_false":   false,
		"bool_true":    true,
		"nil":          nil,
		"array_empty":  []string{},
		"array_maps": []map[string]interface{}{
			{
				"string_foo":   "foo",
				"string_empty": "",
				"int_zero":     0,
				"int":          1,
				"bool_false":   false,
				"bool_true":    true,
				"nil":          nil,
				"array":        []string{},
			},
			{
				"string_foo":   "foo",
				"string_empty": "",
				"int_zero":     0,
				"int":          1,
				"bool_false":   false,
				"bool_true":    true,
				"nil":          nil,
				"array":        []string{},
			},
		},
		"map": map[string]interface{}{
			"string_foo":   "foo",
			"string_empty": "",
			"int_zero":     0,
			"int":          1,
			"bool_false":   false,
			"bool_true":    true,
			"nil":          nil,
			"array_empty":  []string{},
			"array_maps": []map[string]interface{}{
				{
					"string_foo":   "foo",
					"string_empty": "",
					"int_zero":     0,
					"int":          1,
					"bool_false":   false,
					"bool_true":    true,
					"nil":          nil,
					"array":        []string{},
				},
				{
					"string_foo":   "foo",
					"string_empty": "",
					"int_zero":     0,
					"int":          1,
					"bool_false":   false,
					"bool_true":    true,
					"nil":          nil,
					"array":        []string{},
				},
			},
		},
	}
	m2 := removeItem(m)
	data, _ := json.Marshal(m2)
	fmt.Println(string(data))
}
