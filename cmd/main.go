package main

import (
	"auth_test/configs"
	"auth_test/internal/handler"
	"auth_test/internal/service"
	"auth_test/internal/store"
	"auth_test/pkg/pb"
	"context"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbStore, err := store.NewPostgresStore(cfg.DBConnectionString())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbStore.Close()

	// userStore := store.NewInMemoryStore()
	userService := service.NewUserService(dbStore, cfg.JWTSecret)
	grpcHandler := handler.NewGRPCHandler(userService)

	ctx := context.Background()
	userService.CreateUser(ctx, "admin", "admin123")
	userService.CreateUser(ctx, "user", "user123")

	go startMetricsServer(cfg.MetricsPort)
	startGRPCServer(grpcHandler, cfg.GRPCPort, true)
}

func startGRPCServer(grpcHandler *handler.GRPCHandler, port string, enableReflection bool) {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(handler.MetricsInterceptor()))
	pb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	// для проверки reflection api
	if enableReflection {
		reflection.Register(grpcServer)
		log.Printf("gRPC reflection enabled")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	log.Printf("gRPC server starting on %s", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startMetricsServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Metrics server starting on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
