package goxios

import (
	"github.com/stretchr/testify/assert"
	"slices"
	"strings"
	"testing"
)

func TestParamsToStr(t *testing.T) {
	params := make(map[string]interface{})
	params["a"] = 1
	params["b"] = "2"
	params["c"] = []string{"3", "4"}
	qStr := QueryParamsToStr(params)
	qStrSplt := strings.Split(qStr, "&")
	slices.Sort(qStrSplt)
	expectedQstr := strings.Split("a=1&b=2&c=3&c=4", "&")
	slices.Sort(expectedQstr)
	assert.Equal(t, expectedQstr, qStrSplt)
}

func TestUrlFor(t *testing.T) {
	base := "https://example.com:1010/api/v1/"
	u := "/some/path"
	res, err := urlfor(base, u)
	assert.Nil(t, err)
	assert.Equal(t, "https://example.com:1010/api/v1/some/path", res)
}

func TestNew(t *testing.T) {
	type Typ struct {
		A int
	}
	typ := New[Typ]()
	typ2 := New[*Typ]()
	assert.Equal(t, Typ{}, typ)
	assert.IsType(t, &Typ{}, typ2)
}
