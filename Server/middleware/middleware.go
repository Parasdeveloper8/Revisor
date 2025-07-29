package middleware

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// This middleware checks is user authenticated or not
func CheckISAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the session
		session := sessions.Default(c)
		// Retrieve the email from the session
		sessionEmail, ok := session.Get("email").(string)
		if !ok || sessionEmail == "" {
			// If email is not set or invalid
			log.Println("Email doesn't exist in session")
			c.JSON(http.StatusNotFound, gin.H{"error": "You have to login first !"})
			c.Abort()
			return
		}
		c.Next()
	}
}
