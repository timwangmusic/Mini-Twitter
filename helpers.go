package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/weihesdlegend/mini_twitter/tweet"
	"github.com/weihesdlegend/mini_twitter/user"
	"time"
)

const (
	UserNotExistErrorMsg = "user %s does not exist"
)

func postTweet(username string, tweetText string) (error, *tweet.Tweet) {
	t := tweet.Tweet{
		User:      username,
		Text:      tweetText,
		CreatedAt: time.Now().UTC(),
	}
	tweets[username].Tweets = append(tweets[username].Tweets, t)
	return nil, &t
}

func follow(followRequest user.Follow) error {
	fromUser, fromUserExist := users[followRequest.From]
	if !fromUserExist {
		return errors.New(fmt.Sprintf(UserNotExistErrorMsg, followRequest.From))
	}
	toUser, toUserExist := users[followRequest.To]
	if !toUserExist {
		return errors.New(fmt.Sprintf(UserNotExistErrorMsg, followRequest.To))
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
		return errors.New(fmt.Sprintf(UserNotExistErrorMsg, unfollowRequest.From))
	}
	toUser, toUserExist := users[unfollowRequest.To]
	if !toUserExist {
		return errors.New(fmt.Sprintf(UserNotExistErrorMsg, unfollowRequest.To))
	}
	usersFollowing, ok := following[fromUser.Username]
	if !ok {
		return errors.New("user is not following any user")
	}
	_, userFollowingTargetUser := usersFollowing[toUser.Username]
	if !userFollowingTargetUser {
		return errors.New("user is not following the target user")
	}
	delete(usersFollowing, toUser.Username)
	return nil
}
