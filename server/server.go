//go:build !windows

package server

import (
	"net"
	"net/http"
	"os"
)

func (s *ser) Serve() {
	sockPath := s.sockPath

	if sockPath == "" {
		panic("sockPath is nil")
	}

	s.logger.Info("saturn server Unix Serve ...")
	if err := os.Remove(sockPath); err != nil && !os.IsNotExist(err) {
		s.logger.Warnf("Failed to remove existing socket file: %v", err)
	}
	server := http.Server{
		Handler: s,
	}
	unixListener, err := net.Listen("unix", sockPath)
	if err != nil {
		panic(err)
	}
	if err := server.Serve(unixListener); err != nil {
		panic(err)
	}
}
