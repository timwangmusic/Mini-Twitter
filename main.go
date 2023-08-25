package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	"github.com/weihesdlegend/mini_twitter/controller"
	"github.com/weihesdlegend/mini_twitter/database"
	_ "github.com/weihesdlegend/mini_twitter/docs"
	"github.com/weihesdlegend/mini_twitter/tweet"
	"github.com/weihesdlegend/mini_twitter/util"
	"net/http"
	"os"

	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	UserDoesNotExist = "user %s does not exist"
)

// @title User API documentation
// @version 1.0.0

// @host localhost:8800
// @BasePath /
func main() {
	router := gin.Default()

	db, dbSetupErr := database.SetupDatabase()
	util.CheckFatal(dbSetupErr)
	database.DB = db

	// create tables if not already exist
	userTableCreationErr := database.CreateUsersTable(db)
	util.CheckErr(userTableCreationErr)

	tweetsTableCreationErr := database.CreateTweetsTable(db)
	util.CheckErr(tweetsTableCreationErr)

	followsTableCreationErr := database.CreateFollowsTable(db)
	util.CheckErr(followsTableCreationErr)

	util.CheckErr(database.LoadUsers(db, database.Users, database.Tweets, database.Following))

	log.Info("starting server")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userRouterGroup := router.Group("/users")
	{
		userRouterGroup.POST("/follow", controller.FollowUser)
		userRouterGroup.POST("/unfollow", controller.UnfollowUser)
		userRouterGroup.POST("/create", controller.CreateUser)
		userRouterGroup.GET("/names", controller.GetAllUsernames)
		userRouterGroup.GET("/following", controller.GetUserFollowing)
	}

	// create new tweet post
	router.POST("/tweets", controller.CreateTweet)

	router.DELETE("/tweets/:id", func(c *gin.Context) {
		id := c.Param("id")
		u := c.Query("user")
		if u == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user name cannot be empty"})
		}
		if _, ok := database.Tweets[u].Tweets[id]; !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "only tweet owner can delete the tweet"})
		}
		if err := database.DeleteTweet(db, id, database.Tweets, u); err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while deleting tweet: " + err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	// get all the tweets from a specific user in reversed order of post creation
	router.GET("/tweets/:username", func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user name not specified"})
		} else if _, ok := database.Users[username]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf(UserDoesNotExist, username)})
		} else {
			ts := make([]*tweet.Tweet, 0)
			for _, t := range database.Tweets[username].Tweets {
				ts = append(ts, t)
			}
			tweet.By(tweet.SortByCreationTime).Sort(ts)
			c.JSON(http.StatusOK, gin.H{"tweets": ts})
		}
	})
	// get timeline for a specific user in reversed order of post creation
	router.GET("/timeline/:username", func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user name not specified"})
		} else if _, ok := database.Users[username]; !ok {
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
