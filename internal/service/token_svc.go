package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"ecom-backend/internal/model"
	"encoding/base32"
	"time"
)

type TokenScope string

var (
	ScopeAuthentication TokenScope = "authentication"
)

type TokenService struct {
	db         *sql.DB
	tokenModel *model.TokenModel
	userModel  *model.UserModel
}

func NewTokenService(db *sql.DB, tokenModel *model.TokenModel, userModel *model.UserModel) *TokenService {
	return &TokenService{db: db, tokenModel: tokenModel, userModel: userModel}
}

func (svc *TokenService) GenerateHash(tokenPlaintext string) []byte {
	hash := sha256.Sum256([]byte(tokenPlaintext))
	return hash[:]
}

func (svc *TokenService) NewToken(ctx context.Context, userId string, lifeDuration time.Duration, scope TokenScope) (*TokenDTO, error) {
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)

	if err != nil {
		return nil, err
	}

	tokenPaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := svc.GenerateHash(tokenPaintext)

	tokenRecord := model.TokenRecord{
		UserId: userId,
		Expiry: time.Now().Add(lifeDuration),
		Scope:  string(scope),
		Hash:   hash[:],
	}

	err = svc.tokenModel.Insert(ctx, svc.db, &tokenRecord)

	if err != nil {
		return nil, err
	}

	tokenDto := TokenDTO{
		Plaintext: tokenPaintext,
		Hash:      tokenRecord.Hash,
		UserId:    tokenRecord.UserId,
		Expiry:    tokenRecord.Expiry,
		Scope:     TokenScope(tokenRecord.Scope),
	}

	return &tokenDto, err
}

func (svc *TokenService) GetUserByToken(ctx context.Context, tokenPlaintext string, scope TokenScope) (*model.UserRecord, error) {
	hash := svc.GenerateHash(tokenPlaintext)

	return svc.tokenModel.GetUserByToken(ctx, svc.db, hash, string(scope), time.Now())
}

type TokenDTO struct {
	Plaintext string     `json:"token"`
	Hash      []byte     `json:"-"`
	UserId    string     `json:"-"`
	Expiry    time.Time  `json:"expiry"`
	Scope     TokenScope `json:"-"`
}
