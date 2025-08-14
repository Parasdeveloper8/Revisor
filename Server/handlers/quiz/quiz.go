package quiz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// This function generates quiz using sonar model
// Data is received from frontend
// This will send quiz back to frontend
func GenerateQuiz(c *gin.Context) {
	apiKey := os.Getenv("PERPLEXITY_API_KEY")
	if apiKey == "" {
		log.Printf("Empty PERPLEXITY_API_KEY in .env file")
		return
	}

	type Data struct {
		Heading string `json:"Heading"`
		Value   string `json:"Value"`
	}
	//struct to hold frontend received data
	var ReceivedData struct {
		TopicName string `json:"topicName"`
		Notes     []Data `json:"data"`
	}
	//Bind json data to ReceivedData
	err := c.ShouldBindJSON(&ReceivedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No JSON data provided or incomplete data provided"})
		log.Printf("No JSON data provided or incomplete data provided %v\n", err)
		return
	}
	//Send data to model and get quiz data
	url := "https://api.perplexity.ai/chat/completions"
	//Convert notes into a readable string
	var notesBuilder strings.Builder
	for _, note := range ReceivedData.Notes {
		notesBuilder.WriteString(fmt.Sprintf("Heading: %s\nValue: %s\n\n", note.Heading, note.Value))
	}
	noteContent := notesBuilder.String()

	//Create request payload
	requestPayload := map[string]interface{}{
		"model": "sonar",
		"messages": []map[string]string{
			{
				"role": "system",
				"content": `
				You are a quiz generator. 
				Generate quiz questions based on the data provided by the user from the given notes.
				 Don't give answers.
				 Don't generate facts.
				 Questions must not be from outside of given notes.
				 Generate questions from data as provided whether data is less or much.
				 Example if data contains only some lines then questions must be based on given lines.
				 You have to treat as a teacher who asks question from notebook like from NCERT not from google.
				 Separate all questions by 1,2,...,n numbers.
				 Do not generate more than 15 questions and less than 5 questions.`,
			},
			{
				"role":    "user",
				"content": "Here is the data:\n\n" + noteContent,
			},
		},
	}

	//Convert payload to JSON
	jsonBody, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Failed to marshal JSON payload: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error while preparing request"})
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		log.Printf("Failed to build HTTP request : %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build HTTP request"})
		return
	}
	authHeaderValue := "Bearer " + apiKey //value of Authorization header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeaderValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send HTTP request %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send HTTP request"})
		return
	}
	defer resp.Body.Close()

	// Read body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}
	//struct to hold response from model
	type ModelResponse struct {
		Choices []struct {
			Message struct {
				Content string
			}
		}
	}
	modelRes := &ModelResponse{}
	err = json.Unmarshal(respBody, &modelRes)
	if err != nil {
		log.Printf("Failed to unmarshal %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error from Server"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": modelRes, "topic": ReceivedData.TopicName})
}
