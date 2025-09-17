package client

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Kingson4Wu/saturncli/base"
	"github.com/Kingson4Wu/saturncli/utils"
	"io"
	"os"
	"sort"
	"strings"
)

// NewCmd constructs a CLI command wrapper bound to the provided logger and socket path.
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
	c.RunWithArgs(os.Args[1:])
}

func (c *cmd) RunWithArgs(arguments []string) {
	opts, err := c.parse(arguments)
	if err != nil {
		c.logger.Errorf("saturn client parse arguments failure: %+v", err)
		fmt.Fprintln(os.Stderr, "Execution Failure")
		os.Exit(1)
		return
	}
	if opts == nil {
		return
	}

	c.logger.Infof("saturn client cmd task: %s, args:%s, params:%v", opts.name, opts.args, opts.params)

	result := NewClient(c.logger,
		c.sockPath).Run(&Task{
		Name:      opts.name,
		Args:      opts.args,
		Stop:      opts.stop,
		Signature: opts.signature,
		Params:    cloneStringMap(opts.params),
	})

	switch result {
	case base.SUCCESS:
		fmt.Fprintln(os.Stderr, "Execution Success")
	case base.INTERRUPT:
		fmt.Fprintln(os.Stderr, "Execution Interrupted")
	default:
		fmt.Fprintln(os.Stderr, "Execution Failure")
		os.Exit(1)
	}
}

type cmdOptions struct {
	name      string
	args      string
	stop      bool
	signature string
	params    map[string]string
}

func (c *cmd) parse(arguments []string) (*cmdOptions, error) {
	opts := &cmdOptions{}
	fs := flag.NewFlagSet("saturn-cli", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	fs.StringVar(&opts.name, "name", "", "Input Job Name")
	fs.StringVar(&opts.args, "args", "", "Input Job Args")
	fs.BoolVar(&opts.stop, "stop", false, "Input Job Stop Flag")
	fs.StringVar(&opts.signature, "signature", "", "Input Job Stop Signature")
	var paramFlag keyValueFlag
	fs.Var(&paramFlag, "param", "Key=Value pair to include in request; can be repeated")

	usage := func() {
		fmt.Fprint(os.Stderr, `
1. Input Job Name
2. Input Job Args
3. Input Job Stop Flag

Options:
`)
		fs.SetOutput(os.Stderr)
		fs.PrintDefaults()
		fs.SetOutput(io.Discard)
	}
	fs.Usage = usage

	if err := fs.Parse(arguments); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			usage()
			return nil, nil
		}
		usage()
		return nil, err
	}

	if fs.NArg() > 0 {
		usage()
		return nil, fmt.Errorf("unexpected arguments: %v", fs.Args())
	}

	opts.params = cloneStringMap(paramFlag.values)

	return opts, nil
}

// keyValueFlag collects repeated --param flags into a map.
type keyValueFlag struct {
	values map[string]string
}

func (f *keyValueFlag) String() string {
	if f == nil || len(f.values) == 0 {
		return ""
	}
	parts := make([]string, 0, len(f.values))
	for k, v := range f.values {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func (f *keyValueFlag) Set(value string) error {
	if value == "" {
		return errors.New("param requires key=value")
	}
	pieces := strings.SplitN(value, "=", 2)
	if len(pieces) != 2 {
		return fmt.Errorf("invalid param %q, expect key=value", value)
	}
	key := strings.TrimSpace(pieces[0])
	if key == "" {
		return fmt.Errorf("invalid param %q, key is empty", value)
	}
	if f.values == nil {
		f.values = make(map[string]string)
	}
	f.values[key] = pieces[1]
	return nil
}

// cloneStringMap copies a map to avoid leaking shared mutable state across requests.
func cloneStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	out := make(map[string]string, len(input))
	for k, v := range input {
		out[k] = v
	}
	return out
}
