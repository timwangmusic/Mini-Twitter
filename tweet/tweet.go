package tweet

import (
	"sort"
	"time"
)

type Tweet struct {
	Text      string    `json:"text"`
	User      string    `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

// UserTweets contains all the tweets from an user
type UserTweets struct {
	Tweets []Tweet `json:"tweets"`
}

var SortByCreationTime = func(t1, t2 *Tweet) bool {
	return t1.CreatedAt.Sub(t2.CreatedAt) > 0
}

type By func(p1, p2 *Tweet) bool

func (by By) Sort(tweets *UserTweets) {
	sorter := &Sorter{
		tweets: tweets,
		by:     by,
	}
	sort.Sort(sorter)
}
