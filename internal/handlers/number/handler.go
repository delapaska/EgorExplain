package number

import (
	"github.com/delapaska/EgorExplain/internal/services/number"
	"github.com/gin-gonic/gin"
)

type handler struct {
	service number.Service
}
type Handler interface {
	Add(c *gin.Context)
	RegisterRoutes(router *gin.Engine)
}

func New(service number.Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(router *gin.Engine) {

	api := router.Group("/api/v1")
	{
		api.POST("/add", h.Add)

	}

}
