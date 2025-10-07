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
		opts, err := reusable.UnmarshalJSONtoStringSlice[[][]string](response.Options)
		if err != nil {
			fmt.Printf("Failed to parse data %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error from server"})
			return
		}
		ques, err := reusable.UnmarshalJSONtoStringSlice[[]string](response.Question)
		if err != nil {
			fmt.Printf("Failed to parse data %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error from server"})
			return
		}
		var finalRes ModelResJSON
		for index, text := range ques {
			quesOpts := QuesOpt{
				Question: text,
				Options:  opts[index],
			}
			finalRes.Quesopts = append(finalRes.Quesopts, quesOpts)
		}
		finalRes.QuizId = response.QuizId
		c.JSON(http.StatusOK, gin.H{"response": finalRes, "topic": ReceivedData.TopicName})
		return
	} else if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No matching data found %v\n", err)
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
            You are a strict quiz generator.
            Generate quiz questions strictly from the data provided by the user. 
            Follow these rules exactly:

            1. Generate 1 to 10 questions depending on the length of the notes.
            2. Each question must have **exactly 4 options**, labeled A), B), C), D) — no more, no less.
            3. Only **one option should be correct**. Mention the correct option clearly like: "Right answer is __<text>__".
            4. Do not include any additional facts, explanations, or content outside the given notes.
            5. Questions must cover the given notes and nothing else.
            6. Separate questions by numbering: 1, 2, 3, ...
            7. Output must be formatted clearly, easy to parse for storing in database.
            8. Never generate more than 4 options per question.
            `,
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
	//fmt.Printf("Raw response : %v\n", string(respBody))
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
	var allOpts [][]string
	for _, part := range parts {
		var options []string
		if part == "" {
			continue // first split may be empty
		}
		//Split by right answer
		re := regexp.MustCompile(`Right answer is`)
		npart := re.Split(part, -1)
		fmt.Printf("Npart : %v\n", npart)
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

		reg := regexp.MustCompile(`^[A-D]\)\s*(.*)$`) // match A) B) C) D)
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

		allOpts = append(allOpts, options)
		quesOpts := QuesOpt{
			Question: que,
			Options:  options,
		}
		//fmt.Printf("Options :%v\n", options)
		// Append instead of overwrite
		if modelResJsonVar == nil {
			modelResJsonVar = &ModelResJSON{QuizId: quizUID.String(), Quesopts: []QuesOpt{}}
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
	err = db.InsertQuizData(conn, email, answers, questions, allOpts, quizUID.String(), ReceivedData.NoteId)
	if err != nil {
		log.Printf("Failed to insert quiz data in database : %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in server side"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": modelResJsonVar, "topic": ReceivedData.TopicName})
}

// this function sends user their marks of quiz given by them
// Answers from the db and user's answers will be compared using a quizID
func EvaluateQuiz(c *gin.Context) {
	//create a db connection
	conn := db.GetDB()

	//struct to hold frontend received data
	var ReceivedData struct {
		UserAnswers []string `json:"userAnswers"`
		QuizID      string   `json:"quizId"`
		Time        int      `json:"time"`
	}

	//Bind json data to ReceivedData
	err := c.ShouldBindJSON(&ReceivedData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No JSON data provided or incomplete data provided"})
		log.Printf("No JSON data provided or incomplete data provided %v\n", err)
		return
	}
	//fetch answers from db using QuizID
	query := "select answers from revisor.quiz where quizId = ?"
	result := conn.QueryRow(query, ReceivedData.QuizID)
	type DBResponse struct {
		Answers []byte `json:"answers"`
	}
	var dbresponse DBResponse
	err = result.Scan(&dbresponse.Answers)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No matching data found %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error from server"})
			return
		} else if err != sql.ErrNoRows {
			fmt.Printf("Failed to scan row %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error from server"})
			return
		}
	}
	//convert json data from db into string slice ---->
	answers, err := reusable.UnmarshalJSONtoStringSlice[[]string](dbresponse.Answers)
	if err != nil {
		fmt.Printf("Failed to parse data %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error from server"})
		return
	}

	//Compare answers and user's answers
	if len(answers) != len(ReceivedData.UserAnswers) {
		fmt.Printf("Answers are missing \n [length of answers(db):%v] \n [length of answers(users):%v]", len(answers), len(ReceivedData.UserAnswers))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Answers are missing"})
		return
	}
	var marks int
	var totalMarks = len(answers) * 2
	for i := 0; i < len(answers); i++ {
		ans := strings.TrimSpace(answers[i])
		userAns := strings.TrimSpace(ReceivedData.UserAnswers[i])
		//fmt.Printf("Answers by db : %v\n Answers by user : %v\n", answers, ReceivedData.UserAnswers)
		if strings.Contains(ans, userAns) || strings.Contains(userAns, ans) {
			marks += 2
		}
	}
	for v := 0; v < len(ReceivedData.UserAnswers); v++ {
		if ReceivedData.UserAnswers[v] == "" {
			marks -= 2
		}
	}
	fmt.Println(marks)
	c.JSON(http.StatusOK, gin.H{"total marks": totalMarks, "marks": marks, "time": ReceivedData.Time})
}

// this function deletes quiz on basis of quizID given by frontend
func DeleteQuiz(c *gin.Context) {
	//create a db connection
	conn := db.GetDB()
	var quizID string
	//Bind json data to quizID
	err := c.ShouldBindJSON(&quizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No JSON data provided"})
		log.Printf("No JSON data provided %v\n", err)
		return
	}
	err = db.DeleteQuiz(quizID, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quiz"})
		log.Printf("Failed to delete quiz %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"info": "Quiz deleted successfully"})
}
