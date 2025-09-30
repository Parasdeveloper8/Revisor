package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
)

// this function inserts user in db
func InsertUser(connection *sql.DB, email string, familyname string, givenName string, id string, name string, picture string, verifiedEmail bool) error {
	query := `insert into revisor.user 
	(email,family_name,given_name,id,name,picture,verified_email)
	values (?,?,?,?,?,?,?)`
	_, err := connection.Exec(query, email, familyname, givenName, id, name, picture, verifiedEmail)
	if err != nil {
		return err
	}
	return nil
}

// this function inserts flashcard data in db
func InsertFlashCardData(connection *sql.DB, email string, topicName string, data []map[string]string, uniqueId string) error {
	query := "insert into revisor.flashCardData(email,topicName,data,uid) values(?,?,?,?)"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = connection.Exec(query, email, topicName, jsonData, uniqueId)
	if err != nil {
		return err
	}
	return nil
}

// This function inserts quiz data in db
func InsertQuizData(connection *sql.DB, email string, answers []string, question []string, options [][]string, quizId string, noteId string) error {
	if len(answers) != len(question) {
		if len(answers) > len(question) {
			err := errors.New("questions are missing")
			return err
		} else if len(question) > len(answers) {
			err := errors.New("answers are missing")
			return err
		}
	}
	if len(options) == 0 {
		err := errors.New("options are missing")
		return err
	}
	query := "insert into revisor.quiz(quizId,question,answers,options,email,noteId) values (?,?,?,?,?,?)"
	//convert data into json --->
	jsonAnswers, err := json.Marshal(answers)
	if err != nil {
		return err
	}
	jsonQuestions, err := json.Marshal(question)
	if err != nil {
		return err
	}
	jsonOptions, err := json.Marshal(options)
	if err != nil {
		return err
	}
	_, err = connection.Exec(query, quizId, jsonQuestions, jsonAnswers, jsonOptions, email, noteId)
	if err != nil {
		return err
	}
	log.Println("Quiz data saved in database")
	return nil
}
