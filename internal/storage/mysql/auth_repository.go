package mysql

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// AuthRepository handles persistence for authentication-related models.
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository constructs an AuthRepository backed by GORM.
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// FindUserByEmail retrieves a user by their email address.
func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUserToken stores a new authentication token for the given user.
func (r *AuthRepository) CreateUserToken(ctx context.Context, token *UserToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindUserToken loads a user token and its associated user by token string.
func (r *AuthRepository) FindUserToken(ctx context.Context, token string) (*UserToken, *User, error) {
	var userToken UserToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&userToken).Error; err != nil {
		return nil, nil, err
	}

	if userToken.ExpiresAt != nil && userToken.ExpiresAt.Before(time.Now()) {
		// Token expired; remove it eagerly.
		_ = r.db.WithContext(ctx).Delete(&userToken).Error
		return nil, nil, gorm.ErrRecordNotFound
	}

	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", userToken.UserID).First(&user).Error; err != nil {
		return nil, nil, err
	}

	return &userToken, &user, nil
}

// DeleteUserToken removes a persisted token by its string value.
func (r *AuthRepository) DeleteUserToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&UserToken{}).Error
}

// DB exposes the underlying database handle for advanced queries.
func (r *AuthRepository) DB() *gorm.DB {
	return r.db
}
