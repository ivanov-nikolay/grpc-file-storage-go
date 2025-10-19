package main

import (
	"log"
	"path/filepath"

	"github.com/grpc-file-storage-go/internal/config"
	"github.com/grpc-file-storage-go/internal/repository"
	"github.com/grpc-file-storage-go/internal/usecase"
	"github.com/grpc-file-storage-go/pkg/database"
	"github.com/grpc-file-storage-go/pkg/utils"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	migrationManager := database.NewMigrationManager(db)
	projectRoot, err := utils.GetProjectRoot()
	if err != nil {
		log.Fatalf("failed to get project root: %v", err)
	}

	migrationsDir := filepath.Join(projectRoot, "migrations")

	if err := migrationManager.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fileRepo := repository.NewPostgresFileRepository(db)
	defer db.Close()

	fileUseCase := usecase.NewFileUseCase(fileRepo, cfg.StoragePath)
	_ = fileUseCase
}
