package usecase

import (
	"context"
	"io"

	"github.com/grpc-file-storage-go/internal/domain"
)

type FileUseCase interface {
	UploadFile(ctx context.Context, filename string, data io.Reader) (*domain.File, error)
	DownLoadFile(ctx context.Context, filename string) (*domain.File, io.Reader, error)
	ListFiles(ctx context.Context, page, pageSize int) (*domain.FileList, error)
}
