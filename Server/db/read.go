package db

import (
	utils "Revisor/Utils"
	"database/sql"
	"encoding/json"
	"fmt"
)

// struct to contain data of FlashCard
type Data struct { //This struct will be used as slice
	Heading string `db:"heading"`
	Value   string `db:"value"`
}
type FlashCardData struct {
	Email         string  `db:"email"`
	TopicName     string  `db:"topicName"`
	Time          []uint8 `db:"time"`
	FormattedTime string  //this will contain formatted time
	Data          []Data  `db:"data"`
	Uid           string  `db:"uid"`
}

// this function fetches flashcard data from db
func FetchFlashCardData(connection *sql.DB, email string) ([]FlashCardData, error) {
	query := "select email,topicName,time,data,uid from revisor.flashCardData where email = ?"
	rows, err := connection.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var flashCardData []FlashCardData
	var jsonData []byte //Variable to scan returning data from data field
	for rows.Next() {
		var flashCD FlashCardData
		err := rows.Scan(&flashCD.Email, &flashCD.TopicName, &flashCD.Time, &jsonData, &flashCD.Uid)
		if err != nil {
			fmt.Printf("Failed to scan row %v", err)
			return nil, err
		}
		// Decode JSON into []Data
		err = json.Unmarshal(jsonData, &flashCD.Data)
		if err != nil {
			fmt.Printf("Failed to decode data JSON: %v\n", err)
			return nil, err
		}
		formattedTime, err := utils.Uint8ToTime(flashCD.Time)
		if err != nil {
			fmt.Printf("Failed to convert []uint8 to time.Time : %v", err)
		}
		flashCD.FormattedTime = formattedTime.Format("2006-01-02")
		flashCardData = append(flashCardData, flashCD)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return flashCardData, err
}
