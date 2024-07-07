package model

import (
	"context"
	"database/sql"
	"ecom-backend/pkg/sqldb"
	"errors"
	"time"
)

type FileRecord struct {
	Id           string
	OriginalName string
	MimeType     string
	Extension    string
	Size         int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FileModel struct{}

func NewFileModel() *FileModel {
	return &FileModel{}
}

func (f *FileModel) Insert(ctx context.Context, conn sqldb.Connection, file *FileRecord) (*FileRecord, error) {
	q := `INSERT INTO file (original_name, mime_type,extension, size) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := conn.QueryRowContext(ctx, q, file.OriginalName, file.MimeType, file.Extension, file.Size).Scan(&file.Id, &file.CreatedAt, &file.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *FileModel) FindById(ctx context.Context, conn sqldb.Connection, id string) (*FileRecord, error) {
	q := `SELECT id, original_name, mime_type, extension, size, created_at, updated_at FROM file WHERE id = $1`

	file := &FileRecord{}

	err := conn.QueryRowContext(ctx, q, id).Scan(&file.Id, &file.OriginalName, &file.MimeType, &file.Extension, &file.Size, &file.CreatedAt, &file.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return file, nil
}
