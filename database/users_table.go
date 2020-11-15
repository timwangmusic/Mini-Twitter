package database

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/weihesdlegend/mini_twitter/tweet"
	"github.com/weihesdlegend/mini_twitter/user"
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
		_, _ = stmt.Exec()
	}
	return nil // ignore index creation error
}

func CreateUser(db *sql.DB, user user.User) error {
	stmt, err := db.Prepare("INSERT INTO Users (username, password, email) values (?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Username, user.Password, user.Email)
	return err
}

func LoadUsers(db *sql.DB, users map[string]user.User, tweets map[string]*tweet.UserTweets,
	follows map[string]map[string]bool) error {
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
			tweets[username] = &tweet.UserTweets{
				Tweets: make([]tweet.Tweet, 0),
			}
			if err := LoadTweets(db, username, tweets[username]); err != nil {
				log.Error(err)
			}
			follows[username] = make(map[string]bool)
			if err := LoadFollows(db, username, follows[username]); err != nil {
				log.Error(err)
			}
		}
	}

	_ = rows.Close()
	return nil
}
