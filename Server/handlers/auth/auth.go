package auth

import (
	utils "Revisor/Utils"
	"Revisor/db"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// This function extracts detail of logged-in user using token sent by frontend
func Login(c *gin.Context) {
	//get code from json body
	var Data struct {
		Code string `json:"code"`
	}
	//bind json data to struct
	err := c.ShouldBindJSON(&Data)
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
	//expiry of token
	tokenExpiry := token.Expiry.Format(time.RFC3339)
	//This will tell frontend to start a timer
	//After the time completes(token expires) it will mark user as logged out

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
	//before storing data in db we have to confirm that duplicate data isn't stored already
	conn := db.GetDB()
	query := "select name from revisor.user where email = ?"
	result := conn.QueryRow(query, userInfo.Email)
	var name string
	err = result.Scan(&name)

	//prepare data to store in session
	sessionData := []utils.SessionKeyValue{
		{Key: "name", Value: userInfo.Name},
		{Key: "email", Value: userInfo.Email},
		{Key: "token", Value: token.AccessToken},
		{Key: "tokenExpiresAt", Value: tokenExpiry},
	}

	if err == sql.ErrNoRows {
		//store data in db
		err := db.InsertUser(conn, userInfo.Email, userInfo.FamilyName, userInfo.GivenName, userInfo.Id, userInfo.Name, userInfo.Picture, userInfo.VerifiedEmail)
		if err != nil {
			log.Printf("Failed to insert user %v\n", err)
			return
		}
		log.Println("User saved")
		//log.Println(userInfo.Email) //debugging line
		err = utils.SessionSet(c, sessionData)
		if err != nil {
			log.Printf("Failed to save session %v", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": userInfo, "info": "User saved", "token": token.AccessToken, "tokenExpiresAt": tokenExpiry})
		return
	} else if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Failed to scan row %v", err)})
		log.Printf("Failed to scan row %v", err)
		return
	}
	log.Println("user already exist")
	//log.Println(userInfo.Email) //debugging line
	err = utils.SessionSet(c, sessionData)
	if err != nil {
		log.Printf("Failed to save session %v", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": gin.H{"name": userInfo.Name, "email": userInfo.Email}, "info": "User already exists", "token": token.AccessToken, "tokenExpiresAt": tokenExpiry})
}

// This function revokes token to mark user as logged out
func Logout(c *gin.Context) {
	session := sessions.Default(c)

	//get token from json body
	var Data struct {
		Token string `json:"token"`
	}
	//bind json data to struct
	err := c.ShouldBindJSON(&Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		log.Printf("No token provided %v\n", err)
		return
	}

	//Let's revoke token by sending request to endpoint
	endpoint := "https://oauth2.googleapis.com/revoke"
	data := url.Values{}
	data.Set("token", Data.Token)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create http request %v\n", err)})
		log.Printf("Failed to create http request %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send http request %v\n", err)})
		log.Printf("Failed to send request %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to revoke token .Status : %s\n", resp.Status)})
		fmt.Printf("Failed to revoke token. Status: %s\n", resp.Status)
		fmt.Println(Data.Token)
		return
	}
	session.Clear()
	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"info": "Token successfully revoked . User logged out"})
	fmt.Println("Token successfully revoked.")
}

// This function sends user's data from session if session has data
func Me(c *gin.Context) {
	//fetch all data from session
	session := sessions.Default(c)
	sessionName := session.Get("name")
	sessionEmail := session.Get("email")
	sessionToken := session.Get("token")
	sessionTokenExpiresAt := session.Get("tokenExpiresAt")

	//To prevent panic we are using , ok syntax
	email, ok1 := sessionEmail.(string)
	name, ok2 := sessionName.(string)
	token, ok3 := sessionToken.(string)
	tokenExpiresAt, ok4 := sessionTokenExpiresAt.(string)
	if (!ok1 || email == "") || (!ok2 || name == "") || (!ok3 || token == "") || (!ok4 || tokenExpiresAt == "") {
		log.Println("Session is empty.User is not logged in")
		c.JSON(http.StatusUnauthorized, gin.H{"info": "You have to Login first"})
		return
	}
	log.Println("Data has been sent...")
	c.JSON(http.StatusOK, gin.H{"user": gin.H{"name": name, "email": email}, "token": token, "tokenExpiresAt": tokenExpiresAt})
}
