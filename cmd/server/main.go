// Package main ...
package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/podanypepa/wbrestapi/internal/adapter/handler"
	"github.com/podanypepa/wbrestapi/internal/adapter/repository"
	"github.com/podanypepa/wbrestapi/internal/application/usecase"
	"github.com/podanypepa/wbrestapi/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultPort = "3000"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found. Falling back to system environment variables.")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatal(err)
	}

	repo := &repository.UserGormRepository{DB: db}

	saveUC := &usecase.SaveUserUseCase{Repo: repo}
	getUC := &usecase.GetUserUseCase{Repo: repo}

	app := fiber.New()

	h := &handler.UserHandler{SaveUC: saveUC, GetUC: getUC}
	h.RegisterRoutes(app)

	port := cmp.Or(os.Getenv("PORT"), defaultPort)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Shutting down server: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Gracefully shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Panicf("Server shutdown failed: %v", err)
	}
}
