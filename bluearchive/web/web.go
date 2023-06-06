package web

import (
	"errors"
	"net/http"
	"time"
)

func MakeRequest(url string, headers map[string]string) (*http.Response, error) {
	cli := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = cli.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
	}

	if resp == nil {
		return nil, errors.New("response is nil")
	}

	return resp, nil
}
