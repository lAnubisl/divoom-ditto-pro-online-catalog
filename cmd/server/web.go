package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"divoom-gallery-client/internal/handlers"
	"divoom-gallery-client/internal/rotation"
)

func serveWeb(ctx context.Context, address string, ui *handlers.Server, runner *rotation.Runner) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	autoCtx, cancelAuto := context.WithCancel(ctx)
	defer cancelAuto()
	autoDone, err := runner.Start(autoCtx)
	if err != nil {
		listener.Close()
		return err
	}
	srv := &http.Server{Handler: ui.Handler(), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second,
		WriteTimeout: 520 * time.Second, MaxHeaderBytes: 16 * 1024, BaseContext: func(net.Listener) context.Context { return ctx }}
	log.Print("Starting web interface and background animation schedule; keep one controller per account, device and media directory.")
	host, port, _ := net.SplitHostPort(listener.Addr().String())
	if host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	log.Printf("Open Categories: http://%s/categories", net.JoinHostPort(host, port))
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if srv.Shutdown(shutdown) != nil {
				srv.Close()
			}
		case <-done:
		}
	}()
	err = srv.Serve(listener)
	close(done)
	cancelAuto()
	<-autoDone
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
