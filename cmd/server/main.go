// Package main ...
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/podanypepa/wbrestapi/internal/adapter/handler"
	"github.com/podanypepa/wbrestapi/internal/adapter/repository"
	"github.com/podanypepa/wbrestapi/internal/application/usecase"
	"github.com/podanypepa/wbrestapi/internal/config"
	"github.com/podanypepa/wbrestapi/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found. Falling back to system environment variables.")
	}

	// Load configuration
	cfg := config.Load()

	// Setup structured logging
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: nil, // Disable GORM's default logger
	})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get database instance:", err)
	}
	configureDatabasePool(sqlDB, cfg)

	// Auto-migrate schema
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	slog.Info("database connection established and migrated")

	// Initialize repositories and use cases
	repo := &repository.UserGormRepository{DB: db}
	saveUC := &usecase.SaveUserUseCase{Repo: repo}
	getUC := &usecase.GetUserUseCase{Repo: repo}

	// Setup Fiber app with middleware
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
		AppName:      "wbrestapi",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// Rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.Server.RateLimitMax,
		Expiration: cfg.Server.RateLimitWindow,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "too many requests",
			})
		},
	}))

	// Register routes
	h := &handler.UserHandler{
		SaveUC: saveUC,
		GetUC:  getUC,
		Logger: logger,
	}
	h.RegisterRoutes(app)

	// Start server in goroutine
	go func() {
		addr := ":" + cfg.Server.Port
		slog.Info("starting server", "address", addr)
		if err := app.Listen(addr); err != nil {
			log.Printf("Server error: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server gracefully...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Close database connection
	if err := sqlDB.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	slog.Info("server stopped")
}

func configureDatabasePool(db *sql.DB, cfg *config.Config) {
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "internal server error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	slog.Error("request error",
		"path", c.Path(),
		"method", c.Method(),
		"error", err.Error(),
		"status", code,
	)

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}
