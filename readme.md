## GEX In-Terminal Gaia Explorer

GEX is a real time in-terminal explorer for Cosmos SDK blockchains. See the [Check out your Cosmos SDK blockchain in a terminal-based block explorer](https://blog.cosmos.network/gaia-explorer-in-terminal-f37a4ea52e3c) blog post to learn more about GEX.

The GEX explorer displays blocks, transactions, validator, network status, and more information. Use the GEX block explorer to see the status of peers, connection, version, and other useful information to have a quick peek into your own node. GEX works with Ignite.

## Install GEX

The GEX installation requires Go. If you don't already have Go installed, see https://golang.org/dl. Download the binary release that is suitable for your system and follow the installation instructions.

To install the GEX binary:

```shell
go install github.com/ignite/gex@latest
```

## Run GEX

To launch a GEX explorer in your terminal window, type:

```shell
gex explorer
```

and hit enter.

## Optional Host

Configure an optional host, instead of using the default RPC host `http://localhost:26657`

```shell
gex explorer 192.168.0.1
```

## Optional Port

Configure an optional port

```shell
gex explorer 192.168.0.1:27657
```

## Print help
```shell
Usage:
  gex [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  explorer    Gex is a cosmos explorer for terminals
  help        Help about any command
  version     Print the current build information

Flags:
  -h, --help   help for gex

Use "gex [command] --help" for more information about a command.

```

## Preview

![Terminal Screenshot](./screenshot.png "Screenshot Application")

## Run In Development

To manually run GEX, clone the `github.com/ignite/gex` repository and then cd into the `gex` directory. Then to run GEX manually, type this command in a terminal window:

`go run main.go explorer`

## Contribute

See [CONTRIBUTING.md](./CONTRIBUTING.md) to learn about how to contribute and how the code is structured.

Thanks for contributing!
