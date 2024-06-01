package client

import (
	"context"
	"flag"
	"fmt"
	"github.com/Kingson4Wu/saturncli/base"
	"github.com/Kingson4Wu/saturncli/utils"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"sync"
)

var name = flag.String("name", "", "Input Job Name")
var args = flag.String("args", "", "Input Job Args")
var stop = flag.Bool("stop", false, "Input Job Stop Flag")
var signature = flag.String("signature", "", "Input Job Stop Signature")

func NewCmd(logger utils.Logger, sockPath string) *cmd {
	return &cmd{
		logger:   logger,
		sockPath: sockPath,
	}
}

type cmd struct {
	logger   utils.Logger
	sockPath string
}

func (c *cmd) Run() {

	flag.Usage = func() {
		fmt.Println(`
1. Input Job Name
2. Input Job Args
3. Input Job Stop Flag

Options: 
    `)
		flag.PrintDefaults()
	}

	flag.Parse()

	c.logger.Infof("saturn client cmd task: %s, args:%s", *name, *args)

	result := NewClient(c.logger,
		c.sockPath).Run(&NotifyTask{
		Name:      *name,
		Args:      *args,
		Stop:      *stop,
		Signature: *signature,
	})

	if result == base.SUCCESS {
		fmt.Fprintln(os.Stderr, "Execution Success")
	} else if result == base.INTERRUPT {
		fmt.Fprintln(os.Stderr, "Execution Interrupted")
	} else {
		fmt.Fprintln(os.Stderr, "Execution Failure")
		os.Exit(1)
	}
}

type NotifyTask struct {
	Name      string
	Args      string
	Stop      bool
	Signature string
}

type cli struct {
	logger   utils.Logger
	sockPath string
}

func NewClient(logger utils.Logger, sockPath string) *cli {
	return &cli{
		logger:   logger,
		sockPath: sockPath,
	}
}

func (c *cli) Run(task *NotifyTask) string {

	if task.Name == "" {
		c.logger.Warnf("saturn client run, task name is empty, args:%v", task.Args)
		return base.FAILURE
	}
	c.logger.Infof("saturn client run, task: %v, args: %v", task.Name, task.Args)

	httpc := c.buildHttpClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	url := task.buildUrl()
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

func (c *cli) stop(task *NotifyTask, signature string) {
	httpc := c.buildHttpClient()
	url := task.buildUrl()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
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
	defer response.Body.Close()
	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		c.logger.Errorf("saturn client [stop] read resp body failure from server, task: %s, signature: %s, request server failure, err: %+v", task.Name, signature, err)
		return
	}
	c.logger.Warnf("saturn client [stop] receive result from server, task: %s, signature: %s, resp: %s", task.Name, signature, string(bodyData))
}
