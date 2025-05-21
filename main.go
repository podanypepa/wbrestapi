// Package main ...
package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/podanypepa/wbrestapi/pkg/api"
	"github.com/podanypepa/wbrestapi/pkg/repository"
)

const (
	defaultPort = "3000"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found. Falling back to system environment variables.")
	}

	db, err := initDB()
	if err != nil {
		slog.Error("initDB", "err", err)
		os.Exit(1)
	}

	userRepository, err := repository.NewUserRepository(repository.UserRepositoryConfig{
		Db: db,
	})
	if err != nil {
		slog.Error("NewUserRepository", "err", err)
		os.Exit(2)
	}

	apiServer := api.NewServer(api.Config{
		UserRepository: userRepository,
	})

	port := cmp.Or(os.Getenv("PORT"), defaultPort)
	go func() {
		if err := apiServer.Listen(":" + port); err != nil {
			slog.Error("app.Listen", "err", err, "port", port)
			os.Exit(2)
		}
	}()

	slog.Info(fmt.Sprintf("app listen on http://127.0.0.1:%s", port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Gracefully shutting down...")

	if err := apiServer.Shutdown(); err != nil {
		slog.Error("app.Shutdown", "err", err)
	}

	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		slog.Error("sqlDB.Close", "err", err)
	}

	slog.Info("Server exiting")
}
