//go:build !windows

package client

import (
	"context"
	"net"
	"net/http"
)

func (c *cli) buildHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", c.sockPath)
			},
		},
	}
}

func (task *Task) buildUrl() string {
	return "http://unix/" + task.Name + "?" + task.Args
}
