package main

import (
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin/json"
)

type login struct {
	usuario     string `form:"usuario" binding:"required"`
	clave       string `form:"clave" binding:"required"`
	hashCliente string `form:"cliente" binding:"required"`
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
				"FROM usuarios WHERE nombre = $1;", form.usuario).Scan(&idClienteInt, &clave)
			if err != nil {
				log.Fatal(err)
				c.JSON(500, gin.H{"resultado": false, "mensaje": err})
			} else {
				clienteDbByte := []byte(idClienteInt)
				clienteDbHash, err := bcrypt.GenerateFromPassword(clienteDbByte, bcrypt.DefaultCost)
				if err != nil {
					log.Fatal(err)
					c.JSON(500, gin.H{"resultado": false, "mensaje": "Error en BCRYPT"})
				} else {
					clienteDbHashString := string(clienteDbHash)
					if clienteDbHashString == form.hashCliente {
						claveUsuarioByte := []byte(form.clave)
						claveUsuarioHash, err := bcrypt.GenerateFromPassword(claveUsuarioByte, bcrypt.DefaultCost)
						if err != nil {
							log.Fatal(err)
							c.JSON(500, gin.H{"resultado": false, "mensaje": "Error en BCRYPT"})
						} else {
							hashString := string(claveUsuarioHash)
							if hashString == clave {
								roles, err := db.Query("SELECT roles.nombre FROM roles INNER JOIN roles_por_usuario " +
									"ON roles.id = roles_por_usuario.id_rol " +
									"INNER JOIN usuarios on roles_por_usuario.usuario = usuarios.nombre " +
									"WHERE usuarios.nombre = $1;"	, form.usuario)
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
										Usuario:    form.usuario,
										Roles:      rolesUsuario }
									token := jwt.New(jwt.SigningMethodHS512)
									tokenString, err := json.Marshal(tokenRaw)
									if err != nil {
										log.Fatal(err)
										c.JSON(500, gin.H{"resultado": false, "mensaje": "Error en JSON"})
									}
									token.Raw = string(tokenString)
									secret := os.Getenv("JWT_SECRET")
									tokenFinal, err := token.SignedString(secret)
									if err != nil {
										log.Fatal(err)
										c.JSON(500, gin.H{"resultado": false, "mensaje": "Error en JWT"})
									}
									c.JSON(200, gin.H{"resultado": true, "mensaje": "Sesión iniciada!", "token": tokenFinal})
								}
							} else {
								c.JSON(401, gin.H{"resultado": false, "mensaje": "Usuario o contraseña no válida"})
							}
						}
					} else {
						c.JSON(401, gin.H{"resultado": false, "mensaje": "Usuario con hash de cliente inválido"})
					}
				}
			}
		} else {
			c.JSON(400, gin.H{"resultado": false, "mensaje": "Faltan datos para iniciar sesión"})
		}
	}
}
