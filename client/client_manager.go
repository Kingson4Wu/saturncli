package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kingson4Wu/saturncli/base"
	"github.com/Kingson4Wu/saturncli/utils"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Task describes the parameters required to invoke a remote job.
type Task struct {
	Name      string
	Args      string
	Params    map[string]string
	Stop      bool
	Signature string
}

type cli struct {
	logger   utils.Logger
	sockPath string
}

const (
	defaultRequestTimeout = 60 * time.Second
	stopRequestTimeout    = 10 * time.Second
)

// NewClient constructs a client capable of communicating with the Saturn server over the provided socket path.
func NewClient(logger utils.Logger, sockPath string) *cli {
	return &cli{
		logger:   logger,
		sockPath: sockPath,
	}
}

func (c *cli) Run(task *Task) string {

	if task == nil {
		c.logger.Errorf("saturn client run received nil task")
		return base.FAILURE
	}

	if task.Name == "" {
		c.logger.Warnf("saturn client run, task name is empty, args:%v", task.Args)
		return base.FAILURE
	}
	c.logger.Infof("saturn client run, task: %v, args: %v, params: %v", task.Name, task.Args, task.Params)

	requestURL, err := task.buildURL()
	if err != nil {
		c.logger.Errorf("saturn client build url failure, task: %s, args:%s, err: %+v", task.Name, task.Args, err)
		return base.FAILURE
	}

	httpc := c.buildHTTPClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
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
		defer utils.StopSignal(signalChan)
		select {
		case <-requestFinishChan:
			c.logger.Infof("saturn client listen signal, response finish : %s, signature: %s, args:%s", task.Name, runSignature, task.Args)
		case signal := <-signalChan:
			if signal == nil {
				return
			}
			c.logger.Warnf("saturn client listen signal: %s, request interrupt : %s, signature: %s, args:%s", signal, task.Name, runSignature, task.Args)
			if !task.Stop && runSignature != "" {
				c.stop(task, runSignature)
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

	defer func() {
		if err := response.Body.Close(); err != nil {
			c.logger.Warnf("saturn client failed to close response body: %v", err)
		}
	}()
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

func (c *cli) stop(task *Task, signature string) {
	requestURL, err := task.buildURL()
	if err != nil {
		c.logger.Errorf("saturn client [stop] build url failure, task: %s, signature: %s, err: %+v", task.Name, signature, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), stopRequestTimeout)
	defer cancel()

	httpc := c.buildHTTPClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		c.logger.Errorf("saturn client [stop] create request failure, task: %s, signature: %s, request server failure, err: %+v", task.Name, signature, err)
		return
	}
	addStopOption(req, signature)
	response, err := httpc.Do(req)
	if err != nil {
		c.logger.Errorf("saturn client [stop] receive result from server, task: %s, signature: %s, request server failure, err: %+v", task.Name, signature, err)
		return
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			c.logger.Warnf("saturn client [stop] failed to close response body: %v", err)
		}
	}()
	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Errorf("saturn client [stop] read resp body failure from server, task: %s, signature: %s, request server failure, err: %+v", task.Name, signature, err)
		return
	}
	c.logger.Warnf("saturn client [stop] receive result from server, task: %s, signature: %s, resp: %s", task.Name, signature, string(bodyData))
}

func (task *Task) buildURLForHost(host string) (string, error) {
	if task == nil {
		return "", errors.New("task is nil")
	}
	if strings.TrimSpace(task.Name) == "" {
		return "", errors.New("task name is empty")
	}

	query, err := task.queryString()
	if err != nil {
		return "", err
	}

	path := "/" + url.PathEscape(task.Name)
	u := url.URL{
		Scheme:   "http",
		Host:     host,
		Path:     path,
		RawQuery: query,
	}
	return u.String(), nil
}

func (task *Task) queryString() (string, error) {
	values := url.Values{}
	for k, v := range task.Params {
		values.Set(k, v)
	}

	raw := strings.TrimPrefix(task.Args, "?")
	if raw != "" {
		parsed, err := url.ParseQuery(raw)
		if err != nil {
			return "", fmt.Errorf("invalid args: %w", err)
		}
		for k, v := range parsed {
			if len(v) == 0 {
				continue
			}
			if _, exists := values[k]; exists {
				continue
			}
			values.Set(k, v[0])
		}
	}

	if len(values) == 0 {
		return "", nil
	}

	return values.Encode(), nil
}
