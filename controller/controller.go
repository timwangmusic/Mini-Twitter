package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weihesdlegend/mini_twitter/database"
	"github.com/weihesdlegend/mini_twitter/tweet"
	"github.com/weihesdlegend/mini_twitter/user"
	"github.com/weihesdlegend/mini_twitter/util"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"slices"
)

// CreateUser ... Create a new user
// @Summary Create a new user
// @Description Create a new user
// @Tags User
// @Accept json
// @Param user body user.User true "User Data"
// @Success 200 {array} string
// @Failure 400,303
// @Router /users/create [post]
func CreateUser(c *gin.Context) {
	var newUser user.User
	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if _, ok := database.Users[newUser.Username]; ok {
		c.JSON(http.StatusSeeOther, gin.H{"error": "user already exists"})
	} else {
		encodedPsw, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		newUser.Password = string(encodedPsw)
		database.Users[newUser.Username] = newUser // make sure keys and usernames in records are the same

		// create an empty slice when creating user so after an API verifies the user exists
		// it does not to further check the Tweets table
		database.Tweets[newUser.Username] = &tweet.UserTweets{Tweets: make(map[string]*tweet.Tweet)}

		newUser.Level = user.RegularUser
		if slices.Contains(database.Admins, newUser.Username) {
			newUser.Level = user.AdminUser
		}
		// persist user in database
		util.CheckErr(database.CreateUser(database.DB, newUser))

		c.JSON(http.StatusCreated, gin.H{"user is created": newUser.Username})
	}
}

// GetAllUsernames ... Get all the known usernames
// @Summary Get all the known usernames
// @Description Get all the known usernames
// @Tags User
// @Success 200 {array} string
// @Router /users/names [get]
func GetAllUsernames(c *gin.Context) {
	var usernames []string
	for name := range database.Users {
		usernames = append(usernames, name)
	}
	c.JSON(http.StatusOK, gin.H{"users": usernames})
}

// UnfollowUser ... Unfollow a user
// @Summary Unfollow a user
// @Description Unfollow a user
// @Tags User
// @Accept json
// @Param follow body user.Follow true "Follow Data"
// @Success 200 {array} string
// @Failure 400,404
// @Router /users/unfollow [post]
func UnfollowUser(c *gin.Context) {
	var f user.Follow
	err := c.BindJSON(&f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if err = unfollow(f); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		util.CheckErr(database.UnfollowUser(database.DB, f))
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s unfollowed %s", f.From, f.To)})
	}
}

// FollowUser ... Follow a user
// @Summary Follow a user
// @Description Follow a user
// @Tags User
// @Accept json
// @Param follow body user.Follow true "Follow Data"
// @Success 200 {array} string
// @Failure 400,404
// @Router /users/follow [post]
func FollowUser(c *gin.Context) {
	var f user.Follow
	err := c.BindJSON(&f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if err = follow(f); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		util.CheckErr(database.CreateFollow(database.DB, f))
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%s is following %s", f.From, f.To)})
	}
}

// GetUserFollowing ... List the Users a user is Following
// @Summary List all Users a user is Following
// @Description List all Users a user is Following
// @Tags User
// @Success 200 {array} user.User
// @Failure 404 {object} object
// @Router /users/following [get]
func GetUserFollowing(c *gin.Context) {
	username := c.Query("user")
	if _, exists := database.Users[username]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "user does not exist"})
	}
	results := make([]string, 0)
	for u := range database.Following[username] {
		results = append(results, u)
	}
	c.JSON(http.StatusOK, gin.H{"following": results})
}

// CreateTweet ... Create a new tweet
// @Summary Create a new tweet for a user
// @Description Create a new tweet for a user
// @Tags Tweet
// @Accept json
// @Param follow body tweet.Tweet true "Tweet Data"
// @Success 200 {array} string
// @Failure 400,404,500
// @Router /tweets [post]
func CreateTweet(c *gin.Context) {
	var newPost tweet.Tweet
	err := c.BindJSON(&newPost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if _, userExists := database.Users[newPost.User]; !userExists {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf(UserNotExistErrorMsg, newPost.User)})
	} else if err, newTweet := postTweet(newPost.User, newPost.Text); err == nil {
		err = database.CreateTweet(database.DB, newTweet, database.Tweets)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": "Tweet post success!"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
}
