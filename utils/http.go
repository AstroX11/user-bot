package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 15 * time.Second}

var commonHeaders = map[string]string{
	"DNT":                       "1",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
	"Accept":                    "application/json",
}

func GetBuffer(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	setCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func GetJson(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	setCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("GET %s failed: %s", url, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func PostJson(url string, payload any, target any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("POST %s failed: %s", url, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func setCommonHeaders(req *http.Request) {
	for k, v := range commonHeaders {
		req.Header.Set(k, v)
	}
}
