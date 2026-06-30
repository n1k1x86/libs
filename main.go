package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/n1k1x86/libs/http_server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	hs := http_server.NewMux()
	hs.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello world"))
		}
	})

	s_cfg := http_server.NewHTTPServerConfig().
		WithAddr("localhost:9999").
		WithReadTimeout(time.Second * 5).
		WithWriteTimeout(time.Second * 5).
		WithIdleTimeout(time.Second * 5)

	s := http_server.NewHTTPServer(s_cfg).WithMux(hs)

	errChan := make(chan error, 1)

	go func() {
		errChan <- s.Start()
	}()

	select {
	case <-ctx.Done():
		log.Println("shut down by signal")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer shutdownCancel()

		err := s.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
	case err := <-errChan:
		log.Println("shut down by err from chan")
		if err != nil {
			log.Fatal(err)
		}
	}
}
