package routing

import (
	"container/heap"
	"fmt"
	"time"
)

type vertex struct {
	data        *StopGroup
	neighbors   []edge
	weight      time.Time
	index       int
	predecessor *vertex
}

func (v *vertex) String() string {
	return fmt.Sprintf("StopGroup-Vertex(\"%s\")", v.data.String())
}

type edge struct {
	weight func(time time.Time) time.Duration
	target *vertex
}

type graph struct {
	vertices []*vertex
}

func (g *graph) shortestPath(s *vertex, t *vertex, start time.Time) []*vertex {
	var sFound, tFound bool
	priorityQueue := &priorityQueue{}
	for _, vertex := range g.vertices {
		tFound = tFound || vertex == t
		if vertex == s {
			vertex.weight = start
			sFound = true
		} else {
			vertex.weight = time.Time{}
		}
		priorityQueue.Push(vertex)
	}
	heap.Init(priorityQueue)
	if !sFound {
		panic(fmt.Sprintf("the provided start vertex %v was not found in the graph", s))
	}
	if !tFound {
		panic(fmt.Sprintf("the provided target vertex %v was not found in the graph", t))
	}
	for len(*priorityQueue) != 0 {
		v := heap.Pop(priorityQueue).(*vertex)
		for _, edge := range v.neighbors {
			neighbour := edge.target
			waitTime := v.weight.Add(edge.weight(v.weight))
			if (neighbour.weight == time.Time{} || waitTime.Before(neighbour.weight)) {
				neighbour.weight = waitTime
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
