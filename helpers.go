package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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
	tweets[username].Tweets[t.ID] = t
	return nil, t
}

func follow(followRequest user.Follow) error {
	fromUser, fromUserExist := users[followRequest.From]
	if !fromUserExist {
		return fmt.Errorf(UserNotExistErrorMsg, followRequest.From)
	}
	toUser, toUserExist := users[followRequest.To]
	if !toUserExist {
		return fmt.Errorf(UserNotExistErrorMsg, followRequest.To)
	}
	log.Info("both users exist")
	if _, ok := following[fromUser.Username]; !ok {
		following[fromUser.Username] = make(map[string]bool)
	}
	following[fromUser.Username][toUser.Username] = true
	return nil
}

func unfollow(unfollowRequest user.Follow) error {
	fromUser, fromUserExist := users[unfollowRequest.From]
	if !fromUserExist {
		return fmt.Errorf(UserNotExistErrorMsg, unfollowRequest.From)
	}
	toUser, toUserExist := users[unfollowRequest.To]
	if !toUserExist {
		return fmt.Errorf(UserNotExistErrorMsg, unfollowRequest.To)
	}
	usersFollowing, ok := following[fromUser.Username]
	if !ok {
		return fmt.Errorf("user %s is not following any other user", fromUser.Username)
	}
	_, userFollowingTargetUser := usersFollowing[toUser.Username]
	if !userFollowingTargetUser {
		return fmt.Errorf("user %s is not following target user %s", fromUser.Username, toUser.Username)
	}
	delete(usersFollowing, toUser.Username)
	return nil
}
