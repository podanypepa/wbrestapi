// Package main ...
package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const (
	idleTimeout = 5 * time.Second
	defaultPort = "3000"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found. Falling back to system environment variables.")
	}

	if err := initDB(); err != nil {
		slog.Error("initDB", "err", err)
		os.Exit(1)
	}

	app := apiSetup()

	port := cmp.Or(os.Getenv("PORT"), defaultPort)
	go func() {
		if err := app.Listen(":" + port); err != nil {
			slog.Error("app.Listen", "err", err, "port", port)
			os.Exit(2)
		}
	}()

	slog.Info(fmt.Sprintf("app listen on http://127.0.0.1:%s", port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Gracefully shutting down...")

	if err := app.Shutdown(); err != nil {
		slog.Error("app.Shutdown", "err", err)
	}

	slog.Info("Server exiting")
}
