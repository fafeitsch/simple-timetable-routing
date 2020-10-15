package routing

import (
	"container/heap"
	"fmt"
	"time"
)

type vertex struct {
	currentLine *Line
	data        *Stop
	neighbors   []edge
	weight      time.Time
	index       int
	predecessor *vertex
}

func (v *vertex) String() string {
	return fmt.Sprintf("Vertex[%s]", v.data.String())
}

type edgeWeight func(time time.Time, currentLine *Line) (time.Duration, *Line, bool)

type edge struct {
	weight edgeWeight
	target *vertex
}

type graph struct {
	vertices []*vertex
}

func (g *graph) shortestPath(s *vertex, t *vertex, start time.Time) []*vertex {
	priorityQueue := &priorityQueue{}
	for _, vertex := range g.vertices {
		vertex.weight = time.Time{}
		priorityQueue.Push(vertex)
	}
	s.weight = start
	heap.Init(priorityQueue)
	for len(*priorityQueue) != 0 {
		v := heap.Pop(priorityQueue).(*vertex)
		for _, edge := range v.neighbors {
			neighbour := edge.target
			weight, line, ok := edge.weight(v.weight, v.currentLine)
			if !ok {
				// there is no suitable departure to that neighbour any more, skip it.
				continue
			}
			waitTime := v.weight.Add(weight)
			if (neighbour.weight == time.Time{} || waitTime.Before(neighbour.weight)) {
				neighbour.weight = waitTime
				neighbour.currentLine = line
				neighbour.predecessor = v
				priorityQueue.update(neighbour)
			}
		}
	}
	result := make([]*vertex, 0, 0)
	predecessor := t
	for predecessor != nil {
		result = append(result, predecessor)
		predecessor = predecessor.predecessor
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
