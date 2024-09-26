package service

import (
	"context"
	"database/sql"
	"ecom-backend/internal/model"
	"ecom-backend/internal/validator"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db         *sql.DB
	userModel  *model.UserModel
	tokenModel *model.TokenModel
	tokenSvc   *TokenService
}

func NewAuthService(db *sql.DB, userModel *model.UserModel, tokenModel *model.TokenModel, tokenSvc *TokenService) *AuthService {
	return &AuthService{db: db, userModel: userModel, tokenModel: tokenModel, tokenSvc: tokenSvc}
}

type RegisterUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (input *RegisterUserInput) Validate(v *validator.Validator) {
	v.Check(input.Name != "", "name", "must not be empty")

	v.Check(input.Email != "", "email", "must be provided")
	v.Check(validator.Matches(input.Email, validator.EmailRX), "email", "must be valid")

	v.Check(input.Password != "", "password", "must be provided")
	v.Check(len(input.Password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(input.Password) <= 72, "password", "must not be more than 72 bytes long")
}

func (svc *AuthService) RegisterUser(ctx context.Context, input *RegisterUserInput) (*model.UserRecord, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)

	if err != nil {
		return nil, err
	}

	user := &model.UserRecord{Name: input.Name, Email: input.Email, PasswordHash: hash, Activated: false}

	user, err = svc.userModel.Insert(ctx, svc.db, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (input LoginUserInput) Validate(v *validator.Validator) {
	v.Check(input.Email != "", "email", "must be provided")
	v.Check(validator.Matches(input.Email, validator.EmailRX), "email", "must be valid")
	v.Check(input.Password != "", "password", "must be provided")
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var ErrAccountActivationRequired = errors.New("user activation required")

func (svc *AuthService) LoginUser(ctx context.Context, input LoginUserInput) (*TokenDTO, error) {
	user, err := svc.userModel.FindByEmail(ctx, svc.db, input.Email)

	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// prevent user login if user is not activated
	if !user.Activated {
		fmt.Println("here:", user.Activated)
		return nil, ErrAccountActivationRequired
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(input.Password))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return nil, ErrInvalidCredentials
		default:
			return nil, err
		}
	}

	token, err := svc.tokenSvc.NewToken(ctx, user.Id, time.Minute*60*24, ScopeAuthentication)

	if err != nil {
		return nil, err
	}

	return token, nil

}
