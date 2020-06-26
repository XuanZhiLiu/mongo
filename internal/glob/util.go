package glob

import (
	"reflect"
)

func ToSlice(slice interface{}) (iSlice []interface{}, ok bool) {
	s := reflect.ValueOf(slice)
	ok = s.Kind() == reflect.Slice
	if !ok {
		return
	}

	iSlice = make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		iSlice[i] = s.Index(i).Interface()
	}
	return
}
