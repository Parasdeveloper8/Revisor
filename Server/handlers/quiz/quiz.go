package quiz

import (
	"Revisor/db"
	"Revisor/reusable"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// This function generates quiz using sonar model
// Data is received from frontend
// This will send quiz back to frontend
func GenerateQuiz(c *gin.Context) {
	//create a db connection
	conn := db.GetDB()
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
		NoteId    string `json:"noteId"`
	}

	//these structs will hold data (model response) and will be stored in database
	type QuesOpt struct {
		Question string
		Options  []string
	}
	type ModelResJSON struct {
		QuizId   string
		Quesopts []QuesOpt
	}

	//Bind json data to ReceivedData
	err := c.ShouldBindJSON(&ReceivedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No JSON data provided or incomplete data provided"})
		log.Printf("No JSON data provided or incomplete data provided %v\n", err)
		return
	}
	//if NoteId is already present in database
	query := "select quizId,question,options from revisor.quiz where quiz.noteId = ?"
	result := conn.QueryRow(query, ReceivedData.NoteId)
	//if any row returned then send that returned row to frontend
	//this prevents model API calling
	type Response struct {
		QuizId    string `db:"quizId"`
		Question  []byte `db:"question"`
		Options   []byte `db:"options"`
		TopicName string
	}
	var response Response
	err = result.Scan(&response.QuizId, &response.Question, &response.Options)
	if err == nil {

		//convert json data from db into string slice ---->
		opts, err := reusable.UnmarshalJSONtoStringSlice(response.Options)
		if err != nil {
			fmt.Printf("Failed to parse data %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error from server"})
			return
		}
		ques, err := reusable.UnmarshalJSONtoStringSlice(response.Question)
		if err != nil {
			fmt.Printf("Failed to parse data %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error from server"})
			return
		}
		var options [][]string
		chunkSize := 4
		//split slice by 4 elements each
		for i := 0; i < len(opts); i += chunkSize {
			end := i + chunkSize
			if end > len(opts) {
				end = len(opts)
			}
			options = append(options, opts[i:end])

		}
		var finalRes ModelResJSON
		for index, text := range ques {
			quesOpts := QuesOpt{
				Question: text,
				Options:  options[index],
			}
			finalRes.Quesopts = append(finalRes.Quesopts, quesOpts)
		}
		finalRes.QuizId = response.QuizId
		c.JSON(http.StatusOK, gin.H{"response": finalRes, "topic": ReceivedData.TopicName})
		return
	} else if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No matching datafound %v\n", err)
			//Do further process because no existing data found in database
		} else if err != sql.ErrNoRows {
			fmt.Printf("Failed to scan row %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error from server"})
			return
		}
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
				Give four options.
				Give right answer also like right answer is __text of right answer__.
				 Don't generate facts.
				 Questions must not be from outside of given notes.
				 Generate questions from data as provided whether data is less or much.
				 Example if data contains only some lines then questions must be based on given lines.
				 You have to treat as a teacher who asks question from notebook like from NCERT not from google.
				 Separate all questions by 1,2,...,n numbers.
				 if notes are short then you can generate less questions but questions must be from inside the topic.
				 Cover whole topic in 1 to 10 questions.`,
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
	//generate a unique id for every quiz
	quizUID := uuid.New()

	//Now in further code we will split data and store in struct
	re := regexp.MustCompile(`\d+\.\s`)
	parts := re.Split(modelRes.Choices[0].Message.Content, -1)

	var modelResJsonVar *ModelResJSON
	var answers []string
	var questions []string
	var options []string
	for _, part := range parts {
		if part == "" {
			continue // first split may be empty
		}
		//Split by right answer
		re := regexp.MustCompile(`Right answer is`)
		npart := re.Split(part, -1)
		if len(npart) < 2 {
			log.Printf("Skipping part because no answer found: %s\n", part)
			continue
		}
		queoptpart := strings.TrimSpace(npart[0])
		anspart := strings.TrimSpace(npart[1])

		answer := strings.TrimPrefix(anspart, "Right answer is ")
		answer = strings.ReplaceAll(answer, "*", "")
		answer = strings.ReplaceAll(answer, "_", "")
		answer = strings.TrimSpace(answer) //answer
		//separate questions from options
		lines := strings.Split(queoptpart, "\n")
		if len(lines) < 2 {
			log.Printf("Skipping part because options missing: %s\n", part)
			continue
		}

		que := lines[0] // first line is question

		reg := regexp.MustCompile(`^[a-d]\)\s*(.*)$`) // match a) b) c) d)
		for _, l := range lines[1:] {                 //[start:end] [starts with options]
			l = strings.TrimSpace(l)
			if l == "" { //if any option is empty
				continue
			}
			match := reg.FindStringSubmatch(l)
			if len(match) == 2 {
				options = append(options, match[1]) // add the option text
			}
		}

		quesOpts := QuesOpt{
			Question: que,
			Options:  options,
		}

		// Append instead of overwrite
		if modelResJsonVar == nil {
			modelResJsonVar = &ModelResJSON{QuizId: quizUID.String()}
		}
		modelResJsonVar.Quesopts = append(modelResJsonVar.Quesopts, quesOpts)
		answers = append(answers, answer)
		questions = append(questions, que)
	}
	//insert data of quiz including answers in database --->
	session := sessions.Default(c)
	sessionEmail := session.Get("email")

	//To prevent panic we are using , ok syntax
	email, ok := sessionEmail.(string)
	if !ok || email == "" {
		log.Println("Email in session is empty.User is not logged in")
		fmt.Println(email)
		c.JSON(http.StatusUnauthorized, gin.H{"info": "You have to Login first to save quiz"})
		return
	}
	err = db.InsertQuizData(conn, email, answers, questions, options, quizUID.String(), ReceivedData.NoteId)
	if err != nil {
		log.Printf("Failed to insert quiz data in database : %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in server side"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": modelResJsonVar, "topic": ReceivedData.TopicName})
}
