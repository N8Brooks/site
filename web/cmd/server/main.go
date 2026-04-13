package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	siteassets "github.com/N8Brooks/site/web"
	"github.com/N8Brooks/site/web/internal/server"
)

const httpAddr = ":8080"

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run site server: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	processCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	handler, err := server.New(siteassets.Dist)
	if err != nil {
		return fmt.Errorf("create handler: %w", err)
	}

	httpServer := &http.Server{
		Addr:              httpAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		<-processCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(processCtx, "shutdown site server failed", slog.Any("error", err))
		}
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen and serve: %w", err)
	}
	return nil
}
