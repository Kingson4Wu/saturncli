package server

import (
	"fmt"
	"strings"
	"testing"
)

func TestTrimPrefix(t *testing.T) {

	name := "/xxx"
	if strings.HasPrefix(name, "/") {
		name = name[1:]
		fmt.Println(name)
	}

	fmt.Println(strings.TrimPrefix("xxx", "/"))
}

// 测试添加重复任务的情况
func TestAddDuplicateJob(t *testing.T) {
	// 添加第一个任务
	err1 := AddJob("duplicate_test", func(m map[string]string, signature string) bool {
		return true
	})

	if err1 != nil {
		t.Errorf("Failed to add first job: %v", err1)
	}

	// 尝试添加同名任务，应该返回错误
	err2 := AddJob("duplicate_test", func(m map[string]string, signature string) bool {
		return true
	})

	if err2 == nil {
		t.Error("Expected error when adding duplicate job, but got none")
	}
}

// 测试添加和停止可停止任务
func TestAddStoppableJob(t *testing.T) {
	err := AddStoppableJob("stoppable_test", func(m map[string]string, signature string, quit chan struct{}) bool {
		// 模拟任务执行
		select {
		case <-quit:
			fmt.Println("Stoppable job received quit signal")
			return true
		default:
			fmt.Println("Stoppable job running")
			return true
		}
	})

	if err != nil {
		t.Errorf("Failed to add stoppable job: %v", err)
	}
}