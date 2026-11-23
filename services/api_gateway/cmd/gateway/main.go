package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/app"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/config"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/database"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/httpapi"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("gateway stopped", "reason", err.Error())
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if command == "openapi" {
		return json.NewEncoder(os.Stdout).Encode(httpapi.OpenAPI())
	}
	if command != "serve" && command != "migrate" && command != "healthcheck" {
		return fmt.Errorf("usage: gateway [serve|migrate|healthcheck|openapi]")
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if command == "healthcheck" {
		_, port, _ := net.SplitHostPort(cfg.HTTPAddr)
		client := &http.Client{Timeout: 3 * time.Second}
		response, err := client.Get("http://127.0.0.1:" + port + "/readyz")
		if err != nil {
			return fmt.Errorf("HTTP readiness check failed")
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("service is not ready")
		}
		return nil
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	if command == "migrate" {
		migrationCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()
		if err := database.Migrate(migrationCtx, db); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
		logger.Info("database migrations applied")
		return nil
	}
	readyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	err = database.Ready(readyCtx, db)
	cancel()
	if err != nil {
		return fmt.Errorf("database not ready; run migrations first: %w", err)
	}
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: app.Handler(db, func(c context.Context) error { return database.Ready(c, db) }, logger), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	fail := make(chan error, 1)
	go func() { logger.Info("HTTP server starting", "address", cfg.HTTPAddr); fail <- server.ListenAndServe() }()
	select {
	case err := <-fail:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdown); err != nil {
			return err
		}
	}
	return nil
}
