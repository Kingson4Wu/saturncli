package server_test

import (
	"github.com/Kingson4Wu/saturn_cli_go/client"
	"github.com/Kingson4Wu/saturn_cli_go/server"
	"github.com/Kingson4Wu/saturn_cli_go/utils"
	"testing"
	"time"
)

/*func TestNotifyServe(t *testing.T) {
	saturn.AddJob("hello", func(m map[string]string, signature string) bool {
		return true
	})
	saturn.NewServer(&utils.DefaultLogger{},
		"/tmp/notify.sock").Serve()
}*/

func BenchmarkNotifyServe(b *testing.B) {

	socketPath := "/tmp/notify.sock"
	_ = server.AddJob("hello", func(m map[string]string, signature string) bool {
		time.Sleep(50 * time.Millisecond)
		return true
	})
	go server.NewServer(&utils.DefaultLogger{},
		socketPath).Serve()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		client.NewClient(&utils.DefaultLogger{},
			socketPath).Run(&client.NotifyTask{
			Name: "hello",
			Args: "id=45&tel=48893",
		})
	}

	//b.ReportMetric(float64(b.Elapsed())/1000000, "ms/op")

	//BenchmarkNotifyServe-8   	      25	  51485200 ns/op	   31452 B/op	     163 allocs/op

}

func BenchmarkParallelNotifyServe(b *testing.B) {

	socketPath := "/tmp/notify.sock"
	_ = server.AddJob("hello", func(m map[string]string, signature string) bool {
		time.Sleep(50 * time.Millisecond)
		return true
	})
	go server.NewServer(&utils.DefaultLogger{},
		socketPath).Serve()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client.NewClient(&utils.DefaultLogger{},
				socketPath).Run(&client.NotifyTask{
				Name: "hello",
				Args: "id=45&tel=48893",
			})
		}
	})

	//BenchmarkParallelNotifyServe-8   	     242	   6411410 ns/op	   30432 B/op	     159 allocs/op

}
