package tool

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"time"
)

func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func MD55(input string) string {
	has := md5.Sum([]byte(input))
	md5Str := fmt.Sprintf("%x", has)

	return md5Str
}

func ArrSearch(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}
