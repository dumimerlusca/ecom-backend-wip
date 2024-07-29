package model

import (
	"context"
	"database/sql"
	"ecom-backend/pkg/sqldb"
	"errors"
	"time"
)

type WishlistRecord struct {
	Id             string    `json:"id"`
	UserIdentifier string    `json:"user_identifier"` // user id for registered users and session id for guests
	VariantId      string    `json:"variant_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type WishlistModel struct{}

func NewWishlistModel() *WishlistModel {
	return &WishlistModel{}
}

func (m *WishlistModel) Insert(ctx context.Context, conn sqldb.Connection, record *WishlistRecord) (*WishlistRecord, error) {
	q := `INSERT INTO wishlist(user_identifier, variant_id) VALUES($1, $2) RETURNING id, created_at`

	err := conn.QueryRowContext(ctx, q, record.UserIdentifier, &record.VariantId).Scan(&record.Id, &record.CreatedAt)

	if err != nil {

		if err.Error() == `pq: duplicate key value violates unique constraint "duplicate_product_in_favorites_not_allowed"` {
			return nil, ErrProductAlreadyWishlisted
		}
		return nil, err
	}

	return record, nil
}

func (m *WishlistModel) FindById(ctx context.Context, conn sqldb.Connection, id string) (*WishlistRecord, error) {
	q := `SELECT id, user_identifier, variant_id, created_at FROM wishlist WHERE id = $1`

	var record WishlistRecord

	err := conn.QueryRowContext(ctx, q, id).Scan(&record.Id, &record.UserIdentifier, &record.VariantId, &record.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &record, nil
}

func (m *WishlistModel) DeleteById(ctx context.Context, conn sqldb.Connection, id string) error {
	q := `DELETE FROM wishlist WHERE id = $1`

	res, err := conn.ExecContext(ctx, q, id)

	if err != nil {
		return err
	}

	if rowsAff, _ := res.RowsAffected(); rowsAff == 0 {
		return ErrRecordNotFound
	}

	return nil
}
