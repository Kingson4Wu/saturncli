package client_test

import (
	"fmt"
	"github.com/Kingson4Wu/saturn_cli_go/client"
	"github.com/Kingson4Wu/saturn_cli_go/server"
	"github.com/Kingson4Wu/saturn_cli_go/utils"
	"os"
	"testing"
	"time"
)

func TestNewCmd(t *testing.T) {

	go func() {
		_ = server.AddJob("hello", func(m map[string]string, signature string) bool {
			return true
		})
		server.NewServer(&utils.DefaultLogger{},
			"/tmp/notify.sock").Serve()
	}()

	time.Sleep(3 * time.Second)

	// -name hello
	os.Args = append(os.Args, "-name")
	os.Args = append(os.Args, "hello")
	client.NewCmd(&utils.DefaultLogger{},
		"/tmp/notify.sock").Run()
}

/*func TestNewStoppableServerJob(t *testing.T) {

}*/

func TestNewStoppableJob(t *testing.T) {

	go func() {
		_ = server.AddStoppableJob("hello_stoppable", func(m map[string]string, signature string, quit chan struct{}) bool {
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
		})
		server.NewServer(&utils.DefaultLogger{},
			"/tmp/notify.sock").Serve()
	}()

	time.Sleep(3 * time.Second)

	// -name hello
	os.Args = append(os.Args, "-name")
	os.Args = append(os.Args, "hello_stoppable")
	client.NewCmd(&utils.DefaultLogger{},
		"/tmp/notify.sock").Run()
}

/*func TestStoppableJob(t *testing.T) {

	// -name hello
	os.Args = append(os.Args, "-name")
	os.Args = append(os.Args, "hello_stoppable")
	os.Args = append(os.Args, "-stop")
	//os.Args = append(os.Args, "-signature")
	//os.Args = append(os.Args, "9a675e22-1d92-11ef-bdf5-9e40e26b9695")
	saturn.NewCmd(&utils.DefaultLogger{},
		"/tmp/notify.sock").Run()

	//./saturn_cli -name hello_stoppable -stop
}*/
