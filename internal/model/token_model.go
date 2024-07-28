package model

import (
	"context"
	"database/sql"
	"ecom-backend/pkg/sqldb"
	"errors"
	"time"
)

type TokenRecord struct {
	Hash   []byte
	UserId string
	Expiry time.Time
	Scope  string
}

type TokenModel struct {
}

func NewTokenModel() *TokenModel {
	return &TokenModel{}
}

func (m *TokenModel) Insert(ctx context.Context, conn sqldb.Connection, record *TokenRecord) error {
	q := `INSERT INTO token(hash, user_id, expiry, scope) VALUES($1, $2, $3, $4)`

	_, err := conn.ExecContext(ctx, q, record.Hash, record.UserId, record.Expiry, record.Scope)

	return err
}

func (m *TokenModel) GetUserByToken(ctx context.Context, conn sqldb.Connection, hash []byte, scope string, expiry time.Time) (*UserRecord, error) {
	var user UserRecord

	q := `SELECT u.id, u.name, u.email, u.is_admin, u.activated, u.created_at, u.updated_at FROM users as u
	INNER JOIN token ON token.user_id = u.id
	WHERE token.hash = $1 AND token.scope = $2 AND token.expiry > $3`

	err := conn.QueryRowContext(ctx, q, hash, scope, expiry).Scan(&user.Id, &user.Name, &user.Email, &user.IsAdmin, &user.Activated, &user.CreatedAt, &user.UpdatedAt)

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
