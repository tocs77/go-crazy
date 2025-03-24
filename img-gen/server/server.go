package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func setupUserRoutes(r *gin.Engine) {
	userRoute := r.Group("/user")
	{
		userRoute.GET("hello/:name", func(c *gin.Context) {
			user := c.Param("name")
			c.String(http.StatusOK, fmt.Sprintf("Hello %s", user))
		})
		userRoute.POST("/post", func(c *gin.Context) {
			body := Message{}
			if err := c.BindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"name": body.Name, "email": body.Email})
		})
		userRoute.POST("/upload", func(c *gin.Context) {
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.SaveUploadedFile(file, "upload/"+file.Filename)
			c.JSON(http.StatusOK, gin.H{"filename": file.Filename})
		})
	}
}

func router() *gin.Engine {
	r := gin.Default()
	setupUserRoutes(r)

	return r
}

func Start() {
	r := router()
	r.Run()
}
