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

// Request sends a request with a given method, url, payload and config
func request(method string, url string, payload RequestPayload, config RequestConfig, client HttpClient) (*Response, error) {
	builder := newRequestBuilder()
	builder.
		Method(method).
		Url(url).
		Payload(payload).
		Headers(config.Headers).
		QueryParams(config.QueryParams)

	return builder.Do(client)
}

// Get request sends GET request with a given url and config
func Get(url string, config *RequestConfig) (*Response, error) {
	return DefaultGoxiosInstance.Get(url, config)
}

// Post request sends POST request with a given url, payload and config
func Post(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return DefaultGoxiosInstance.Post(url, payload, config)
}

// Put request sends PUT request with a given url, payload and config
func Put(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return DefaultGoxiosInstance.Put(url, payload, config)
}

// Patch request sends PATCH request with a given url, payload and config
func Patch(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return DefaultGoxiosInstance.Patch(url, payload, config)
}

// Delete request sends DELETE request with a given url, payload and config
func Delete(url string, payload RequestPayload, config *RequestConfig) (*Response, error) {
	return DefaultGoxiosInstance.Delete(url, payload, config)
}

type RequestBuilder struct {
	method      string
	url         string
	queryParams map[string]any
	payload     *bytes.Buffer
	contentType string
	headers     http.Header
	error       error
}

func (b *RequestBuilder) Build() (*http.Request, error) {
	if b.error != nil {
		return nil, b.error
	}
	if b.payload == nil {
		b.payload = bytes.NewBuffer([]byte{})
	}
	req, err := http.NewRequest(b.method, b.BuildQueryParams(), b.payload)
	if err != nil {
		return nil, err
	}
	if b.headers != nil {
		for k, v := range b.headers {
			for _, s := range v {
				req.Header.Set(k, s)
			}
		}

	}
	if b.contentType != "" {
		req.Header.Set("Content-Type", b.contentType)
	}
	return req, nil
}

func (b *RequestBuilder) Do(client HttpClient) (*Response, error) {
	req, err := b.Build()
	if err != nil {
		return nil, err
	}
	if client == nil {
		client = DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	response := Response{Response: *resp}
	return &response, nil
}

func (b *RequestBuilder) BuildQueryParams() string {
	if len(b.queryParams) == 0 {
		return b.url
	}
	return b.url + "?" + QueryParamsToStr(b.queryParams)
}

func (b *RequestBuilder) Method(method string) *RequestBuilder {
	if method == "" {
		b.error = ErrEmptyMethod
	}
	b.method = method
	return b
}

func (b *RequestBuilder) Url(url string) *RequestBuilder {
	if url == "" {
		b.error = ErrEmptyUrl
	}
	b.url = url
	return b
}

func (b *RequestBuilder) QueryParams(params map[string]any) *RequestBuilder {
	if params != nil {
		b.queryParams = params
	}
	return b
}

func (b *RequestBuilder) Payload(payload RequestPayload) *RequestBuilder {
	if payload != nil {
		buff, err := payload.toBuff()
		if err != nil {
			b.error = err
		}
		b.payload = buff
		b.contentType = payload.contentType()
	}
	return b
}

func (b *RequestBuilder) Headers(headers http.Header) *RequestBuilder {
	b.headers = headers
	return b
}
func newRequestBuilder() *RequestBuilder {
	return &RequestBuilder{}
}
