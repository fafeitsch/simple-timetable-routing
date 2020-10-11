package routing

import (
	"container/heap"
	"time"
)

type priorityQueue []*vertex

func (p priorityQueue) Len() int {
	return len(p)
}

func (p priorityQueue) Less(i int, j int) bool {
	empty := time.Time{}
	if p[j].weight == empty {
		return true
	}
	if p[i].weight == empty {
		return false
	}
	return p[i].weight.Before(p[j].weight)
}

func (p priorityQueue) Swap(i int, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

func (p *priorityQueue) Push(x interface{}) {
	n := len(*p)
	item := x.(*vertex)
	item.index = n
	*p = append(*p, item)
}

func (p *priorityQueue) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	item.index = -1
	old[n-1] = nil
	*p = old[0 : n-1]
	return item
}

func (p *priorityQueue) update(vertex *vertex) {
	heap.Fix(p, vertex.index)
}
