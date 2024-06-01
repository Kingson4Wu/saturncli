//go:build windows

package client

import (
	"fmt"
	"net/http"
)

func (c *cli) buildHttpClient() *http.Client {
	return &http.Client{}
}

func (task *NotifyTask) buildUrl() string {
	return fmt.Sprintf("http://127.0.0.1:8096/%s?%s", task.Name, task.Args)
}
