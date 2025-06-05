package auth

import (
	"Revisor/db"
	"database/sql"
	"encoding/json"
	"fmt"
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

	//struct to contain user data
	type UserInfo struct {
		Email         string `json:"email"`
		FamilyName    string `json:"family_name"`
		GivenName     string `json:"given_name"`
		Id            string `json:"id"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}
	var userInfo = &UserInfo{}
	err = json.NewDecoder(resp.Body).Decode(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		log.Printf("Failed to decode user info %v\n", err)
		return
	}

	//Store data in db once
	//before storing data in db we have to confirm that duplicate data isn't stroed already
	conn := db.GetDB()
	query := "select name from revisor.user where email = ?"
	result := conn.QueryRow(query, userInfo.Email)
	var name string
	err = result.Scan(&name)
	if err == sql.ErrNoRows {
		//store data in db
		err := db.InsertUser(conn, userInfo.Email, userInfo.FamilyName, userInfo.GivenName, userInfo.Id, userInfo.Name, userInfo.Picture, userInfo.VerifiedEmail)
		if err != nil {
			log.Printf("Failed to insert user %v\n", err)
			return
		}
		log.Println("User saved")
		c.JSON(http.StatusOK, gin.H{"user": userInfo, "info": "User saved", "login": true})
		return
	} else if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Failed to scan row %v", err)})
		log.Printf("Failed to scan row %v", err)
		return
	}
	log.Println("user already exist")
	c.JSON(http.StatusOK, gin.H{"user": gin.H{"name": userInfo.Name, "email": userInfo.Email}, "info": "User already exists", "login": true})
}
