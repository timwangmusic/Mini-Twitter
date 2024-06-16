package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/weihesdlegend/mini_twitter/controller"
	"github.com/weihesdlegend/mini_twitter/database"
	_ "github.com/weihesdlegend/mini_twitter/docs"
	"github.com/weihesdlegend/mini_twitter/tweet"
	"github.com/weihesdlegend/mini_twitter/util"
	"net/http"
	"os"
	"path"
	"sort"
)

const (
	UserDoesNotExist = "user %s does not exist"

	DefaultConfigPath     = "."
	DefaultConfigFileName = "config.yaml"
)

type Config struct {
	Admin []string `mapstructure:"admin"`
}

// @title User API documentation
// @version 1.0.0

// @host localhost:8800
// @BasePath /
func main() {
	router := gin.Default()
	if _, err := os.Stat(path.Join(DefaultConfigPath, DefaultConfigFileName)); err != nil {
		log.Infof("Creating default config file at %s", path.Join(DefaultConfigPath, DefaultConfigFileName))
		if os.IsNotExist(err) {
			_, err = os.Create(path.Join(DefaultConfigPath, DefaultConfigFileName))
			if err != nil {
				log.Fatalf("failed to create config file: %s", DefaultConfigFileName)
			}
		}
	}

	viper.AddConfigPath(DefaultConfigPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName(DefaultConfigFileName)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}

	var config Config
	if err = viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to unmarshal into config: %s", err.Error())
	}

	log.Debugf("admin users are: %+v", config.Admin)

	db, dbSetupErr := database.SetupDatabase(config.Admin)
	util.CheckFatal(dbSetupErr)
	database.DB = db

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
			sort.Slice(ts, func(i, j int) bool { return ts[i].CreatedAt.After(ts[j].CreatedAt) })
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
