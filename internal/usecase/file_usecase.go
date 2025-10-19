package usecase

import (
	"context"
	"io"

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
	return nil, nil
}

func (uc *fileUseCase) DownLoadFile(ctx context.Context, filename string) (*domain.File, io.Reader, error) {
	return nil, nil, nil
}
func (uc *fileUseCase) ListFiles(ctx context.Context, page, pageSize int) (*domain.FileList, error) {
	return nil, nil
}
