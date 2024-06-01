[![Go Report Card](https://goreportcard.com/badge/github.com/kingson4wu/saturn_cli_go)&nbsp;](https://goreportcard.com/report/github.com/kingson4wu/saturn_cli_go)![GitHub top language](https://img.shields.io/github/languages/top/kingson4wu/saturn_cli_go)&nbsp;[![GitHub stars](https://img.shields.io/github/stars/kingson4wu/saturn_cli_go)&nbsp;](https://github.com/kingson4wu/saturn_cli_go/stargazers)[![codecov](https://codecov.io/gh/kingson4wu/saturn_cli_go/branch/main/graph/badge.svg)](https://codecov.io/gh/kingson4wu/saturn_cli_go) [![Go Reference](https://pkg.go.dev/badge/github.com/kingson4wu/saturn_cli_go.svg)](https://pkg.go.dev/github.com/kingson4wu/saturn_cli_go) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#database) [![LICENSE](https://img.shields.io/github/license/kingson4wu/saturn_cli_go.svg?style=flat-square)](https://github.com/kingson4wu/saturn_cli_go/blob/main/LICENSE)

English| [简体中文](https://github.com/kingson4wu/saturn_cli_go/blob/main/README-CN.md)

A command-line client that communicates with a server process via Linux domain sockets can be embedded into a Go-written server, facilitating the use of Shell-type scheduled tasks in the open-source project [Saturn](https://github.com/vipshop/Saturn).

## Design Overview

![](https://github.com/kingson4wu/saturn_cli_go/blob/main/resource/img/design-overview-saturn-cli-go.png)

## Quick Start

**1. embedded usage:** see [examples](https://github.com/kingson4wu/saturn_cli_go/tree/main/examples) 

**2. command line usage:**
1. `make`
2. `./saturn_svr`
3. `./saturn_cli -name hello -args 'id=33&ver=22'`
4. `./saturn_cli -name hello_stoppable` 
5. `./saturn_cli -name hello_stoppable -stop` OR `CRTL + C` to stop the task `hello_stoppable`

## Documentation

See [wiki](https://github.com/kingson4wu/saturn_cli_go/wiki)

## Contributing

If you are interested in contributing to saturn_cli_go, see [CONTRIBUTING](https://github.com/kingson4wu/saturn_cli_go/blob/main/CONTRIBUTING.md) 

## License

saturn_cli_go is licensed under the term of the [Apache 2.0 License](https://github.com/kingson4wu/saturn_cli_go/blob/main/LICENSE)

