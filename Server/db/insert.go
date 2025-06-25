package db

import (
	"database/sql"
	"encoding/json"
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
func InsertFlashCardData(connection *sql.DB, email string, topicName string, data []map[string]string) error {
	query := "insert into revisor.flashCardData(email,topicName,data) values(?,?,?)"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = connection.Exec(query, email, topicName, jsonData)
	if err != nil {
		return err
	}
	return nil
}
