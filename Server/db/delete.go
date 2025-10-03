package db

import "database/sql"

//Delete quiz from database
func DeleteQuiz(quizId string, connection *sql.DB) error {
	query := "delete from revisor.quiz where quizId = ?"
	_, err := connection.Exec(query, quizId)
	if err != nil {
		return err
	}
	return nil
}
