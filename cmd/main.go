package main

import (
	"auth_test/configs"
	"auth_test/internal/handler"
	"auth_test/internal/service"
	"auth_test/internal/store"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// sugar := logger.Sugar()

	cfg := configs.LoadConfig()

	userStore := store.NewInMemoryStore()
	userService := service.NewUserService(userStore, cfg.JWTSecret)
	loginHandler := handler.NewLoginHandler(userService)
	verifyHandler := handler.NewVerifyHandler(userService)

	http.HandleFunc("/login", loginHandler.Handle)
	http.HandleFunc("/verify", verifyHandler.Handle)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))

	// логи
	// sugar.Info("Auth service starting...")
	// sugar.Infof("Server port: %s", cfg.Port)

	// sugar.Info("Service stopped")

}
