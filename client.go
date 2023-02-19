package goxios

import (
	"bytes"
	"net/http"
	"time"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
	SetTimeout(duration time.Duration)
}

type httpClient struct {
	http.Client
}

func (client *httpClient) SetTimeout(duration time.Duration) {
	client.Timeout = duration
}

var DefaultHeaders = http.Header{}
var DefaultGoxiosInstance = NewGoxiosInstance(Config{
	BaseURL: "",
	Headers: DefaultHeaders,
	Client:  &httpClient{},
})

type GoxiosInstance struct {
	baseUrl string
	headers http.Header
	client  HttpClient
}

func newGoxiosInstance(config Config) *GoxiosInstance {
	if config.Headers == nil {
		config.Headers = http.Header{}
	}
	if config.Client == nil {
		config.Client = &httpClient{}
	}
	return &GoxiosInstance{
		baseUrl: config.BaseURL,
		headers: config.Headers,
		client:  config.Client,
	}
}

func NewGoxiosInstance(config Config) *GoxiosInstance {
	return newGoxiosInstance(config)
}

func (g *GoxiosInstance) urlFor(u string) (string, error) {
	return urlfor(g.baseUrl, u)
}

func (g *GoxiosInstance) SetBaseUrl(baseUrl string) {
	g.baseUrl = baseUrl
}

func (g *GoxiosInstance) SetClient(client HttpClient) {
	g.client = client
}

func (g *GoxiosInstance) SetHeaders(headers http.Header) {
	g.headers = headers
}

func (g *GoxiosInstance) SetHeader(k, v string) {
	g.headers.Set(k, v)
}

func (g *GoxiosInstance) GetHeader(k string) string {
	return g.headers.Get(k)
}

func (g *GoxiosInstance) GetHeaders() http.Header {
	return g.headers
}

func (g *GoxiosInstance) AddHeader(k, v string) {
	g.headers.Add(k, v)
}

func (g *GoxiosInstance) SetRequestTimeout(duration time.Duration) {
	g.client.SetTimeout(duration)
}

func (g *GoxiosInstance) mergeHeaders(headers http.Header) http.Header {
	h := make(http.Header)
	if g.headers != nil {
		h = g.headers.Clone()
	}
	if headers != nil {
		for k, v := range headers {
			for _, s := range v {
				h.Add(k, s)
			}
		}
	}
	return h
}

func (g *GoxiosInstance) Request(method string, url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	u, err := g.urlFor(url)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &RequestConfig{}
	}
	config.Headers = g.mergeHeaders(config.Headers)
	return request(method, u, payload, *config, g.client)
}

func (g *GoxiosInstance) Get(url string, config *RequestConfig) (*Response, error) {
	return g.Request("GET", url, nil, config)
}

func (g *GoxiosInstance) Post(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return g.Request("POST", url, payload, config)
}

func (g *GoxiosInstance) Put(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return g.Request("PUT", url, payload, config)
}

func (g *GoxiosInstance) Patch(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return g.Request("PATCH", url, payload, config)
}

func (g *GoxiosInstance) Delete(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return g.Request("DELETE", url, payload, config)
}

func request(method string, url string, payload RequestPayload, config RequestConfig, client HttpClient) (*Response, error) {
	var buff = bytes.NewBuffer([]byte{})
	var contentType = "text/plain"
	var err error
	if payload != nil {
		buff, err = payload.toBuff()
		if err != nil {
			return nil, err
		}
		contentType = payload.contentType()
	}
	if err != nil {
		return nil, err
	}
	if config.QueryParams != nil {
		url = url + "?" + QueryParamsToStr(config.QueryParams)
	}
	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return nil, err
	}
	if config.Headers != nil {
		for k, v := range config.Headers {
			for _, s := range v {
				req.Header.Set(k, s)
			}
		}
	}
	req.Header.Set("Content-Type", contentType)
	if client == nil {
		client = &httpClient{}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	response := Response{Response: *resp}
	return &response, nil
}

func Get(url string, config *RequestConfig) (*Response, error) {
	return DefaultGoxiosInstance.Get(url, config)
}

func Post(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return DefaultGoxiosInstance.Post(url, payload, config)
}

func Put(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return DefaultGoxiosInstance.Put(url, payload, config)
}

func Patch(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return DefaultGoxiosInstance.Patch(url, payload, config)
}

func Delete(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return DefaultGoxiosInstance.Delete(url, payload, config)
}
