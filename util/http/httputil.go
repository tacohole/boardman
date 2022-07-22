package httpHelpers

import (
	"fmt"
	"log"
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

	log.Printf("sending request to %s \n", url)

	response, err := client.Do(request)
	if err != nil {
		return response, fmt.Errorf("%s request to %s failed: %d", method, url, response.StatusCode)
	}
	return response, nil
}
