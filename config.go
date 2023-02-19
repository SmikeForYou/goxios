package goxios

import "net/http"

type Config struct {
	BaseURL string
	Headers http.Header
	Client  HttpClient
}

type RequestConfig struct {
	Headers     http.Header
	QueryParams map[string]any
}

func (receiver RequestConfig) queryParamsPrepared() string {
	return QueryParamsToStr(receiver.QueryParams)
}
