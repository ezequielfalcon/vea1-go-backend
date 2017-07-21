package main

import (

	"github.com/gin-gonic/gin"
	"os"
	_ "github.com/lib/pq"
)

func main() {
	router := gin.Default()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hola")
	})

	router.POST("/login", seguridadLogin)

	router.Run(":" + port)
}
