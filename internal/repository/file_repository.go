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
	query := `INSERT INTO files (id, filename, size, path, created_at, updated_at) 
				VALUES($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		file.ID,
		file.Filename,
		file.Size,
		file.Path,
		file.CreatedAt,
		file.UpdatedAt,
	)

	return err
}
func (r *postgresFileRepository) GetByFileName(ctx context.Context, fileName string) (*domain.File, error) {
	query := `SELECT id, filename, size, path, created_at, updated_at FROM files WHERE filename = $1`

	file := &domain.File{}
	err := r.db.QueryRowContext(ctx, query, fileName).Scan(
		&file.ID,
		&file.Filename,
		&file.Size,
		&file.Path,
		&file.CreatedAt,
		&file.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return file, nil
}
func (r *postgresFileRepository) List(ctx context.Context, page, pageSize int) (*domain.FileList, error) {
	offset := (page - 1) * pageSize

	var total int
	countQuery := `SELECT COUNT(*) FROM files`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, err
	}

	query := `
				SELECT id, filename, size, path, created_at, updated_at 
				FROM files
				ORDER BY created_at DESC
				LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]domain.File, 0)
	for rows.Next() {
		var file domain.File
		err := rows.Scan(
			&file.ID,
			&file.Filename,
			&file.Size,
			&file.Path,
			&file.CreatedAt,
			&file.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return &domain.FileList{
		Files: files,
		Total: total,
	}, nil
}
