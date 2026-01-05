package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/neovim/go-client/nvim"
)

func main() {
	os.Exit(run())
}

func run() int {
	addr := os.Getenv("NVIM")
	if addr == "" {
		addr = os.Getenv("NVIM_LISTEN_ADDRESS")
	}
	if addr == "" {
		fmt.Fprintln(os.Stderr, "NVIM or NVIM_LISTEN_ADDRESS environment variable must be set")
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
		"openwin": func() (cli.Command, error) {
			return openWinCommand{
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
	return "execute command"
}

func (exCommand) Help() string {
	return `Execute arguments as NeoVim command. No need to prefix ':'.

EXAMPLES:

	$ nvc ex echo \'abc\'

	It executes ":echo 'abc'" in NeoVim so that NeoVim shows "abc" in status line.

	$ nvc ex e a.txt

	It executes ":e a.txt" in NeoVim so that NeoVim opens "a.txt" in current buffer.
`
}

func (c exCommand) Run(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "missing command to execute as argument")
		return 1
	}

	out, err := c.nv.Exec(strings.Join(args, " "), true)
	fmt.Print(out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

type openWinCommand struct {
	nv *nvim.Nvim
}

func (openWinCommand) Synopsis() string { return "open window" }

func (openWinCommand) Help() string {
	return `Open Window. See https://neovim.io/doc/user/api.html#nvim_open_win() for detail.
The width and height should be specified as WxH.

EXAMPLES:

	$ nvc openwin 30x5 --relative cursor --row 3 --col 3
`
}

func (c openWinCommand) Run(args []string) int {
	fs := flag.NewFlagSet("open-window", flag.ExitOnError)
	var bufN, winN int
	var enter, focusable, external bool
	var relative, anchor string
	var row, col int
	fs.IntVar(&bufN, "buffer", 0, "")
	fs.BoolVar(&enter, "enter", true, "")
	fs.StringVar(&relative, "relative", "", "")
	fs.IntVar(&winN, "win", 0, "")
	fs.StringVar(&anchor, "anchor", "", "")
	fs.IntVar(&row, "row", 0, "")
	fs.IntVar(&col, "col", 0, "")
	fs.BoolVar(&focusable, "focusable", false, "")
	fs.BoolVar(&external, "external", false, "")
	var pargs []string
	for {
		if err := fs.Parse(args); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid argument: %v\n", err)
			return 1
		}
		if fs.NArg() == 0 {
			break
		}
		pargs = append(pargs, args[0])
		args = args[1:]
	}

	var buf nvim.Buffer
	if isSet(fs, "buffer") {
		var err error
		buf, err = c.nv.CurrentBuffer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get the current buffer: %v\n", err)
		}
	} else {
		buf = nvim.Buffer(bufN)
	}

	if len(pargs) < 1 {
		fmt.Fprintln(os.Stderr, "Please specify the width and height as WxH")
		return 1
	}
	if len(pargs) > 1 {
		fmt.Fprintln(os.Stderr, "Too many positional arguments")
		return 1
	}
	var width, height int
	if _, err := fmt.Sscanf(pargs[0], "%dx%d", &width, &height); err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse the width and height: %v\n", err)
		return 1
	}

	config := new(nvim.WindowConfig)
	config.Width = width
	config.Height = height
	if isSet(fs, "relative") {
		config.Relative = relative
	}
	if isSet(fs, "win") {
		config.Win = nvim.Window(winN)
	}
	if isSet(fs, "anchor") {
		config.Anchor = anchor
	}
	if isSet(fs, "row") {
		config.Row = row
	}
	if isSet(fs, "col") {
		config.Col = col
	}
	if isSet(fs, "focusable") {
		config.Focusable = focusable
	}
	if isSet(fs, "external") {
		config.External = external
	}
	w, err := c.nv.OpenWindow(buf, enter, config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Printf("Window %d is successfully opened.\n", w)
	return 0
}

func isSet(fs *flag.FlagSet, name string) bool {
	var set bool
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			set = true
		}
	})
	return set
}
