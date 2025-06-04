package db

import "database/sql"

//this function insert user in db
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
