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
	return signalChan
}
