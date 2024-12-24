package router

import (
	"net/http"
	"project01/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/similarity", func(c *gin.Context) {
			var request struct {
				Text1 string `json:"text1" binding:"required"`
				Text2 string `json:"text2" binding:"required"`
			}
			if err := c.ShouldBindBodyWithJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			similarity, err := controller.TextSimilarity(request.Text1, request.Text2) // 这个参数是不是得用小写的？因为`json:"text1"`
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
				return 
			}
			c.JSON(http.StatusOK, gin.H{"similarity": similarity})
		})
	}

	return router
}
