package handler

import (
	"net/http"
	"project/ikan"

	"github.com/gin-gonic/gin"
)

type ikanHandler struct {
	ikanService ikan.Service
}

func NewIkanHandler(ikanService ikan.Service) *ikanHandler {
	return &ikanHandler{ikanService}
}

func (h *ikanHandler) RegisterIkan(c *gin.Context) {
	var body ikan.Ikan
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Body is invalid.",
			"error":   err.Error(),
		})
		return
	}
	newIkan, err := h.ikanService.TambahIkan(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error when inserting into the database.",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Ikan Berhasil Ditambahkan",
		"status":  "Sukses",
		"data":    newIkan,
	})
}
