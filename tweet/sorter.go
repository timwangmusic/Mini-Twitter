package tweet

type Sorter struct {
	tweets []*Tweet
	by     func(t1, t2 *Tweet) bool
}

func (s *Sorter) Len() int {
	return len(s.tweets)
}

func (s *Sorter) Swap(i, j int) {
	s.tweets[i], s.tweets[j] = s.tweets[j], s.tweets[i]
}

func (s *Sorter) Less(i, j int) bool {
	return s.by(s.tweets[i], s.tweets[j])
}
