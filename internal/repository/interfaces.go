package repository

import (
	"context"

	"github.com/grpc-file-storage-go/internal/domain"
)

type FileRepository interface {
	Save(ctx context.Context, file *domain.File) error
	GetByFileName(ctx context.Context, fileName string) (*domain.File, error)
	List(ctx context.Context, page, pageSize int) (*domain.FileList, error)
	Delete(ctx context.Context, fileName string) error
	Update(ctx context.Context, file *domain.File) error
}
