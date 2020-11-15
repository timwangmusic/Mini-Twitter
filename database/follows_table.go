package database

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/weihesdlegend/mini_twitter/user"
)

func CreateFollowsTable(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS Follows (id INTEGER PRIMARY KEY, from_user TEXT, to_user TEXT )")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	// create index on from_user
	stmt, err = db.Prepare("CREATE INDEX from_user ON Follows (from_user)")
	if err == nil {
		_, _ = stmt.Exec()
	}

	// create composite index on from_user and to_user for fast unfollow
	stmt, err = db.Prepare("CREATE INDEX following_relationship ON Follows (from_user, to_user)")
	if err == nil {
		_, _ = stmt.Exec()
	}
	return nil // ignore index creation error
}

func CreateFollow(db *sql.DB, follow user.Follow) error {
	stmt, err := db.Prepare("INSERT INTO Follows (from_user, to_user) VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(follow.From, follow.To)
	return err
}

func UnfollowUser(db *sql.DB, unfollow user.Follow) error {
	stmt, err := db.Prepare("DELETE FROM Follows WHERE from_user=? AND to_user=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(unfollow.From, unfollow.To)
	return err
}

func LoadFollows(db *sql.DB, username string, followings map[string]bool) error {
	rows, err := db.Query("SELECT * FROM Follows where from_user=?", username)
	if err != nil {
		return err
	}
	var id int
	var currentUser string
	var followedUser string
	for rows.Next() {
		if err = rows.Scan(&id, &currentUser, &followedUser); err != nil {
			log.Error(err)
		} else {
			followings[followedUser] = true
		}
	}
	return nil
}
