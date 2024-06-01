//go:build windows

package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Kingson4Wu/saturn_cli_go/base"
)

func (c *cli) Run(task *NotifyTask) string {

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8096/%s?%s", task.Name, task.Args))
	if err != nil {
		c.logger.Errorf("saturn client error, name:%s, args: %s", task.Name, task.Args)
		return base.FAILURE
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Errorf("saturn client error, name:%s, args: %s", task.Name, task.Args)
		return base.FAILURE
	}
	c.logger.Infof("saturn client, name:%s, args: %s, resp: %s", task.Name, task.Args, string(bytes))
	return base.SUCCESS
}
