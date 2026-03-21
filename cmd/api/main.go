package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go_pg_http/internal/bootstrap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	initCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	app, err := bootstrap.BuildApp(initCtx)
	if err != nil {
		log.Fatalf("build app: %v", err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			log.Printf("close app: %v", err)
		}
	}()

	addr := fmt.Sprintf("%s:%d", app.Config.HTTP.Host, app.Config.HTTP.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: app.Router,
	}

	go func() {
		log.Printf(
			"http server started: addr=%s app=%s env=%s",
			addr,
			app.Config.App.Name,
			app.Config.App.Env,
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Printf("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutdown: %v", err)
	}
}
