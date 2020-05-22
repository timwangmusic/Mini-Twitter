package database

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"mini_twitter/user"
)

func CreateUsersTable(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS Users (username TEXT PRIMARY KEY, password TEXT, email TEXT)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("CREATE UNIQUE INDEX username ON Users (username)")
	if err == nil {
		_, err = stmt.Exec()
	}
	return nil  // ignore index creation error
}

func CreateUser(db *sql.DB, user user.User) error {
	stmt, err := db.Prepare("INSERT INTO Users (username, password, email) values (?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Username, user.Password, user.Email)
	return err
}

func LoadUsers(db *sql.DB, users map[string]user.User) error {
	rows, queryErr := db.Query("SELECT * FROM Users")
	if queryErr != nil {
		return queryErr
	}

	var username string
	var password string
	var email string
	for rows.Next() {
		if err := rows.Scan(&username, &password, &email); err != nil {
			log.Error(err)
		} else {
			users[username] = user.User{
				Username: username,
				Password: password,
				Email:    email,
			}
		}

	}

	defer rows.Close()
	return nil
}
