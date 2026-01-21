package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/davidgarcia/digital-wallet/internal/config"
	"github.com/davidgarcia/digital-wallet/internal/handler"
	"github.com/davidgarcia/digital-wallet/internal/middleware"
	"github.com/davidgarcia/digital-wallet/internal/repository"
	"github.com/davidgarcia/digital-wallet/internal/usecase"
	"github.com/davidgarcia/digital-wallet/pkg/database"
)

func main() {
	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewPostgresUserRepository(db)
	transactionRepo := repository.NewPostgresTransactionRepository(db)
	walletUsecase := usecase.NewWalletUsecase(db, userRepo, transactionRepo)
	walletHandler := handler.NewWalletHandler(walletUsecase)

	router := mux.NewRouter()
	walletHandler.RegisterRoutes(router)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}).Methods(http.MethodGet)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.AppPort),
		Handler:      middleware.Logging(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on http://localhost:%d", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
}
