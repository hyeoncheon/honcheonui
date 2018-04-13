package utils

import (
	"errors"
	"reflect"
	"strings"
)

// ToInterface make []interface{} from []Type
// See: https://golang.org/doc/faq#convert_slice_of_interface
func ToInterface(src interface{}) ([]interface{}, error) {
	srcValue := reflect.ValueOf(src)
	srcType := srcValue.Type()
	if srcType.Kind() != reflect.Array && srcType.Kind() != reflect.Slice {
		return nil, errors.New("not array or slice")
	}
	s := make([]interface{}, srcValue.Len())
	for i := 0; i < srcValue.Len(); i++ {
		s[i] = srcValue.Index(i).Interface()
	}
	return s, nil
}

// Has returns true if the item exists on the array. otherwise returns false.
func Has(array []interface{}, item interface{}) bool {
	for _, e := range array {
		if e == item {
			return true
		}
	}
	return false
}

// Remove search the item on the array and remove it.
// NOTE: given array also modified. see https://blog.golang.org/go-slices-usage-and-internals
func Remove(array []interface{}, item interface{}) []interface{} {
	for i, e := range array {
		if e == item {
			ret := append(array[:i], array[i+1:]...)
			return ret
		}
	}
	return array
}

// Cleaner has bad name but works for me
func Cleaner(a []string) []string {
	var na []string
	for _, e := range a {
		na = append(na, strings.TrimSpace(e))
	}
	return na
}
