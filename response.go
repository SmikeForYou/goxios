package goxios

import (
	"encoding/json"
	"io"
	"net/http"
)

var DefaultUnmarshaller = json.Unmarshal

type Response struct {
	http.Response
}

func (r *Response) ReadBody() ([]byte, error) {
	return io.ReadAll(r.Body)
}

type JsonResponse[T any] struct {
	*Response
	Value *T
}

func (r *JsonResponse[T]) Body() *T {
	return r.Value
}

func Json[T any](r *Response) (*JsonResponse[T], error) {
	return NewJsonParser[T](DefaultUnmarshaller).Parse(r)
}

type JsonParser[T any] struct {
	unmarshaller func(data []byte, v interface{}) error
}

// NewJsonParser creates a new JSON parser
func NewJsonParser[T any](unmarshaler func([]byte, any) error) *JsonParser[T] {
	return &JsonParser[T]{
		unmarshaller: unmarshaler,
	}
}

// Parse parses the response body into a JSON object
func (p *JsonParser[T]) Parse(r *Response) (*JsonResponse[T], error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	target := New[T]()
	err = p.unmarshaller(body, &target)
	if err != nil {
		return nil, err
	}
	return &JsonResponse[T]{
		Response: r,
		Value:    &target,
	}, nil
}
