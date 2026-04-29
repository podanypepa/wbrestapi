// Package main ...
package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"
	"github.com/podanypepa/wbrestapi/internal/adapter/handler"
	"github.com/podanypepa/wbrestapi/internal/adapter/repository"
	"github.com/podanypepa/wbrestapi/internal/application/usecase"
	"github.com/podanypepa/wbrestapi/internal/config"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS


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

	db, err := gorm.Open(gormpg.Open(dsn), &gorm.Config{
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

	// Run database migrations
	if err := runMigrations(sqlDB); err != nil {
		log.Fatal("failed to run database migrations:", err)
	}

	slog.Info("database connection established and migrated")

	// Initialize validator
	v := validator.New()

	// Use JSON tag names for validation errors
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

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
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust this in production to specific domains
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Prometheus metrics
	prometheus := fiberprometheus.New("wbrestapi")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

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

	// Swagger UI
	app.Static("/openapi.yaml", "./api/openapi.yaml")
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/openapi.yaml",
	}))

	// Register routes
	h := &handler.UserHandler{
		SaveUC:    saveUC,
		GetUC:     getUC,
		Logger:    logger,
		Validator: v,
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

func runMigrations(db *sql.DB) error {
	d, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration target: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("could not create migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", d)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	return nil
}
