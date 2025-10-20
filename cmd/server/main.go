package main

import (
	"log"
	"net"

	"github.com/grpc-file-storage-go/api/proto"
	"github.com/grpc-file-storage-go/internal/config"
	handlergrpc "github.com/grpc-file-storage-go/internal/handler/grpc"
	"github.com/grpc-file-storage-go/internal/repository"
	"github.com/grpc-file-storage-go/internal/usecase"
	"github.com/grpc-file-storage-go/pkg/database"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	migrationManager := database.NewMigrationManager(db)

	if err := migrationManager.RunMigrations(cfg.MigrationsPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fileRepo := repository.NewPostgresFileRepository(db)

	fileUseCase := usecase.NewFileUseCase(fileRepo, cfg.StoragePath)

	fileHandler := handlergrpc.NewFileHandler(fileUseCase)

	limiter := handlergrpc.NewConcurrencyLimiter(
		cfg.UploadLimit,
		cfg.DownloadLimit,
		cfg.ListLimit,
	)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(limiter.UnaryInterceptor()),
		grpc.StreamInterceptor(limiter.StreamInterceptor()),
	)

	proto.RegisterFileServiceServer(server, fileHandler)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("serving gRPC starting on port %s", cfg.GRPCPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
