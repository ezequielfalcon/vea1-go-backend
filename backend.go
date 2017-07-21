package main

import (

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hola")
	})

	router.POST("/login", seguridadLogin)

	router.Run(":8080")
}
