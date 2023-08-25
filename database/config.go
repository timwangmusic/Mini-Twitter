package database

import (
	"database/sql"
	"github.com/weihesdlegend/mini_twitter/tweet"
	"github.com/weihesdlegend/mini_twitter/user"
	"github.com/weihesdlegend/mini_twitter/util"
	"os"
)

const DefaultDatabaseName = "mini-twitter.DB"

var DB *sql.DB
var Following map[string]map[string]bool // user to users the user is Following
var Users map[string]user.User           // user details
var Tweets map[string]*tweet.UserTweets  // user to Tweets

func SetupDatabase() (*sql.DB, error) {
	Tweets = make(map[string]*tweet.UserTweets)
	Users = make(map[string]user.User)
	Following = make(map[string]map[string]bool)

	_, err := os.Stat(DefaultDatabaseName)
	if os.IsNotExist(err) {
		_, creationErr := os.Create(DefaultDatabaseName)
		return nil, creationErr
	}

	db, dbConnectionErr := sql.Open("sqlite3", DefaultDatabaseName)
	if dbConnectionErr != nil {
		return nil, dbConnectionErr
	}

	// create tables if not already exist
	userTableCreationErr := CreateUsersTable(db)
	util.CheckErr(userTableCreationErr)

	tweetsTableCreationErr := CreateTweetsTable(db)
	util.CheckErr(tweetsTableCreationErr)

	followsTableCreationErr := CreateFollowsTable(db)
	util.CheckErr(followsTableCreationErr)

	util.CheckErr(LoadUsers(db, Users, Tweets, Following))

	return db, nil
}
