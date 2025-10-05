package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"merch-app-codex/internal/storage/mysql"
)

type contextKey string

const (
	// ContextUserKey is used to store the authenticated user in Gin's context.
	ContextUserKey contextKey = "authenticated_user"
	// ContextTokenKey stores the active token model in the request context.
	ContextTokenKey contextKey = "authenticated_token"
)

// TokenAuthMiddleware validates bearer tokens from the Authorization header.
func TokenAuthMiddleware(repo Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenValue := strings.TrimSpace(authHeader)
		if strings.HasPrefix(strings.ToLower(tokenValue), "bearer ") {
			tokenValue = strings.TrimSpace(tokenValue[7:])
		}

		if tokenValue == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		token, user, err := repo.FindUserToken(c.Request.Context(), tokenValue)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Set(string(ContextUserKey), user)
		c.Set(string(ContextTokenKey), token)

		c.Next()
	}
}

// CurrentUser retrieves the authenticated user from the context when available.
func CurrentUser(c *gin.Context) (*mysql.User, bool) {
	value, ok := c.Get(string(ContextUserKey))
	if !ok {
		return nil, false
	}
	user, ok := value.(*mysql.User)
	return user, ok
}
