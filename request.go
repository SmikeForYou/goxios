package goxios

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

var DefaultMarshaller = json.Marshal

type RequestPayload interface {
	contentType() string
	toBuff() (*bytes.Buffer, error)
}

type FormDataRequest struct {
	values map[string][]io.Reader
	cType  string
}

func NewFormDataRequest() *FormDataRequest {
	return &FormDataRequest{values: make(map[string][]io.Reader)}
}

func (r *FormDataRequest) contentType() string {
	if r.cType == "" {
		r.cType = "multipart/form-data"
	}
	return r.cType
}

func (r *FormDataRequest) toBuff() (*bytes.Buffer, error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	defer w.Close()
	for key, fd := range r.values {
		for _, reader := range fd {
			var fw io.Writer
			if x, ok := reader.(io.Closer); ok {
				defer x.Close()
			}
			if x, ok := reader.(*os.File); ok {
				if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
					return nil, err
				}
			} else {
				if fw, err = w.CreateFormField(key); err != nil {
					return nil, err
				}
			}
			if _, err = io.Copy(fw, reader); err != nil {
				return nil, err
			}
		}

	}
	r.cType = w.FormDataContentType()
	return &b, nil
}

func (r *FormDataRequest) AddValueStr(key, val string) {
	if r.values[key] == nil {
		r.values[key] = make([]io.Reader, 0)
	}
	r.values[key] = append(r.values[key], strings.NewReader(val))
}

func (r *FormDataRequest) AddValue(key string, val []byte) {
	if r.values[key] == nil {
		r.values[key] = make([]io.Reader, 0)
	}
	r.values[key] = append(r.values[key], bytes.NewReader(val))
}

func (r *FormDataRequest) AddFile(key string, file *os.File) {
	if r.values[key] == nil {
		r.values[key] = make([]io.Reader, 0)
	}
	r.values[key] = append(r.values[key], file)
}

type JsonRequest[T any] struct {
	data       T
	marshaller func([]byte, T) error
}

// NewJsonRequest creates a JsonRequest
func NewJsonRequest[T any](data T) JsonRequest[T] {
	return JsonRequest[T]{
		data: data,
	}
}

// SetMarshaller sets the marshaller for the JsonRequest
func (j JsonRequest[T]) SetMarshaller(m func([]byte, T) error) {
	j.marshaller = m
}

func (j JsonRequest[T]) contentType() string {
	return "application/json"
}

func (j JsonRequest[T]) toBuff() (*bytes.Buffer, error) {
	raw, err := DefaultMarshaller(j.data)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(raw), nil
}
