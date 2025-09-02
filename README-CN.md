[![CI/CD Pipeline](https://github.com/Kingson4Wu/saturncli/actions/workflows/go.yml/badge.svg)](https://github.com/Kingson4Wu/saturncli/actions/workflows/go.yml)[![Go Report Card](https://goreportcard.com/badge/github.com/kingson4wu/saturncli)&nbsp;](https://goreportcard.com/report/github.com/kingson4wu/saturncli)![GitHub top language](https://img.shields.io/github/languages/top/kingson4wu/saturncli)&nbsp;[![GitHub stars](https://img.shields.io/github/stars/kingson4wu/saturncli)&nbsp;](https://github.com/kingson4wu/saturncli/stargazers)[![codecov](https://codecov.io/gh/kingson4wu/saturncli/branch/main/graph/badge.svg)](https://codecov.io/gh/kingson4wu/saturncli) [![Go Reference](https://pkg.go.dev/badge/github.com/kingson4wu/saturncli.svg)](https://pkg.go.dev/github.com/kingson4wu/saturncli) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#database) [![LICENSE](https://img.shields.io/github/license/kingson4wu/saturncli.svg?style=flat-square)](https://github.com/kingson4wu/saturncli/blob/main/LICENSE)

[English](https://github.com/kingson4wu/saturncli#saturncli) | 简体中文

## 简介

一个通过Linux域套接字与服务器进程通信的命令行客户端。服务器通信逻辑可以直接嵌入到用Golang编写的服务器代码中，便于在开源项目[Saturn](https://github.com/vipshop/Saturn)中使用shell类型的计划任务。

## 架构设计

![架构设计图](https://github.com/kingson4wu/saturncli/blob/main/resource/img/design-overview-saturn-cli-go.png)

## 快速开始

### 1. 嵌入式使用
参见 [examples](https://github.com/kingson4wu/saturncli/tree/main/examples)

### 2. 命令行使用
1. `make`
2. `./saturn_svr`
3. `./saturn_cli -name hello -args 'id=33&ver=22'`
4. `./saturn_cli -name hello_stoppable`
5. `./saturn_cli -name hello_stoppable -stop` 或者使用 `CRTL + C` 来停止任务 `hello_stoppable`

## 文档

请参考 [wiki](https://github.com/kingson4wu/saturncli/wiki)

## 参与贡献

感谢你的参与，完整的步骤及规范，请参考：[CONTRIBUTING](https://github.com/kingson4wu/saturncli/blob/main/CONTRIBUTING.md)

## License

saturncli 根据 Apache 2.0 License 许可证授权，有关完整许可证文本，请参阅 [LICENSE](https://github.com/kingson4wu/saturncli/blob/main/LICENSE)。