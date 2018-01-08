# nvc: Neovim Remote Client

A CLI to control NeoVim via its RPC API.

Inspired by [mhinz/neovim-remote](https://github.com/mhinz/neovim-remote).

## Usage

Now we have only `ex` sub-command, which runs command in ex-mode and prints the output. No need to prefix `:`.

To show message in command line,

```
$ nvc ex echo \'abc\'
```

To open file in current buffer,

```
$ nvc ex e a.txt
```

(Note that you must pass relative path from NeoVim's current directory or absolute path)

The `NVIM_LISTEN_ADDRESS` environment variable must be set. In terminal mode of NeoVim, it should set normally.

## Installation

```
$ go get -u github.com/acomagu/nvc
```

## What is the advantage over [neovim-remote](https://github.com/mhinz/neovim-remote)?

Speed.

Comparison of running `:echo 'abc'` as ex-command from CLI:

|               | Average Time |    Ratio   |
|---------------|--------------|------------|
| nvc           |       0.007s |         1x |
| neovim-remote |       0.494s |      70.6x |

**70x faster!**

## Contributing

Welcome!

We have a lot of missing functions.

## Author

[acomagu](https://github.com/acomagu)
