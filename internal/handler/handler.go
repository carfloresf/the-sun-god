package handler

import (
	"github.com/carfloresf/the-sun-god/internal/constants"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"github.com/carfloresf/the-sun-god/internal/model"
	"github.com/carfloresf/the-sun-god/internal/service"
)

type Service interface {
	CreatePost(post model.Post) (uuid.UUID, error)
	UpdatePostScore(pID uuid.UUID, score int) error
	GetFeed(page, limit int) (service.GetFeedResponse, error)
}

type PostHandler struct {
	Service Service
}

type UpdateScoreRequest struct {
	ID    uuid.UUID `json:"id" binding:"required"`
	Score int       `json:"score" binding:"required"`
}

func (p *PostHandler) UpdateScore(c *gin.Context) {
	var request UpdateScoreRequest
	if err := c.BindJSON(&request); err != nil {
		log.Errorf("validation errors: %+v", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	err := p.Service.UpdatePostScore(request.ID, request.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"error": nil})
}

type CreateResponse struct {
	ID    uuid.UUID `json:"id"`
	Error error     `json:"error,omitempty"`
}

func (p *PostHandler) Create(c *gin.Context) {
	var post model.Post

	if err := c.BindJSON(&post); err != nil {
		log.Errorf("validation errors: %+v", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	pID, err := p.Service.CreatePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.IndentedJSON(http.StatusCreated, CreateResponse{
		ID:    pID,
		Error: nil,
	})
}

func (p *PostHandler) GetFeed(c *gin.Context) {
	page := cast.ToInt(c.Query("page"))
	limit := cast.ToInt(c.Query("limit"))

	if page == 0 {
		page = constants.DefaultPage
	}

	if limit == 0 {
		limit = constants.DefaultLimit
	}

	log.Printf("page: %d, limit: %d", page, limit)

	feed, err := p.Service.GetFeed(cast.ToInt(page), cast.ToInt(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, feed)
}
