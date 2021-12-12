package httpHelpers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Endpoint struct {
	PageNumber int `json:"current_page"`
	PageSize   int `json:"per_page"`
	NextPage   int `json:"next_page"`
}

const (
	BaseUrl     = "https://www.balldontlie.io/api/v1/"
	Players     = "players"
	Teams       = "teams"
	Stats       = "stats"
	httpTimeout = 30 * time.Second
)

func MakeHttpRequest(method string, url string, data []byte, token string) (*http.Response, error) {
	client := http.Client{
		Timeout: httpTimeout,
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
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

func setHeaders(request *http.Request, headers map[string]string) error {
	if request == nil {
		return errors.New("request cannot be nil")
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	return nil
}
