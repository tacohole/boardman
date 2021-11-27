package httputil_test

import (
	"boardman/util/httputil"
	"testing"
)

const (
	method = "GET"
	getUrl = "http://www.google.com"
)

func testMakeHttpRequest(t *testing.T) {
	_, err := httputil.MakeHttpRequest(method, getUrl, nil, "")
	if err != nil {
		t.Fail()
	}
}
