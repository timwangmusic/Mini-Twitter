package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/weihesdlegend/mini_twitter/database"
	"github.com/weihesdlegend/mini_twitter/tweet"
	"github.com/weihesdlegend/mini_twitter/user"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

const (
	UserDoesNotExist = "user %s does not exist"
)

var tweets map[string]*tweet.UserTweets  // user to tweets
var users map[string]user.User           // user details
var following map[string]map[string]bool // user to users the user is following

func main() {
	tweets = make(map[string]*tweet.UserTweets)
	users = make(map[string]user.User)
	following = make(map[string]map[string]bool)
	router := gin.Default()

	dbName := "mini-twitter.db"
	_, err := os.Stat(dbName)
	if os.IsNotExist(err) {
		_, creationErr := os.Create(dbName)
		checkFatal(creationErr)
	}

	db, dbConnectionErr := sql.Open("sqlite3", dbName)
	checkFatal(dbConnectionErr)

	// create tables if not already exist
	userTableCreationErr := database.CreateUsersTable(db)
	checkErr(userTableCreationErr)

	tweetsTableCreationErr := database.CreateTweetsTable(db)
	checkErr(tweetsTableCreationErr)

	followsTableCreationErr := database.CreateFollowsTable(db)
	checkErr(followsTableCreationErr)

	checkErr(database.LoadUsers(db, users, tweets, following))

	log.Info("starting server")
	// create a new user
	router.POST("/users", func(c *gin.Context) {
		var newUser user.User
		err := c.BindJSON(&newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if _, ok := users[newUser.Username]; ok {
			c.JSON(http.StatusSeeOther, gin.H{"error": "user already exists"})
		} else {
			encodedPsw, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
			newUser.Password = string(encodedPsw)
			users[newUser.Username] = newUser // make sure keys and usernames in records are the same

			// create an empty slice when creating user so after an API verifies the user exists
			// it does not to further check the tweets table
			tweets[newUser.Username] = &tweet.UserTweets{Tweets: make([]tweet.Tweet, 0)}

			// persist user in database
			checkErr(database.CreateUser(db, newUser))

			c.JSON(http.StatusCreated, gin.H{"user created": newUser.Username})
		}
	})

	router.POST("/follows", func(c *gin.Context) {
		var f user.Follow
		err := c.BindJSON(&f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err = follow(f); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			checkErr(database.CreateFollow(db, f))
			c.JSON(http.StatusOK, fmt.Sprintf("%s is following %s", f.From, f.To))
		}
	})

	router.POST("/unfollows", func(c *gin.Context) {
		var f user.Follow
		err := c.BindJSON(&f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err = unfollow(f); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			checkErr(database.UnfollowUser(db, f))
			c.JSON(http.StatusOK, fmt.Sprintf("%s unfollowed %s", f.From, f.To))
		}
	})

	// create new tweet post
	router.POST("/tweets", func(c *gin.Context) {
		var newPost tweet.Tweet
		err := c.BindJSON(&newPost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if _, userExists := users[newPost.User]; !userExists {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf(UserDoesNotExist, newPost.User)})
		} else if err, newTweet := postTweet(newPost.User, newPost.Text); err == nil {
			err = database.CreateTweet(db, newTweet)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"result": "Tweet post success!"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
	})

	// get all the tweets from a specific user in reversed order of post creation
	router.GET("/tweets/:username", func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user name not specified"})
		} else if _, ok := users[username]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf(UserDoesNotExist, username)})
		} else {
			tweet.By(tweet.SortByCreationTime).Sort(tweets[username])
			c.JSON(http.StatusOK, gin.H{"result": tweets[username]})
		}
	})
	// get timeline for a specific user in reversed order of post creation
	router.GET("/timeline/:username", func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user name not specified"})
		} else if _, ok := users[username]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf(UserDoesNotExist, username)})
		} else {
			timeline := GetTimeLine(username)
			c.JSON(http.StatusOK, gin.H{"result": timeline})
		}
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to use the Mini Twitter")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8800"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
