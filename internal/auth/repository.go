package auth

import (
	"context"

	"merch-app-codex/internal/storage/mysql"
)

// Repository defines the persistence operations required by the auth package.
type Repository interface {
	FindUserByEmail(ctx context.Context, email string) (*mysql.User, error)
	CreateUserToken(ctx context.Context, token *mysql.UserToken) error
	FindUserToken(ctx context.Context, token string) (*mysql.UserToken, *mysql.User, error)
	DeleteUserToken(ctx context.Context, token string) error
}
