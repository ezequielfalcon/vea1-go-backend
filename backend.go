package main

import "github.com/gin-gonic/gin"
import "net/http"

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hola")
	})

	router.Run(":8080")
}
