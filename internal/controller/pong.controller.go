package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
type PongController struct{}

func NewPongController() *PongController{
	return &PongController{}
}
func (uc *PongController) Pong(c *gin.Context) {
	name := c.Param("name")
	uid := c.DefaultQuery("uid", "go ecommerce")
	c.JSON(http.StatusOK, gin.H{
		"message": name + uid,
		"users":   []string{"quan", "thu", "khoi"},
	})
}