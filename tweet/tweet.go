package tweet

import (
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
