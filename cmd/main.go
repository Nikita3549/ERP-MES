package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"erp-mes/internal/configs"
	"erp-mes/internal/health"
	"erp-mes/pkg/db"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conf := configs.LoadConfig()
	DB, DBErr := db.NewDB(conf)
	if DBErr != nil {
		log.Fatalf("Database error: %v", DBErr)
	}

	router := http.NewServeMux()

	// Handlers
	health.NewHandler(router, DB)

	port := conf.Port

	server := &http.Server{
		Handler:           router,
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: conf.ReadHeaderTimeout,
		ReadTimeout:       conf.ReadTimeout,
		WriteTimeout:      conf.WriteTimeout,
		IdleTimeout:       conf.IdleTimeout,
	}

	log.Printf("Server is listening on port %d\n", port)

	var exitErr error

	errCh := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
	case exitErr = <-errCh:
		log.Printf("Fatal error: %v", exitErr)
	}
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown error: %v", err)
		exitErr = err
		_ = server.Close()
	}

	if err := DB.Close(); err != nil {
		log.Printf("Database error: %v\n", err)
	}

	if exitErr != nil {
		os.Exit(1)
	}
	log.Println("Server stopped")
}
