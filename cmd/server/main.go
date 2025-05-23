// Package main ...
package main

import (
	"cmp"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
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
		panic("failed to connect to database")
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
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
