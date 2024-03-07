package goxios

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

type mockClient struct{}

func (m mockClient) Do(req *http.Request) (*http.Response, error) {
	switch req.Method {
	case "GET":
		switch req.URL.String() {
		case "http://e.com/json":
			data := []byte(`{"a": 1, "b": 2, "c": 3}`)
			pipeR, pipeW := io.Pipe()
			go func() {
				//nolint:all
				pipeW.Write(data)
				//nolint:all
				pipeW.Close()
			}()
			resp := &http.Response{
				StatusCode: 200,
				Body:       pipeR,
				Header:     make(http.Header),
			}
			return resp, nil
		case "http://e.com/stream":
			return &http.Response{}, nil
		}
	case "POST", "PUT", "PATCH", "DELETE":
		switch req.URL.String() {
		case "http://e.com/json":
			reqBody, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			statusCode, err := strconv.Atoi(getHeaderOrDefault(req.Header, "X-Return-Status", "200"))
			if err != nil {
				return nil, err
			}
			returnBody := getHeaderOrDefault(req.Header, "X-Return-Payload", string(reqBody))
			resp := &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(returnBody)),
			}
			return resp, nil
		case "http://e.com/form-data":
			return nil, nil
		}
	}
	return &http.Response{
		StatusCode: 200,
		Body:       nil,
		Header:     make(http.Header),
	}, nil
}

func (m mockClient) SetTimeout(_ time.Duration) {}

func TestRequest(t *testing.T) {
	url := "https://httpbin.org/get"
	r, err := request("GET", url, nil, RequestConfig{
		Headers:     nil,
		QueryParams: nil,
	}, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.StatusCode)

}

func TestClientGetJsonReal(t *testing.T) {
	url := "https://httpbin.org"
	c := NewGoxiosInstance(Config{BaseURL: url})
	resp, err := c.Get("/get", &RequestConfig{QueryParams: map[string]any{"a": "1", "b": "2"}})
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	b, err := resp.ReadBody()
	assert.Nil(t, err)
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	assert.Nil(t, err)
	assert.Equal(t, data["args"].(map[string]interface{}), map[string]interface{}{"a": "1", "b": "2"})
}

func TestClientPostFormDataReal(t *testing.T) {
	url := "https://httpbin.org"
	c := NewGoxiosInstance(Config{BaseURL: url})
	fr := NewFormDataRequest()
	fr.AddValueStr("a", "1")
	fr.AddValueStr("a", "2")
	fr.AddValueStr("b", "2")
	resp, err := c.Post("/post", fr, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	b, err := resp.ReadBody()
	assert.Nil(t, err)
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	assert.Nil(t, err)
	assert.Equal(t, data["form"].(map[string]interface{}), map[string]interface{}{"a": []any{"1", "2"}, "b": "2"})
}

func TestClientGetJson(t *testing.T) {
	url := "http://e.com"
	c := newGoxiosInstance(Config{url, nil, &mockClient{}})
	resp, err := c.Get("/json", nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	b, err := resp.ReadBody()
	assert.Nil(t, err)
	assert.Equal(t, `{"a": 1, "b": 2, "c": 3}`, string(b))
}

func TestClientGetJsonStructuredResponse(t *testing.T) {
	url := "http://e.com"
	c := newGoxiosInstance(Config{url, nil, &mockClient{}})
	resp, err := c.Get("/json", nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	type RespStruct struct {
		A int `json:"a"`
		B int `json:"b"`
		C int `json:"c"`
	}
	json_resp, err := Json[RespStruct](resp)
	assert.Nil(t, err)
	assert.Equal(t, *json_resp.Body(), RespStruct{1, 2, 3})
}
