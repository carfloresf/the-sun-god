package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/carfloresf/the-sun-god/internal/handler"
	"github.com/carfloresf/the-sun-god/internal/model"
	"github.com/carfloresf/the-sun-god/internal/service"
)

type storageInterface interface {
	CreatePost(post model.Post) (uuid.UUID, error)
	UpdatePostScore(pID uuid.UUID, score int) error
	GetPosts(limit, offset int) ([]model.Post, int, error)
	GetPromotedPosts(limit int) ([]model.Post, error)
}

func InitRouter(storage storageInterface) (*gin.Engine, error) {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(JSONMiddleware())

	healthController := NewHealthController()
	postHandler := handler.PostHandler{
		Service: &service.Service{Storage: storage},
	}

	r.GET("/health", healthController.Status)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := r.Group("/v1")
	{
		h.POST("/post/score", postHandler.UpdateScore)
		h.POST("/post", postHandler.Create)
		h.GET("/feed", postHandler.GetFeed)
	}

	return r, nil
}

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}
