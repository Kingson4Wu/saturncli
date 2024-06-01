package main

import (
	saturn "github.com/Kingson4Wu/saturncli/client"
	"github.com/Kingson4Wu/saturncli/utils"
)

func main() {
	saturn.NewCmd(&utils.DefaultLogger{},
		"/tmp/notify.sock").Run()

	//go build  -o saturn_cli ./examples/client/client.go
	//./saturn_cli -name hello -args 'id=33&ver=22'

}
