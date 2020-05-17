package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"mini_tweeter/tweet"
	"mini_tweeter/user"
	"net/http"
)

var tweets map[string]*tweet.UserTweets  // user to tweets
var users map[string]user.User           // user details
var following map[string]map[string]bool // user to users the user is following

func main() {
	tweets = make(map[string]*tweet.UserTweets)
	users = make(map[string]user.User)
	following = make(map[string]map[string]bool)
	router := gin.Default()

	log.Println("starting server")
	// create new user
	router.POST("/users", func(c *gin.Context) {
		var newUser user.User
		err := c.BindJSON(&newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			encodedPsw, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
			newUser.Password = string(encodedPsw)
			users[newUser.Username] = newUser // make sure keys and usernames in records are the same
			c.JSON(http.StatusOK, gin.H{})
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
			c.JSON(http.StatusOK, fmt.Sprintf("%s is following %s", f.From, f.To))
		}
	})

	// create new tweet post
	router.POST("/tweets", func(c *gin.Context) {
		var newPost tweet.Tweet
		err := c.BindJSON(&newPost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err = postTweet(newPost.User, newPost.Text); err == nil {
			c.JSON(http.StatusOK, gin.H{"result": "Tweet post success!"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
	})

	// get all the tweets from a specific user in reversed order of creation
	router.GET("/tweets", func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user name not specified"})
		} else if _, ok := users[username]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "user does not exist"})
		} else {
			if _, ok := tweets[username]; !ok {
				tweets[username] = &tweet.UserTweets{Tweets: make([]tweet.Tweet, 0)}
			}
			tweet.By(tweet.SortByCreationTime).Sort(tweets[username])
			c.JSON(http.StatusOK, gin.H{"result": tweets[username]})
		}

	})

	router.GET("/timeline", func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user name not specified"})
		} else if _, ok := users[username]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "user does not exist"})
		} else {
			timeline := GetTimeLine(username)
			c.JSON(http.StatusOK, gin.H{"result": timeline})
		}
	})

	if err := router.Run(":8800"); err != nil {
		log.Fatal(err)
	}

}
