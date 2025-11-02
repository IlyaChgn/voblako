package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	mygrpc "github.com/IlyaChgn/voblako/internal/pkg/auth/delivery/grpc"
	authproto "github.com/IlyaChgn/voblako/internal/pkg/auth/delivery/grpc/protobuf"
	"github.com/IlyaChgn/voblako/internal/pkg/auth/repository"

	"github.com/joho/godotenv"

	"github.com/IlyaChgn/voblako/internal/pkg/config"
	"github.com/IlyaChgn/voblako/internal/pkg/server/dbinit"

	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading env file:", err)
	}

	cfgPath := os.Getenv("CONFIG_PATH")
	cfg := config.ReadConfig(cfgPath)
	if cfg == nil {
		log.Fatalf("Something went wrong while opening config in auth service")
	}

	postgresURL := dbinit.NewConnectionString(cfg.Postgres.Username, cfg.Postgres.Password,
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)
	postgresPool, err := dbinit.NewPostgresPool(postgresURL)
	if err != nil {
		log.Fatal("Something went wrong while creating postgres pool", err)
	}

	if err = postgresPool.Ping(context.Background()); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	err = postgresPool.Ping(context.Background())
	if err != nil {
		log.Fatal("Cannot ping postgres database", err)
	}

	redisClient := dbinit.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("Cannot ping Redis", err)
	}

	authStorage := repository.NewAuthStorage(postgresPool)
	sessionManager := repository.NewSessionManager(redisClient)
	authManager := mygrpc.NewAuthManager(sessionManager, authStorage)

	grpcAddr := fmt.Sprintf("%s:%s", cfg.Auth.InternalHost, cfg.Auth.Port)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Error occurred while listening gRPC service on %s: %v", grpcAddr, err)
	}

	srv := grpc.NewServer()
	authproto.RegisterAuthServer(srv, authManager)

	log.Printf("Starting Auth gRPC service on %s", grpcAddr)

	if err = srv.Serve(listener); err != nil {
		log.Fatalf("gRPC Server failed to serve: %v", err)
	}
}
