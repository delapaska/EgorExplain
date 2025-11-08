package number

import (
	"github.com/delapaska/EgorExplain/internal/models/number"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Add(c *gin.Context) {
	var req number.NumberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	err := h.service.Add(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
