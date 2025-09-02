package client_test

import (
	"fmt"
	"github.com/Kingson4Wu/saturncli/client"
	"github.com/Kingson4Wu/saturncli/server"
	"github.com/Kingson4Wu/saturncli/utils"
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
	// 由于flag.Parse()不是并发安全的，我们需要串行执行这个测试
	// 先启动服务器
	go func() {
		_ = server.AddStoppableJob("test_stoppable", func(m map[string]string, signature string, quit chan struct{}) bool {
			for i := 0; i < 10; i++ {
				select {
				case <-quit:
					fmt.Printf("Job %s with signature %s stopped\n", "test_stoppable", signature)
					return true
				default:
					fmt.Printf("Processing step %d for job %s with signature %s\n", i, "test_stoppable", signature)
					time.Sleep(1 * time.Second)
				}
			}
			return true
		})
		server.NewServer(&utils.DefaultLogger{},
			"/tmp/notify_stop.sock").Serve()
	}()

	time.Sleep(3 * time.Second)

	// 保存原始的os.Args
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs // 恢复原始参数
	}()

	// 启动一个长时间运行的任务
	os.Args = []string{"cmd", "-name", "test_stoppable"}
	client.NewCmd(&utils.DefaultLogger{},
		"/tmp/notify_stop.sock").Run()

	// 等待任务开始执行
	time.Sleep(2 * time.Second)

	// 发送停止信号
	os.Args = []string{"cmd", "-name", "test_stoppable", "-stop"}
	client.NewCmd(&utils.DefaultLogger{},
		"/tmp/notify_stop.sock").Run()
}