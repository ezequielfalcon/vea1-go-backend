package main

import (
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin/json"
	_ "github.com/lib/pq"
)

type login struct {
	Usuario     string `form:"usuario" binding:"required"`
	Clave       string `form:"clave" binding:"required"`
	HashCliente string `form:"cliente" binding:"required"`
}

type TokenStr struct {
	Usuario     string  `json:"usuario"`
	Roles       []string  `json:"roles"`
}

func seguridadLogin(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{"resultado": false, "mensaje": "Error en variables de entorno"})
	} else {
		var form login
		if c.Bind(&form) == nil {
			var idClienteInt string
			var clave string
			err := db.QueryRow("SELECT id_cliente_int, clave " +
				"FROM usuarios WHERE nombre = $1;", form.Usuario).Scan(&idClienteInt, &clave)
			if err != nil {
				log.Println(err)
				c.JSON(500, gin.H{"resultado": false, "mensaje": err})
			} else {
				clienteHashBytes := []byte(form.HashCliente)
				clienteBytes := []byte(idClienteInt)
				comparar := bcrypt.CompareHashAndPassword(clienteHashBytes, clienteBytes)
				if comparar != nil {
					log.Println(comparar)
					c.JSON(401, gin.H{"resultado": false, "mensaje": "Usuario con hash de cliente inválido"})
				} else {
					usuarioHashBytes := []byte(clave)
					usuarioClaveBytes := []byte(form.Clave)
					compararClave := bcrypt.CompareHashAndPassword(usuarioHashBytes, usuarioClaveBytes)
					if compararClave != nil {
						c.JSON(401, gin.H{"resultado": false, "mensaje": "Usuario o contraseña no válida"})
					} else {
						roles, err := db.Query("SELECT roles.nombre FROM roles INNER JOIN roles_por_usuario " +
							"ON roles.id = roles_por_usuario.id_rol " +
							"INNER JOIN usuarios on roles_por_usuario.usuario = usuarios.nombre " +
							"WHERE usuarios.nombre = $1;"	, form.Usuario)
						if err != nil {
							c.JSON(500, gin.H{"resultado": false, "mensaje": err})
						} else {
							rolesUsuario := make([]string, 0)
							for roles.Next() {
								var nombre string
								roles.Scan(&nombre)
								rolesUsuario = append(rolesUsuario, nombre)
							}
							tokenRaw := &TokenStr{
								Usuario:    form.Usuario,
								Roles:      rolesUsuario }
							token := jwt.New(jwt.SigningMethodHS512)
							tokenString, err := json.Marshal(tokenRaw)
							if err != nil {
								log.Println(err)
								c.JSON(500, gin.H{"resultado": false, "mensaje": "Error en JSON"})
							} else {
								token.Raw = string(tokenString)
								secret := os.Getenv("JWT_SECRET")
								tokenFinal, err := token.SignedString([]byte(secret))
								if err != nil {
									log.Println(err)
									c.JSON(500, gin.H{"resultado": false, "mensaje": "Error en JWT"})
								} else {
									c.JSON(200, gin.H{"resultado": true, "mensaje": "Sesión iniciada!", "token": tokenFinal})
								}
							}
						}
					}
				}
			}
		} else {
			c.JSON(400, gin.H{"resultado": false, "mensaje": "Faltan datos para iniciar sesión"})
		}
	}
}
