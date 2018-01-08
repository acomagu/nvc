package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/neovim/go-client/nvim"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("nvc", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"ex": exCommandFactory,
	}

	es, err := c.Run()
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(es)
}

type exCommand struct {}

func (exCommand) Synopsis() string {
	return ""
}

func (exCommand) Help() string {
	return ""
}

func (exCommand) Run(args []string) int {
	addr := os.Getenv("NVIM_LISTEN_ADDRESS")
	if addr == "" {
		fmt.Println("NVIM_LISTEN_ADDRESS not set")
		return 1
	}

	v, err := nvim.Dial(addr)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer v.Close()

	if len(args) < 2 {
		fmt.Println("missing command to execute as argument")
		return 1
	}

	res, err := v.CommandOutput(strings.Join(args, " "))
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Println(res)
	return 0
}

func exCommandFactory() (cli.Command, error) {
	return exCommand{}, nil
}
