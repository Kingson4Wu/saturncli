//go:build windows

package server

import (
	"fmt"
	"net/http"
)

func (s *ser) Serve() {
	s.logger.Info("saturn server Http Serve ...")
	server := http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%s", "8096"),
		Handler: s,
	}
	server.ListenAndServe()

}
