package utils

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {

	fmt.Println(string(Stack(3)))

}

// 测试Stack函数在不同深度下的行为
func TestStackWithDifferentDepths(t *testing.T) {
	// 测试深度为1
	stack1 := Stack(1)
	if len(stack1) == 0 {
		t.Error("Expected non-empty stack with depth 1")
	}

	// 测试深度为3
	stack3 := Stack(3)
	if len(stack3) == 0 {
		t.Error("Expected non-empty stack with depth 3")
	}

	// 我们不强制要求不同深度的栈长度不同，因为这取决于具体的调用栈
	// 只要确保都能返回非空结果即可
	t.Logf("Stack with depth 1 length: %d", len(stack1))
	t.Logf("Stack with depth 3 length: %d", len(stack3))
}

// 测试Logger接口的实现
func TestDefaultLogger(t *testing.T) {
	logger := &DefaultLogger{}

	// 测试各种日志级别是否能正常输出（不会panic）
	logger.Debug("Debug message")
	logger.Debugf("Debug message with format %s", "test")
	logger.Info("Info message")
	logger.Infof("Info message with format %s", "test")
	logger.Warn("Warn message")
	logger.Warnf("Warn message with format %s", "test")
	logger.Error("Error message")
	logger.Errorf("Error message with format %s", "test")

	// 如果没有panic，测试通过
	t.Log("All logger methods executed successfully")
}