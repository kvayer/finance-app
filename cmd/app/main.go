package main

import (
	"context"
	"finance-tracker/internal/config"
	"finance-tracker/internal/db"
	"finance-tracker/internal/middleware"
	"finance-tracker/internal/repository"
	"finance-tracker/internal/server"
	"finance-tracker/internal/service"
	"finance-tracker/internal/session"
	"finance-tracker/internal/transport"
	"finance-tracker/pkg/hash"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	// 2. DB
	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("db init error: %s", err)
	}
	defer database.Close()

	// 3. Repositories
	userRepo := repository.NewUserRepo(database)
	sessionRepo := repository.NewSessionRepo(database)
	expenseRepo := repository.NewExpenseRepo(database)

	// 4. Services & Helpers
	hasher := hash.NewSHA256Hasher()
	authService := service.NewAuthService(userRepo, hasher)
	expenseService := service.NewExpenseService(expenseRepo)
	sessionManager := session.NewManager(sessionRepo)

	// 5. Middleware
	authMiddleware := middleware.NewAuthMiddleware(sessionManager)

	// 6. Handlers
	handler := transport.NewHandler(authService, expenseService, sessionManager, authMiddleware)

	// 7. Server
	srv := server.NewServer(cfg, handler.InitRoutes())

	// Run
	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("server run error: %s", err)
		}
	}()

	log.Printf("App running on port %s", cfg.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("shutdown error: %s", err)
	}
}
