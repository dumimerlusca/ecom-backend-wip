package model

import (
	"context"
	"ecom-backend/pkg/sqldb"
	"time"
)

// reusable junction table for each entity that can own files ( ex: products, users, etc  )
type EntityFileRecord struct {
	EntityId  string // product_id, user_id, etc
	FileId    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EntityFileModel struct{}

func NewEntityFileModel() *EntityFileModel {
	return &EntityFileModel{}
}

func (e *EntityFileModel) Insert(ctx context.Context, conn sqldb.Connection, entityFile *EntityFileRecord) (*EntityFileRecord, error) {
	q := `INSERT INTO entity_file (entity_id, file_id) VALUES ($1, $2) RETURNING created_at, updated_at`

	err := conn.QueryRowContext(ctx, q, entityFile.EntityId, entityFile.FileId).Scan(&entityFile.CreatedAt, &entityFile.UpdatedAt)

	if err != nil {
		msg := err.Error()

		switch {
		case msg == `pq: insert or update on table "entity_file" violates foreign key constraint "entity_file_file_id_fkey"`:
			return nil, ErrFileNotFound
		default:
			return nil, err

		}

	}

	return entityFile, nil
}
