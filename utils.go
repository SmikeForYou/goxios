package goxios

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"unsafe"
)

// s2b converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func s2b(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func QueryParamsToStr(params map[string]any) string {
	var pp = make([]string, 0)
	for key, value := range params {
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			for _, v := range value.([]string) {
				pp = append(pp, fmt.Sprintf("%s=%v", key, v))
			}
			continue
		}
		pp = append(pp, fmt.Sprintf("%s=%v", key, value))
	}
	return strings.Join(pp, "&")
}

func urlfor(base string, u string) (string, error) {
	if base != "" {
		parsed, err := url.Parse(base)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"%s/%s",
			strings.TrimRight(parsed.String(), "/"),
			strings.TrimLeft(u, "/"),
		), nil
	}
	return u, nil
}
