package flashcard

import (
	"Revisor/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// this function stores flashcard data in database
// Data is received from frontend
func StoreFlashcardData(c *gin.Context) {
	session := sessions.Default(c)
	//get data from json body
	var Data struct {
		Topic         string          `json:"topic"`
		FlashCardData json.RawMessage `json:"flashdata"`
	}
	//bind json data to struct
	err := c.ShouldBindJSON(&Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data provided"})
		log.Printf("No data provided %v\n", err)
		return
	}
	//check if received data is empty
	if Data.Topic == "" {
		log.Println("Topic is required")
		c.JSON(http.StatusBadRequest, gin.H{"info": "Topic is required"})
		return
	}

	//decode data
	var decodedData []map[string]string
	err = json.Unmarshal(Data.FlashCardData, &decodedData)
	if err != nil {
		log.Printf("Failed to decode data %v", err)
		return
	}

	var isEmpty bool
	for _, item := range decodedData {
		if strings.TrimSpace(item["heading"]) == "" || strings.TrimSpace(item["value"]) == "" {
			isEmpty = true
			break
		}
	}
	if isEmpty {
		log.Println("Heading and Text required")
		c.JSON(http.StatusBadRequest, gin.H{"info": "Both Heading and Text are required"})
		return
	}
	//fmt.Println(decodedData)
	sessionEmail := session.Get("email")
	fmt.Println(sessionEmail)
	//To prevent panic we are using , ok syntax
	email, ok := sessionEmail.(string)
	if !ok || email == "" {
		log.Println("Email in session is empty.User is not logged in")
		fmt.Println(email)
		c.JSON(http.StatusUnauthorized, gin.H{"info": "You have to Login first to create flashcard"})
		return
	}
	//create a db connection
	conn := db.GetDB()
	err = db.InsertFlashCardData(conn, email, Data.Topic, decodedData)
	if err != nil {
		log.Printf("Failed to store flashCardData in database %v", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"info": "FlashCard Data has been stored in database"})
	log.Println("FlashCard Data has been stored in database")
}

// this function sends flashcard data to frontend
func SendFlashCardData(c *gin.Context) {
	session := sessions.Default(c)
	sessionEmail := session.Get("email")
	//To prevent panic we are using , ok syntax
	email, ok := sessionEmail.(string)
	if !ok || email == "" {
		log.Println("Email in session is empty.User is not logged in")
		fmt.Println(email)
		c.JSON(http.StatusUnauthorized, gin.H{"info": "You have to Login first to create flashcard"})
		return
	}

	//create a db connection
	conn := db.GetDB()
	//fetch data->
	var flashCardData []db.FlashCardData
	flashCardData, err := db.FetchFlashCardData(conn, email)
	if err != nil {
		log.Printf("Error in fetching flashcard data %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error in fetching flashcard data %v", err)})
		return
	}
	log.Println("flashcard Data sent")
	c.JSON(http.StatusOK, gin.H{"flashCardData": flashCardData})
}
