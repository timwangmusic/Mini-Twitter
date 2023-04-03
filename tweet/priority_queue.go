package tweet

type PriorityQueue []Tweet

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].CreatedAt.After(pq[j].CreatedAt)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(item interface{}) {
	*pq = append(*pq, item.(Tweet))
}

func (pq *PriorityQueue) Pop() (item interface{}) {
	prev := *pq
	n := len(prev)
	res := prev[n-1]
	*pq = prev[:n-1]
	return res
}
