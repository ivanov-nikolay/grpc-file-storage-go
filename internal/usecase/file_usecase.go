package usecase

import (
	"context"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/grpc-file-storage-go/internal/domain"
	"github.com/grpc-file-storage-go/internal/repository"
)

type fileUseCase struct {
	repo        repository.FileRepository
	storagePath string
}

func NewFileUseCase(repo repository.FileRepository, storagePath string) FileUseCase {
	return &fileUseCase{
		repo:        repo,
		storagePath: storagePath,
	}
}

func (uc *fileUseCase) UploadFile(ctx context.Context, filename string, data io.Reader) (*domain.File, error) {
	ext := filepath.Ext(filename)
	baseName := filename[:len(filename)-len(ext)]
	uniqueFilename := baseName + "_" + uuid.New().String() + ext
	filePath := filepath.Join(uc.storagePath, uniqueFilename)

	if err := os.MkdirAll(uc.storagePath, 0755); err != nil {
		return nil, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	size, err := io.Copy(file, data)
	if err != nil {
		os.Remove(filePath)
		return nil, err
	}

	fileMetadata := &domain.File{
		ID:        uuid.New().String(),
		Filename:  uniqueFilename,
		Size:      size,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.repo.Save(ctx, fileMetadata); err != nil {
		os.Remove(filePath)
		return nil, err
	}

	return fileMetadata, nil
}

func (uc *fileUseCase) DownLoadFile(ctx context.Context, filename string) (*domain.File, io.Reader, error) {
	return nil, nil, nil
}

func (uc *fileUseCase) ListFiles(ctx context.Context, page, pageSize int) (*domain.FileList, error) {
	return nil, nil
}
