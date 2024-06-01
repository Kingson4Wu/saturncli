//go:build !windows

package client

import (
	"context"
	"github.com/Kingson4Wu/saturn_cli_go/base"
	"github.com/Kingson4Wu/saturn_cli_go/utils"
	"github.com/google/uuid"
	"io"
	"net"
	"net/http"
	"sync"
)

func (c *cli) Run(task *NotifyTask) string {

	if task.Name == "" {
		c.logger.Warnf("saturn client run, task name is empty, args:%v", task.Args)
		return base.FAILURE
	}
	c.logger.Infof("saturn client run, task: %v, args: %v", task.Name, task.Args)

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", c.sockPath)
			},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	url := "http://unix/" + task.Name + "?" + task.Args
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.logger.Errorf("saturn client create request failure, task: %s, args:%s, err: %+v", task.Name, task.Args, err)
		return base.FAILURE
	}
	runSignature := ""
	if v, err := uuid.NewUUID(); err == nil {
		runSignature = v.String()
	}

	if task.Stop {
		addStopOption(req, task.Signature)
	} else {
		if runSignature != "" {
			req.Header.Set(base.RunSignature, runSignature)
		}
	}

	var wg sync.WaitGroup
	var interrupt bool
	requestFinishChan := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		signalChan := utils.ListenSignal()
		select {
		case <-requestFinishChan:
			c.logger.Infof("saturn client listen signal, response finish : %s, signature: %s, args:%s", task.Name, runSignature, task.Args)
		case signal := <-signalChan:
			c.logger.Warnf("saturn client listen signal: %s, request interrupt : %s, signature: %s, args:%s", signal, task.Name, runSignature, task.Args)
			if !task.Stop && runSignature != "" {
				c.stop(task.Name, runSignature)
			}
			interrupt = true
			cancel()
		}
	}()
	var (
		response *http.Response
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(requestFinishChan)
		response, err = httpc.Do(req)
	}()
	wg.Wait()

	if interrupt {
		return base.INTERRUPT
	}

	//response, err := httpc.Do(req)
	if err != nil {
		c.logger.Errorf("saturn client fail to request server, task: %s, signature: %s, args:%s, err: %+v", task.Name, runSignature, task.Args, err)
		return base.FAILURE
	}

	defer response.Body.Close()
	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Errorf("saturn client read resp body failure from server, task: %s, signature: %s, args:%se, err: %+v", task.Name, runSignature, task.Args, err)
		return base.FAILURE
	}
	c.logger.Infof("saturn client receive result from server, task: %s, signature: %s, args:%s, resp: %s", task.Name, runSignature, task.Args, string(bodyData))
	return string(bodyData)
}

func addStopOption(req *http.Request, signature string) {
	req.Header.Set(base.StopJobFlag, "true")
	if signature != "" {
		req.Header.Set(base.StopSignature, signature)
	}
}

func (c *cli) stop(taskName, signature string) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", c.sockPath)
			},
		},
	}
	url := "http://unix/" + taskName
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		c.logger.Errorf("saturn client [stop] create request failure, task: %s, signature: %s, request server failure, err: %+v", taskName, signature, err)
		return
	}
	addStopOption(req, signature)
	response, err := httpc.Do(req)
	if err != nil {
		c.logger.Errorf("saturn client [stop] receive result from server, task: %s, signature: %s, request server failure, err: %+v", taskName, signature, err)
		return
	}
	defer response.Body.Close()
	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Errorf("saturn client [stop] read resp body failure from server, task: %s, signature: %s, request server failure, err: %+v", taskName, signature, err)
		return
	}
	c.logger.Warnf("saturn client [stop] receive result from server, task: %s, signature: %s, resp: %s", taskName, signature, string(bodyData))
}
