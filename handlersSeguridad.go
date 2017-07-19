package main

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	db, err := sql.Open("posgres", os.Getenv("DATABASE_URL"))
}
