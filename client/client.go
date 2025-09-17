//go:build !windows

package client

import (
	"context"
	"net"
	"net/http"
)

func (c *cli) buildHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultRequestTimeout,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, "unix", c.sockPath)
			},
		},
	}
}

func (task *Task) buildURL() (string, error) {
	return task.buildURLForHost("unix")
}
