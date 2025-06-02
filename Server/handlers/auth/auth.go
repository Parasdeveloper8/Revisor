package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// This function extracts detail of logged-in user using token sent by frontend
func Auth(c *gin.Context) {
	//Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//get code from json body
	var Data struct {
		Code string `json:"code"`
	}
	//bind json data to struct
	err = c.ShouldBindJSON(&Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
		log.Printf("No code provided %v\n", err)
		return
	}

	//struct for storing env variables ' value
	type OauthCredentials struct {
		ClientId     string `env:"CLIENT_ID"`
		ClientSecret string `env:"CLIENT_SECRET"`
		RedirectURI  string `env:"REDIRECT_URI"`
		GrantType    string `env:"GRANT_TYPE"`
	}
	//store values
	var authcredentials = &OauthCredentials{
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURI:  os.Getenv("REDIRECT_URI"),
		GrantType:    os.Getenv("GRANT_TYPE"),
	}

	if authcredentials.ClientId == "" && authcredentials.ClientSecret == "" && authcredentials.GrantType == "" && authcredentials.RedirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty env variables", "env": authcredentials})
		log.Println("Empty env variables")
		return
	}
	//oauth configurations
	var OauthConfig = &oauth2.Config{
		ClientID:     authcredentials.ClientId,
		ClientSecret: authcredentials.ClientSecret,
		RedirectURL:  authcredentials.RedirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	//Exchange code for token
	token, err := OauthConfig.Exchange(c, Data.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code "})
		log.Printf("Failed to exchange code %v\n", err)
		return
	}

	//fetch user info using token
	client := OauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		log.Println("User info request failed:", err)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		log.Printf("Failed to decode user info %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userInfo})
}
