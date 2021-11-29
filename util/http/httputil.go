package httpHelpers

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultApiUrl = "http://localhost:3000"
	httpTimeout   = 1 * time.Minute
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

func GetApiUrl() string {
	apiUrl := viper.GetString("API_URL")

	if apiUrl == defaultApiUrl {
		log.Printf("Warning: default API URL in use %s", defaultApiUrl)
	}

	return apiUrl
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
