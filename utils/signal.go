//go:build !windows

package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func ListenSignal() chan os.Signal {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	// terminated : kill -15 [pid]
	// interrupt: kill -2 [pid] OR kill -SIGINT [pid]
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
