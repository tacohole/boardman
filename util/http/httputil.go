package httpHelpers

import (
	"fmt"
	"net/http"
	"time"
)

const (
	BaseUrl     = "https://balldontlie.io/api/v1/"
	Players     = "players"
	Teams       = "teams"
	Games       = "games/"
	Stats       = "stats/"
	httpTimeout = 300 * time.Second
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
	fmt.Printf("Sending request to %s\n", url)

	response, err := client.Do(request)
	if err != nil || response.StatusCode >= 400 {
		return nil, fmt.Errorf("%s request to %s failed: %s", method, url, err)
	}
	return response, nil
}
