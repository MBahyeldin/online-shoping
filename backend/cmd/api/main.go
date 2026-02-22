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

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/online-cake-shop/backend/internal/config"
	"github.com/online-cake-shop/backend/internal/email"
	"github.com/online-cake-shop/backend/internal/handler"
	custmw "github.com/online-cake-shop/backend/internal/middleware"
	"github.com/online-cake-shop/backend/internal/repository/db"
	"github.com/online-cake-shop/backend/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found, using environment variables")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Database connection pool
	pool, err := pgxpool.New(context.Background(), cfg.Database.URL())
	if err != nil {
		logger.Error("failed to create db pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		logger.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	logger.Info("database connected", "host", cfg.Database.Host, "dbname", cfg.Database.Name)

	// Run migrations
	m, err := migrate.New("file://db/migrations", cfg.Database.URL())
	if err != nil {
		logger.Error("failed to create migrator", "error", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}
	logger.Info("database migrations applied")

	// Dependency graph
	queries := db.New(pool)

	var emailSender email.Sender
	if cfg.Email.Provider == "smtp" {
		emailSender = email.NewSMTPSender(cfg.Email)
	} else {
		emailSender = email.NewMockSender(logger)
	}

	authSvc := service.NewAuthService(queries, emailSender, cfg.JWT)
	productSvc := service.NewProductService(queries)
	cartSvc := service.NewCartService(queries)
	orderSvc := service.NewOrderService(pool, queries)

	authHandler := handler.NewAuthHandler(authSvc)
	productHandler := handler.NewProductHandler(productSvc)
	cartHandler := handler.NewCartHandler(cartSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)

	authMiddleware := custmw.NewAuthMiddleware(cfg.JWT.Secret)

	// Router
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(custmw.Logger(logger))
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Server.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Auth (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/verify-otp", authHandler.VerifyOTP)
			r.Post("/resend-otp", authHandler.ResendOTP)
		})

		// Products (public)
		r.Route("/products", func(r chi.Router) {
			r.Get("/", productHandler.List)
			r.Get("/{id}", productHandler.GetByID)
		})
		r.Get("/categories", productHandler.ListCategories)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Route("/cart", func(r chi.Router) {
				r.Get("/", cartHandler.GetCart)
				r.Post("/items", cartHandler.AddItem)
				r.Put("/items/{itemId}", cartHandler.UpdateItem)
				r.Delete("/items/{itemId}", cartHandler.RemoveItem)
				r.Delete("/", cartHandler.ClearCart)
			})

			r.Route("/orders", func(r chi.Router) {
				r.Post("/", orderHandler.CreateOrder)
				r.Get("/", orderHandler.ListOrders)
				r.Get("/{id}", orderHandler.GetOrder)
			})
		})
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("server starting", "port", cfg.Server.Port, "env", cfg.Server.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("server stopped")
}
