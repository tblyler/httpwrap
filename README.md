# httpwrap

A simple HTTP server to wrap around shell executions.

# Configuration

See [the example config](./example_config.json) and the [config.Config structure](./config/config.go)

# Install

## Prerequisites

* [Go](https://golang.org) version >= 1.16

* `export PATH="$(go env GOPATH)/bin:$PATH"`

## Steps

1. `go install github.com/tblyler/httpwrap@latest`

2. Set `CONFIG_FILE_PATH` environment variable to the path for your config file. [See reference example config](./example_config.json)

3. Execute `httpwrap`

# TODO

- [ ] Allow execution of httpwrap to wrap around another command execution like... `httpwrap /bin/my-non-http-program`. Forwarding `kill` signals, `STDOUT`, and `STDERR` appropriately.

- [ ] Add proper logging & error handling

- [ ] Add support for more config formats

- [ ] Add Makefile

- [ ] Add CI pipeline

- [ ] Make artifacts available for download in "Releases"

- [ ] Add unit tests

- [X] Stop judging me, I wrote this for my Raspberry Pi