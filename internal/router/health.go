package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController interface {
	Status(c *gin.Context)
}

type healthController struct{}

func NewHealthController() *healthController {
	return &healthController{}
}

func (h healthController) Status(c *gin.Context) {
	c.String(http.StatusOK, "Still Working!")
}
