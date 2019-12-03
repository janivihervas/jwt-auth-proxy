// Package server implements a HTTP server
package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"time"
)

// RunHTTP server
func RunHTTP(port string, handler http.Handler, logger *log.Logger) error {
	if port == "" {
		return errors.New("port is not defined")
	}

	timeout := time.Second * 15
	server := &http.Server{
		Addr:           "127.0.0.1:" + port,
		Handler:        handler,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		MaxHeaderBytes: 1 << 20,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-quit
		logger.Printf("Received signal %s, shutting down server\n", s)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		errShutdown := server.Shutdown(ctx)
		if errShutdown == nil {
			logger.Println("Server shutdown successful")
			return
		}

		logger.Println("Server shutdown failed, forcing shutdown:", errShutdown)
		errShutdown = server.Close()
		if errShutdown != nil {
			panic(fmt.Sprint("Force shutdown of server failed:", errShutdown))
		}
	}()

	logger.Println("Starting HTTP server on address", server.Addr)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed && err != nil {
		return err
	}

	return nil
}
