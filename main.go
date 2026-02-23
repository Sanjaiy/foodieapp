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
	"github.com/Sanjaiy/foodieapp/internal/helpers"
	"github.com/Sanjaiy/foodieapp/internal/service"
	pgstore "github.com/Sanjaiy/foodieapp/internal/store/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	dbConn, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	if err := runMigrations(dbConn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	validCodesPath := os.Getenv("VALID_CODES_PATH")
	if validCodesPath == "" {
		validCodesPath = "data/valid_codes.txt"
	}

	promoLookup, err := helpers.NewCouponLookup(validCodesPath)
	if err != nil {
		log.Fatalf("Failed to load promo codes: %v", err)
	}
	defer promoLookup.Close()

	rootHandler := setupApp(dbConn, cfg.APIKey, promoLookup)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      rootHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	runServer(srv)
}

func setupApp(dbConn *sql.DB, apiKey string, promoLookup *helpers.CouponLookup) http.Handler {
	promoSvc := service.NewPromoService(promoLookup)

	productStore := pgstore.NewProductStore(dbConn)
	orderStore := pgstore.NewOrderStore(dbConn)

	productSvc := service.NewProductService(productStore)
	orderSvc := service.NewOrderService(orderStore, promoSvc)

	productHandler := handler.NewProductHandler(productSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/product", productHandler.ListProducts)
	mux.HandleFunc("GET /api/product/{productId}", productHandler.GetProduct)

	orderMux := http.NewServeMux()
	orderMux.HandleFunc("POST /api/order", orderHandler.PlaceOrder)
	mux.Handle("/api/order", handler.AuthMiddleware(apiKey, orderMux))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	var root http.Handler = mux
	root = handler.CORSMiddleware(root)
	root = handler.LoggingMiddleware(root)

	return root
}

func runServer(srv *http.Server) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on %s", srv.Addr)
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
