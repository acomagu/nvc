package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/neovim/go-client/nvim"
	"github.com/mitchellh/cli"
)

func main() {
	os.Exit(run())
}

func run() int {
	addr := os.Getenv("NVIM_LISTEN_ADDRESS")
	if addr == "" {
		fmt.Fprintln(os.Stderr, "NVIM_LISTEN_ADDRESS not set")
		return 1
	}

	nv, err := nvim.Dial(addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer nv.Close()

	c := cli.NewCLI("nvc", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"ex": func() (cli.Command, error) {
			return exCommand{
				nv: nv,
			}, nil
		},
	}

	es, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return es
}

type exCommand struct {
	nv *nvim.Nvim
}

func (exCommand) Synopsis() string {
	return ""
}

func (exCommand) Help() string {
	return ""
}

func (c exCommand) Run(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "missing command to execute as argument")
		return 1
	}

	err := c.nv.Command(strings.Join(args, " "))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
