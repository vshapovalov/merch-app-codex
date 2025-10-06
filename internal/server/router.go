package server

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"merch-app-codex/internal/auth"
	"merch-app-codex/internal/config"
	"merch-app-codex/internal/report"
	"merch-app-codex/internal/storage/mysql"
)

// NewRouter wires all HTTP handlers and middleware.
func NewRouter(cfg config.Config, repo *mysql.Repository, authRepo auth.Repository, reportService *report.Service) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.POST("/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := authRepo.FindUserByEmail(c.Request.Context(), req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := user.CheckPassword(req.Password); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		tokenValue, err := auth.GenerateToken()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token := &mysql.UserToken{
			UserID: user.ID,
			Token:  tokenValue,
		}
		token.SetID(mysql.NewID())

		if err := authRepo.CreateUserToken(c.Request.Context(), token); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token.Token})
	})

	secured := api.Group("")
	secured.Use(auth.TokenAuthMiddleware(authRepo))

	registerEntityRoutes[mysql.User, *mysql.User](secured, repo, entityFactory[mysql.User, *mysql.User]{
		path: "/users",
		new:  func() *mysql.User { return &mysql.User{} },
	})

	registerEntityRoutes[mysql.Company, *mysql.Company](secured, repo, entityFactory[mysql.Company, *mysql.Company]{
		path: "/companies",
		new:  func() *mysql.Company { return &mysql.Company{} },
	})

	registerEntityRoutes[mysql.RetailPoint, *mysql.RetailPoint](secured, repo, entityFactory[mysql.RetailPoint, *mysql.RetailPoint]{
		path: "/retail-points",
		new:  func() *mysql.RetailPoint { return &mysql.RetailPoint{} },
	})

	registerEntityRoutes[mysql.Brand, *mysql.Brand](secured, repo, entityFactory[mysql.Brand, *mysql.Brand]{
		path: "/brands",
		new:  func() *mysql.Brand { return &mysql.Brand{} },
	})

	registerEntityRoutes[mysql.Category, *mysql.Category](secured, repo, entityFactory[mysql.Category, *mysql.Category]{
		path: "/categories",
		new:  func() *mysql.Category { return &mysql.Category{} },
	})

	registerEntityRoutes[mysql.Product, *mysql.Product](secured, repo, entityFactory[mysql.Product, *mysql.Product]{
		path: "/products",
		new:  func() *mysql.Product { return &mysql.Product{} },
	})

	registerEntityRoutes[mysql.Visit, *mysql.Visit](secured, repo, entityFactory[mysql.Visit, *mysql.Visit]{
		path: "/visits",
		new: func() *mysql.Visit {
			return &mysql.Visit{VisitedAt: time.Now()}
		},
	})

	registerEntityRoutes[mysql.VisitItem, *mysql.VisitItem](secured, repo, entityFactory[mysql.VisitItem, *mysql.VisitItem]{
		path: "/visit-items",
		new:  func() *mysql.VisitItem { return &mysql.VisitItem{} },
	})

	reports := secured.Group("/reports")
	reports.GET("/companies/:id/visits", func(c *gin.Context) {
		summary, err := reportService.CompanyVisitSummary(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, summary)
	})

	if cfg.StaticDir != "" {
               if info, err := os.Stat(cfg.StaticDir); err == nil && info.IsDir() {
                       router.StaticFS("/static", gin.Dir(cfg.StaticDir, true))

                       router.NoRoute(func(c *gin.Context) {
                               if strings.HasPrefix(c.Request.URL.Path, "/api") {
                                       c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
                                       return
                               }

                               c.File(filepath.Join(cfg.StaticDir, "index.html"))
                       })
               }
       }

	return router
}
