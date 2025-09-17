package client_test

import (
	"fmt"
	"github.com/Kingson4Wu/saturncli/client"
	"github.com/Kingson4Wu/saturncli/server"
	"github.com/Kingson4Wu/saturncli/utils"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestNewCmd(t *testing.T) {
	registry := server.NewRegistry()
	result := make(chan map[string]string, 1)
	if err := registry.AddJob("hello", func(m map[string]string, signature string) bool {
		result <- m
		return true
	}); err != nil {
		t.Fatalf("failed to add job: %v", err)
	}

	socket := tempSocketPath(t, "notify")
	go server.NewServer(&utils.DefaultLogger{}, socket, server.WithRegistry(registry)).Serve()

	time.Sleep(300 * time.Millisecond)

	client.NewCmd(&utils.DefaultLogger{},
		socket).RunWithArgs([]string{"-name", "hello", "-param", "id=42", "-param", "foo=bar"})

	select {
	case args := <-result:
		if args["id"] != "42" {
			t.Fatalf("expected id=42, got %v", args["id"])
		}
		if args["foo"] != "bar" {
			t.Fatalf("expected foo=bar, got %v", args["foo"])
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for job to execute")
	}
}

/*func TestNewStoppableServerJob(t *testing.T) {

}*/

func TestNewStoppableJob(t *testing.T) {
	registry := server.NewRegistry()
	if err := registry.AddStoppableJob("hello_stoppable", func(m map[string]string, signature string, quit chan struct{}) bool {
		list := []int{1, 2, 3}
		for _, value := range list {
			select {
			case <-quit:
				fmt.Println("Received quit signal. Exiting loop.")
				return true
			default:
				fmt.Printf("Processing value :%v, signature: %v \n", value, signature)
				time.Sleep(200 * time.Millisecond)
			}
		}
		return true
	}); err != nil {
		t.Fatalf("failed to add stoppable job: %v", err)
	}

	socket := tempSocketPath(t, "notify-stoppable")
	go server.NewServer(&utils.DefaultLogger{}, socket, server.WithRegistry(registry)).Serve()

	time.Sleep(300 * time.Millisecond)

	client.NewCmd(&utils.DefaultLogger{},
		socket).RunWithArgs([]string{"-name", "hello_stoppable"})
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

// 新增测试用例：测试任务名称为空的情况
// 注意：这个测试会导致调用os.Exit(1)，因此不能直接运行
// func TestNewCmdWithEmptyName(t *testing.T) {
// 	// 重置命令行参数
// 	os.Args = []string{"cmd"}
//
// 	// 添加空的任务名称
// 	os.Args = append(os.Args, "-name")
// 	os.Args = append(os.Args, "") // 空名称
//
// 	// 捕获输出
// 	// 这里应该输出警告信息并返回FAILURE
// 	client.NewCmd(&utils.DefaultLogger{},
// 		"/tmp/notify.sock").Run()
// }

// 新增测试用例：测试停止任务功能
func TestStopJob(t *testing.T) {
	registry := server.NewRegistry()
	if err := registry.AddStoppableJob("test_stoppable", func(m map[string]string, signature string, quit chan struct{}) bool {
		for i := 0; i < 10; i++ {
			select {
			case <-quit:
				fmt.Printf("Job %s with signature %s stopped\n", "test_stoppable", signature)
				return true
			default:
				fmt.Printf("Processing step %d for job %s with signature %s\n", i, "test_stoppable", signature)
				time.Sleep(120 * time.Millisecond)
			}
		}
		return true
	}); err != nil {
		t.Fatalf("failed to add stoppable job: %v", err)
	}

	socket := tempSocketPath(t, "notify-stop")

	go server.NewServer(&utils.DefaultLogger{}, socket, server.WithRegistry(registry)).Serve()

	time.Sleep(300 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		client.NewCmd(&utils.DefaultLogger{},
			socket).RunWithArgs([]string{"-name", "test_stoppable"})
	}()

	time.Sleep(400 * time.Millisecond)

	client.NewCmd(&utils.DefaultLogger{},
		socket).RunWithArgs([]string{"-name", "test_stoppable", "-stop"})

	wg.Wait()
}

func tempSocketPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(os.TempDir(), fmt.Sprintf("saturncli-%s-%d.sock", name, time.Now().UnixNano()))
}
