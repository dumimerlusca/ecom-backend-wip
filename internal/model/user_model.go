package model

import (
	"context"
	"database/sql"
	"ecom-backend/pkg/sqldb"
	"errors"
	"time"
)

type UserRecord struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Activated    bool      `json:"-"`
	IsAdmin      bool      `json:"-"`
}

type UserModel struct{}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func (m *UserModel) Insert(ctx context.Context, conn sqldb.Connection, record *UserRecord) (*UserRecord, error) {
	q := `INSERT INTO users(name, email, password_hash, activated) VALUES($1, $2, $3, $4)
		  RETURNING id,created_at, updated_at, is_admin`

	row := conn.QueryRowContext(ctx, q, record.Name, record.Email, record.PasswordHash, record.Activated)

	err := row.Scan(&record.Id, &record.CreatedAt, &record.UpdatedAt, &record.IsAdmin)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicatedEmail
		default:
			return nil, err
		}
	}

	return record, nil
}

func (m *UserModel) FindByEmail(ctx context.Context, conn sqldb.Connection, email string) (*UserRecord, error) {
	q := `SELECT id, name, email, activated, password_hash, created_at, updated_at FROM users WHERE email = $1`

	row := conn.QueryRowContext(ctx, q, email)

	var user UserRecord

	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Activated, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
