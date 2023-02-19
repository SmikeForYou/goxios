package goxios

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParamsToStr(t *testing.T) {
	params := make(map[string]interface{})
	params["a"] = 1
	params["b"] = "2"
	params["c"] = []string{"3", "4"}
	assert.Equal(t, "a=1&b=2&c=3&c=4", QueryParamsToStr(params))
}

func TestUrlFor(t *testing.T) {
	base := "https://example.com:1010/api/v1/"
	u := "/some/path"
	res, err := urlfor(base, u)
	assert.Nil(t, err)
	assert.Equal(t, "https://example.com:1010/api/v1/some/path", res)
}
