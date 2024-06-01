package main

import (
	saturn "github.com/Kingson4Wu/saturn_cli_go/client"
	"github.com/Kingson4Wu/saturn_cli_go/utils"
)

func main() {
	saturn.NewCmd(&utils.DefaultLogger{},
		"/tmp/notify.sock").Run()

	//go build  -o saturn_cli ./examples/client/client.go
	//./saturn_cli -name hello -args 'id=33&ver=22'

}
