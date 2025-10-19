package repository

import (
	"context"
	"database/sql"

	"github.com/grpc-file-storage-go/internal/domain"
)

type postgresFileRepository struct {
	db *sql.DB
}

func NewPostgresFileRepository(db *sql.DB) FileRepository {
	return &postgresFileRepository{
		db: db,
	}
}

func (r *postgresFileRepository) Save(ctx context.Context, file *domain.File) error {
	return nil
}
func (r *postgresFileRepository) GetByFileName(ctx context.Context, fileName string) (*domain.File, error) {
	return nil, nil
}
func (r *postgresFileRepository) List(ctx context.Context, page, pageSize int) (*domain.FileList, error) {
	return nil, nil
}
func (r *postgresFileRepository) Delete(ctx context.Context, fileName string) error {
	return nil
}
func (r *postgresFileRepository) Update(ctx context.Context, file *domain.File) error {
	return nil
}
