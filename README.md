[![CI/CD Pipeline](https://github.com/Kingson4Wu/saturncli/actions/workflows/go.yml/badge.svg)](https://github.com/Kingson4Wu/saturncli/actions/workflows/go.yml)[![Go Report Card](https://goreportcard.com/badge/github.com/kingson4wu/saturncli)&nbsp;](https://goreportcard.com/report/github.com/kingson4wu/saturncli)![GitHub top language](https://img.shields.io/github/languages/top/kingson4wu/saturncli)&nbsp;[![GitHub stars](https://img.shields.io/github/stars/kingson4wu/saturncli)&nbsp;](https://github.com/kingson4wu/saturncli/stargazers)[![codecov](https://codecov.io/gh/kingson4wu/saturncli/branch/main/graph/badge.svg)](https://codecov.io/gh/kingson4wu/saturncli) [![Go Reference](https://pkg.go.dev/badge/github.com/kingson4wu/saturncli.svg)](https://pkg.go.dev/github.com/kingson4wu/saturncli) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#database) [![LICENSE](https://img.shields.io/github/license/kingson4wu/saturncli.svg?style=flat-square)](https://github.com/kingson4wu/saturncli/blob/main/LICENSE)

English| [简体中文](https://github.com/kingson4wu/saturncli/blob/main/README-CN.md)|[deepwiki](https://deepwiki.com/Kingson4Wu/saturncli)

A command-line client that communicates with a server process via Linux domain sockets. The server communication logic can be directly embedded into server code written in Golang, facilitating the use of shell-type scheduled tasks in the open-source project [Saturn](https://github.com/vipshop/Saturn).

## Design Overview

![](https://github.com/kingson4wu/saturncli/blob/main/resource/img/design-overview-saturn-cli-go.png)

## Quick Start

**1. embedded usage:** see [examples](https://github.com/kingson4wu/saturncli/tree/main/examples) 

**2. command line usage:**
1. `make`
2. `./saturn_svr`
3. `./saturn_cli -name hello -args 'id=33&ver=22'`
4. `./saturn_cli -name hello_stoppable` 
5. `./saturn_cli -name hello_stoppable -stop` OR `CRTL + C` to stop the task `hello_stoppable`

## Documentation

See [wiki](https://github.com/kingson4wu/saturncli/wiki)

## Contributing

If you are interested in contributing to saturncli, see [CONTRIBUTING](https://github.com/kingson4wu/saturncli/blob/main/CONTRIBUTING.md) 

## License

saturncli is licensed under the term of the [Apache 2.0 License](https://github.com/kingson4wu/saturncli/blob/main/LICENSE)

