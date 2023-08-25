package controller

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weihesdlegend/mini_twitter/database"
	"github.com/weihesdlegend/mini_twitter/tweet"
	"github.com/weihesdlegend/mini_twitter/user"
	"github.com/weihesdlegend/mini_twitter/util"
	"time"
)

const (
	UserNotExistErrorMsg = "user %s does not exist"
)

func postTweet(username string, tweetText string) (error, *tweet.Tweet) {
	t := &tweet.Tweet{
		ID:        util.GenID(),
		User:      username,
		Text:      tweetText,
		CreatedAt: time.Now().UTC(),
	}
	database.Tweets[username].Tweets[t.ID] = t
	return nil, t
}

func unfollow(unfollowRequest user.Follow) error {
	fromUser, fromUserExist := database.Users[unfollowRequest.From]
	if !fromUserExist {
		return fmt.Errorf(UserNotExistErrorMsg, unfollowRequest.From)
	}
	toUser, toUserExist := database.Users[unfollowRequest.To]
	if !toUserExist {
		return fmt.Errorf(UserNotExistErrorMsg, unfollowRequest.To)
	}
	usersFollowing, ok := database.Following[fromUser.Username]
	if !ok {
		return fmt.Errorf("user %s is not Following any other user", fromUser.Username)
	}
	_, userFollowingTargetUser := usersFollowing[toUser.Username]
	if !userFollowingTargetUser {
		return fmt.Errorf("user %s is not Following target user %s", fromUser.Username, toUser.Username)
	}
	delete(usersFollowing, toUser.Username)
	return nil
}

func follow(followRequest user.Follow) error {
	fromUser, fromUserExist := database.Users[followRequest.From]
	if !fromUserExist {
		return fmt.Errorf(UserNotExistErrorMsg, followRequest.From)
	}
	toUser, toUserExist := database.Users[followRequest.To]
	if !toUserExist {
		return fmt.Errorf(UserNotExistErrorMsg, followRequest.To)
	}
	logrus.Info("both Users exist")
	if _, ok := database.Following[fromUser.Username]; !ok {
		database.Following[fromUser.Username] = make(map[string]bool)
	}
	database.Following[fromUser.Username][toUser.Username] = true
	return nil
}
