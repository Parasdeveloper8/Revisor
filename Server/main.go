package main

import (
	"Revisor/cors"
	"Revisor/handlers/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//set up cors
	cors.SetupCORS(r)

	r.POST("/auth/google", auth.Login)
	r.POST("/auth/logout", auth.Logout)

	r.Run(":8080")
}
