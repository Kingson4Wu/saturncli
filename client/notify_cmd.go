package client

import (
	"flag"
	"fmt"
	"github.com/Kingson4Wu/saturn_cli_go/base"
	"github.com/Kingson4Wu/saturn_cli_go/utils"
	"os"
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
