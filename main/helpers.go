package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"mini_tweeter/tweet"
	"mini_tweeter/user"
	"time"
)

const (
	UserNotExistErrorMsg = "user %s does not exist"
)

func postTweet(username string, tweetText string) error {
	if _, ok := users[username]; !ok {
		return errors.New(fmt.Sprintf(UserNotExistErrorMsg, username))
	}
	t := tweet.Tweet{
		User:      username,
		Text:      tweetText,
		CreatedAt: time.Now(),
	}
	if _, ok := tweets[username]; !ok {
		tweets[username] = new(tweet.UserTweets)
		tweets[username].Tweets = make([]tweet.Tweet, 0)
	}
	tweets[username].Tweets = append(tweets[username].Tweets, t)
	return nil
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
