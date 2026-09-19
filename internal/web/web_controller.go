package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebController struct {
}

func NewWebController() *WebController {
	return &WebController{}
}

func (controller *WebController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "main.gohtml", gin.H{})
}
