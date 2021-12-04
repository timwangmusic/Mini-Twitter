package tweet

type Sorter struct {
	tweets *UserTweets
	by     func(t1, t2 *Tweet) bool
}

func (s *Sorter) Len() int {
	return len(s.tweets.Tweets)
}

func (s *Sorter) Swap(i, j int) {
	s.tweets.Tweets[i], s.tweets.Tweets[j] = s.tweets.Tweets[j], s.tweets.Tweets[i]
}

func (s *Sorter) Less(i, j int) bool {
	return s.by(&s.tweets.Tweets[i], &s.tweets.Tweets[j])
}
