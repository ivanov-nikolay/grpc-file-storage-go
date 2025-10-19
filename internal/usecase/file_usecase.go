package usecase

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/grpc-file-storage-go/internal/domain"
	"github.com/grpc-file-storage-go/internal/repository"

	"github.com/google/uuid"
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
		Path:      filePath,
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
	file, err := uc.repo.GetByFileName(ctx, filename)
	if err != nil {
		return nil, nil, err
	}

	fileReader, err := os.Open(file.Path)
	if err != nil {
		return nil, nil, err
	}

	return file, fileReader, nil
}

func (uc *fileUseCase) ListFiles(ctx context.Context, page, pageSize int) (*domain.FileList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return uc.repo.List(ctx, page, pageSize)
}
