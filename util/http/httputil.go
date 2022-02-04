package httpHelpers

import (
	"fmt"
	"net/http"
	"time"
)

const (
	httpTimeout = 15 * time.Second
)

func MakeHttpRequest(method string, url string) (*http.Response, error) {
	client := http.Client{
		Timeout: httpTimeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	fmt.Printf("sending request to %s", url)
	response, err := client.Do(request)
	if err != nil || response.StatusCode >= 400 {
		return nil, fmt.Errorf("%s request to %s failed: %d", method, url, response.StatusCode)
	}
	return response, nil
}
