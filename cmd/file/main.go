package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	mygrpc "github.com/IlyaChgn/voblako/internal/pkg/file/delivery/grpc"
	fileproto "github.com/IlyaChgn/voblako/internal/pkg/file/delivery/grpc/protobuf"
	metarepo "github.com/IlyaChgn/voblako/internal/pkg/file/repository/metadata"
	objectrepo "github.com/IlyaChgn/voblako/internal/pkg/file/repository/object"

	"github.com/joho/godotenv"

	"github.com/IlyaChgn/voblako/internal/pkg/config"
	"github.com/IlyaChgn/voblako/internal/pkg/server/dbinit"

	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Println(".env file not found, using OS environment")
	}

	cfgPath := os.Getenv("CONFIG_PATH")
	generalCfg := config.ReadConfig(cfgPath)
	if generalCfg == nil {
		log.Fatalf("Something went wrong while opening config in file service")
	}
	cfg := generalCfg.File

	postgresURL := dbinit.NewConnectionString(cfg.Postgres.Username, cfg.Postgres.Password,
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)
	postgresPool, err := dbinit.NewPostgresPool(postgresURL)
	if err != nil {
		log.Fatal("Something went wrong while creating postgres pool", err)
	}

	err = postgresPool.Ping(context.Background())
	if err != nil {
		log.Fatal("Cannot ping postgres database", err)
	}

	minioURL := dbinit.NewMinioEndpoint(cfg.Minio.Host, cfg.Minio.Port)
	minioClient, err := dbinit.NewMinioClient(minioURL, cfg.Minio.AccessKey, cfg.Minio.SecretKey, cfg.Minio.Bucket)
	if err != nil {
		log.Fatal("Something went wrong while creating minio client", err)
	}

	metadataStorage := metarepo.NewMetadataStorage(postgresPool)
	objectStorage := objectrepo.NewObjectStorage(minioClient, cfg.Minio.Bucket)
	fileManager := mygrpc.NewFileManager(metadataStorage, objectStorage)

	grpcAddr := fmt.Sprintf("%s:%s", cfg.InternalHost, cfg.Port)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Error occurred while listening gRPC service on %s: %v", grpcAddr, err)
	}

	srv := grpc.NewServer()
	fileproto.RegisterFileServer(srv, fileManager)

	log.Printf("Starting File gRPC service on %s", grpcAddr)

	if err = srv.Serve(listener); err != nil {
		log.Fatalf("gRPC Server failed to serve: %v", err)
	}
}
