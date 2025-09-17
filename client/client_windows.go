//go:build windows

package client

import "net/http"

func (c *cli) buildHTTPClient() *http.Client {
	return &http.Client{Timeout: defaultRequestTimeout}
}

func (task *Task) buildURL() (string, error) {
	return task.buildURLForHost("127.0.0.1:8096")
}
