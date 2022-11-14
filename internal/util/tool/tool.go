package tool

import (
	"crypto/md5"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strings"
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

func ClientIp(r *http.Request) string {
	ip := r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if res := net.ParseIP(i); res != nil {
			return res.String()
		}
	}
	ip = r.Header.Get("X-Real-IP")
	if res := net.ParseIP(ip); res != nil {
		return res.String()
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	if net.ParseIP(ip) != nil {
		return ip
	}

	return ""
}
