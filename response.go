package goxios

import (
	"encoding/json"
	"io"
	"net/http"
)

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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	target := New[T]()
	err = json.Unmarshal(body, &target)
	if err != nil {
		return nil, err
	}
	return &JsonResponse[T]{
		Response: r,
		Value:    &target,
	}, nil
}
