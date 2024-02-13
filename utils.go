package goxios

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

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

func New[T any]() T {
	var t T
	type_ := reflect.TypeOf(t)
	if type_.Kind() == reflect.Ptr {
		val := reflect.New(type_.Elem()).Interface()
		return val.(T)
	}
	return t
}
