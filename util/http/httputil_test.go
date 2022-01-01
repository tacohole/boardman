package httpHelpers

import (
	"testing"
)

const (
	method = "GET"
	getUrl = "http://www.google.com"
)

func TestMakeHttpRequest(t *testing.T) {
	_, err := MakeHttpRequest(method, getUrl)
	if err != nil {
		t.Fail()
	}
}
