package main

import (
	"database/sql"
	"log"
	"os"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin/json"
	_ "github.com/lib/pq"
)

type usuario struct {
	Nombre  string  `form:"nombre" binding:"required"`
	Clave   string  `form:"clave" binding:"required"`
	Cliente string  `form:"cliente" binding:"required"`
}

func usuariosNuevo(c *gin.Context) {
	header := c.GetHeader("x-access-token")
	if header == "" {
		c.JSON(401, gin.H{"resultado": false, "mensaje": "No token provided"})
	} else {
		decoded, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			c.JSON(401, gin.H{"resultado": false, "mensaje": "Error de autenticación"})
		} else {
			var tokenDeco map[string]interface{}
			if err := json.Unmarshal([]byte(decoded), &tokenDeco); err != nil {
				log.Println("Error en decoded de token!")
				c.JSON(500, gin.H{"resultado": false, "mensaje": "Error interno decodificando el token"})
			} else {
				usuario := tokenDeco["usuario"].(string)
				roles := tokenDeco["roles"].([]string)
				log.Println(usuario)
				log.Println(roles)
			}
		}
	}
}