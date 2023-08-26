package main

import (
	"container/heap"
	"github.com/weihesdlegend/mini_twitter/database"
	"github.com/weihesdlegend/mini_twitter/tweet"
	"sort"
)

func GetTimeLine(user string) (timeline []tweet.Tweet) {
	timeline = make([]tweet.Tweet, 0)

	pq := &tweet.PriorityQueue{}

	userTweetIdxMap := make(map[string]int) // user to tweet index mapping
	followingUsers, _ := database.Following[user]

	tweetsPerUser := make(map[string][]*tweet.Tweet)
	for u := range followingUsers {
		if _, ok := tweetsPerUser[u]; !ok {
			tweetsPerUser[u] = make([]*tweet.Tweet, 0)
		}
		if len(database.Tweets[u].Tweets) > 0 {
			userTweetIdxMap[u] = 0
			for _, t := range database.Tweets[u].Tweets {
				tweetsPerUser[u] = append(tweetsPerUser[u], t)
			}
			sort.Slice(tweetsPerUser[u], func(i, j int) bool { return tweetsPerUser[u][i].CreatedAt.After(tweetsPerUser[u][j].CreatedAt) })
			heap.Push(pq, *tweetsPerUser[u][0])
		}
	}

	// K路归并
	for len(*pq) > 0 {
		top := heap.Pop(pq).(tweet.Tweet)
		timeline = append(timeline, top)
		curUser := top.User
		userTweetIdxMap[curUser] += 1
		nextTweetIdx := userTweetIdxMap[curUser]
		if userTweetIdxMap[curUser] < len(tweetsPerUser[curUser]) {
			heap.Push(pq, *tweetsPerUser[curUser][nextTweetIdx])
		}
	}

	return
}
