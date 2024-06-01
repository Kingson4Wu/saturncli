package main

import (
	"fmt"
	"github.com/Kingson4Wu/saturn_cli_go/server"
	"github.com/Kingson4Wu/saturn_cli_go/utils"
	"time"
)

func main() {
	if err := server.AddJob("hello", func(m map[string]string, signature string) bool {
		for k, v := range m {
			fmt.Printf("%v: %v\n", k, v)
		}
		return true
	}); err != nil {
		panic(err)
	}

	if err := server.AddStoppableJob("hello_stoppable", func(m map[string]string, signature string, quit chan struct{}) bool {
		list := []int{1, 2, 3, 4, 5}
		for _, value := range list {
			select {
			case <-quit:
				fmt.Println("Received quit signal. Exiting loop.")
				return true
			default:
				fmt.Printf("Processing value :%v, signature: %v \n", value, signature)
				time.Sleep(3 * time.Second)
			}
		}
		return true
	}); err != nil {
		panic(err)
	}

	server.NewServer(&utils.DefaultLogger{},
		"/tmp/notify.sock").Serve()

}
