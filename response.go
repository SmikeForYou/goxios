package goxios

import (
	"io"
	"net/http"
)

type Response struct {
	http.Response
}

func (r *Response) ReadBody() ([]byte, error) {
	return io.ReadAll(r.Body)
}
