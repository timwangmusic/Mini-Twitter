package main

import (
	"container/heap"
	"mini_twitter/tweet"
)

func GetTimeLine(user string) (timeline []tweet.Tweet) {
	timeline = make([]tweet.Tweet, 0)

	pq := &tweet.PriorityQueue{}

	userTweetIdxMap := make(map[string]int)  // user to tweet ID mapping
	followingUsers, _ := following[user]
	for u := range followingUsers {
		if len(tweets[u].Tweets) > 0 {
			tweet.By(tweet.SortByCreationTime).Sort(tweets[u])
			userTweetIdxMap[u] = 0
			heap.Push(pq, tweets[u].Tweets[0])
		}
	}

	// K路归并
	for len(*pq) > 0 {
		top := heap.Pop(pq).(tweet.Tweet)
		timeline = append(timeline, top)
		curUser := top.User
		userTweetIdxMap[curUser] += 1
		nextTweetIdx := userTweetIdxMap[curUser]
		if userTweetIdxMap[curUser] < len(tweets[curUser].Tweets) {
			heap.Push(pq, tweets[curUser].Tweets[nextTweetIdx])
		}
	}

	return
}
