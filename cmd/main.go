package main

import (
	"auth_test/configs"
	"auth_test/internal/handler"
	"auth_test/internal/service"
	"auth_test/internal/store"
	"log"
	"net/http"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	userStore := store.NewInMemoryStore()
	userService := service.NewUserService(userStore, cfg.JWTSecret)
	loginHandler := handler.NewLoginHandler(userService)
	verifyHandler := handler.NewVerifyHandler(userService)

	http.HandleFunc("/login", loginHandler.Handle)
	http.HandleFunc("/verify", verifyHandler.Handle)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
