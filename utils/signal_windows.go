//go:build windows
// +build windows

package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func ListenSignal() chan os.Signal {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	return signalChan
}

// StopSignal cancels signal notifications and closes the provided channel.
func StopSignal(signalChan chan os.Signal) {
	if signalChan == nil {
		return
	}
	signal.Stop(signalChan)
	close(signalChan)
}
