package httpHelpers

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseUrl     = "https://www.balldontlie.io/api/v1/"
	Players     = "players/"
	Teams       = "teams/"
	Games       = "games/"
	Stats       = "stats/"
	httpTimeout = 30 * time.Second
)

func MakeHttpRequest(method string, url string, data []byte, token string) (*http.Response, error) {
	client := http.Client{
		Timeout: httpTimeout,
	}

	request, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content Type", "application/json")
	if token != "" {
		request.Header.Set("Authorization", "Bearer"+token)
	}

	response, err := client.Do(request)
	if err != nil || response.StatusCode <= 400 {
		return response, fmt.Errorf("%s request to %s failed with status %d", method, url, response.StatusCode)
	}
	return response, nil
}
