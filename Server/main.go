package main

import (
	"Revisor/cors"
	"Revisor/handlers/auth"
	"Revisor/handlers/flashcard"
	"Revisor/handlers/quiz"
	"Revisor/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	//set up cors
	cors.SetupCORS(r)

	//Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	store := cookie.NewStore([]byte(os.Getenv("SECRET_KEY")))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600,  // 1 hour
		HttpOnly: true,  // 🔒 Prevent frontend access
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("RevisorSession", store))

	r.POST("/auth/google", auth.Login)
	r.POST("/auth/logout", middleware.CheckISAuthenticated(), auth.Logout)
	r.GET("/auth/me", middleware.CheckISAuthenticated(), auth.Me)
	r.POST("/flashcard/store/data", middleware.CheckISAuthenticated(), flashcard.StoreFlashcardData)
	r.GET("/flashcard/get/data", middleware.CheckISAuthenticated(), flashcard.SendFlashCardData)
	r.POST("/generate/quiz", middleware.CheckISAuthenticated(), quiz.GenerateQuiz)
	r.POST("/evaluate/quiz", middleware.CheckISAuthenticated(), quiz.EvaluateQuiz)
	r.Run(":8080")
}
