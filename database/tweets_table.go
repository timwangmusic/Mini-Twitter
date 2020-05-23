package database

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"mini_twitter/tweet"
	"time"
)

func CreateTweet(db *sql.DB, tweet tweet.Tweet) error {
	stmt, err := db.Prepare("INSERT INTO Tweets (text, username, created_at) values (?,?,?)")
	if err != nil {
		return err
	}
	logrus.Info("tweet creation time is:", tweet.CreatedAt.Format(time.RFC3339))
	_, err = stmt.Exec(tweet.Text, tweet.User, tweet.CreatedAt.Format(time.RFC3339))
	return err
}

func LoadTweets(db *sql.DB, username string, userTweets *tweet.UserTweets) error {
	rows, err := db.Query("SELECT * FROM Tweets WHERE username=?", username)
	if err != nil {
		return err
	}

	var id int
	var text string
	var u string
	var createdAt time.Time
	for rows.Next() {
		err = rows.Scan(&id, &text, &u, &createdAt)
		if err != nil {
			logrus.Error(err)
		}
		userTweets.Tweets = append(userTweets.Tweets, tweet.Tweet{
			Text:      text,
			User:      u,
			CreatedAt: createdAt,
		})
	}
	_ = rows.Close()
	return nil
}

func CreateTweetsTable(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS Tweets (id INTEGER PRIMARY KEY," +
		"text TEXT, username TEXT, created_at datetime)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = db.Prepare("CREATE INDEX username_in_tweets ON Tweets (username)")
	if err == nil {
		_, _ = stmt.Exec()
	}
	return nil // ignore index creation error
}
