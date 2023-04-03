package tweet

import (
	"sort"
	"time"
)

type Tweet struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// UserTweets contains all the tweets from a user
type UserTweets struct {
	Tweets map[string]*Tweet `json:"tweets"`
}

var SortByCreationTime = func(t1, t2 *Tweet) bool {
	return t1.CreatedAt.Sub(t2.CreatedAt) > 0
}

type By func(p1, p2 *Tweet) bool

func (by By) Sort(tweets []*Tweet) {
	sorter := &Sorter{
		tweets: tweets,
		by:     by,
	}
	sort.Sort(sorter)
}
