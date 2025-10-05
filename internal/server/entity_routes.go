package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"merch-app-codex/internal/storage/mysql"
)

type entityFactory[Model any, Ptr interface {
	*Model
	mysql.Entity
}] struct {
	path string
	new  func() Ptr
}

func registerEntityRoutes[Model any, Ptr interface {
	*Model
	mysql.Entity
}](group *gin.RouterGroup, repo *mysql.Repository, factory entityFactory[Model, Ptr]) {
	route := group.Group(factory.path)

	route.POST("", func(c *gin.Context) {
		entity := factory.new()
		if err := c.ShouldBindJSON(entity); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if entity.GetID() == "" {
			entity.SetID(mysql.NewID())
		}

		if err := repo.Create(c.Request.Context(), entity); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, entity)
	})

	route.GET("", func(c *gin.Context) {
		var list []Model
		if err := repo.List(c.Request.Context(), &list); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	})

	route.GET(":id", func(c *gin.Context) {
		entity := factory.new()
		if err := repo.FindByID(c.Request.Context(), entity, c.Param("id")); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, entity)
	})

	route.PUT(":id", func(c *gin.Context) {
		id := c.Param("id")
		entity := factory.new()
		if err := repo.FindByID(c.Request.Context(), entity, id); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if err := c.ShouldBindJSON(entity); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		entity.SetID(id)
		if err := repo.Update(c.Request.Context(), entity); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, entity)
	})

	route.DELETE(":id", func(c *gin.Context) {
		if err := repo.DeleteByID(c.Request.Context(), factory.new(), c.Param("id")); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
