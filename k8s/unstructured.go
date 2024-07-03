package k8s

import (
	"fmt"
	"strconv"
)

func M(data interface{}, keys ...string) interface{} {
	result := data
	for _, key := range keys {
		switch v := result.(type) {
		case map[string]interface{}:
			result = v[key]
		case map[interface{}]interface{}:
			result = v[key]
		}
		if result == nil {
			return result
		}
	}
	return result
}
func Ms(data interface{}, keys ...string) string {
	return fmt.Sprintf("%s", M(data, keys...))
}

func Mb(data interface{}, keys ...string) bool {
	return M(data, keys...).(bool)
}

func Mns(data interface{}, keys ...string) string {
	return strconv.Itoa(int(M(data, keys...).(float64)))

}

func Mn(data interface{}, keys ...string) int {
	return int(M(data, keys...).(float64))

}

func L(data interface{}) []interface{} {
	if data == nil {
		return make([]interface{}, 0)
	}
	return data.([]interface{})
}
