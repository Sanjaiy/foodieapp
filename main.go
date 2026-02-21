package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pressly/goose/v3"

	"github.com/Sanjaiy/foodieapp/db"
	"github.com/Sanjaiy/foodieapp/internal/config"
	"github.com/Sanjaiy/foodieapp/internal/database"
	"github.com/Sanjaiy/foodieapp/internal/handler"
	"github.com/Sanjaiy/foodieapp/internal/service"
	pgstore "github.com/Sanjaiy/foodieapp/internal/store/postgres"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	ctx := context.Background()
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	promoSvc := service.NewPromoService()

	// Initialize stores (DI — inject DB connection)
	productStore := pgstore.NewProductStore(db)
	orderStore := pgstore.NewOrderStore(db)

	productSvc := service.NewProductService(productStore)
	orderSvc := service.NewOrderService(orderStore, promoSvc)

	productHandler := handler.NewProductHandler(productSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// Setup routes
	mux := http.NewServeMux()

	// Product routes (public)
	mux.HandleFunc("GET /api/product", productHandler.ListProducts)
	mux.HandleFunc("GET /api/product/{productId}", productHandler.GetProduct)

	// Order routes (authenticated)
	orderMux := http.NewServeMux()
	orderMux.HandleFunc("POST /api/order", orderHandler.PlaceOrder)
	mux.Handle("/api/order", handler.AuthMiddleware(cfg.APIKey, orderMux))

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Apply global middleware
	var root http.Handler = mux
	root = handler.CORSMiddleware(root)
	root = handler.LoggingMiddleware(root)

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      root,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-done
	log.Println("Server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped gracefully")
}

func runMigrations(dbConn *sql.DB) error {
	goose.SetBaseFS(db.MigrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(dbConn, "migrations")
}
