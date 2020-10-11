package routing

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPriorityQueue_Len(t *testing.T) {
	pq := priorityQueue{}
	assert.Equal(t, 0, pq.Len(), "Initial Length")
	vertex1 := &vertex{weight: time.Now()}
	vertex2 := &vertex{weight: time.Now()}
	vertex3 := &vertex{weight: time.Now()}
	pq.Push(vertex1)
	pq.Push(vertex2)
	pq.Push(vertex3)
	assert.Equal(t, 3, pq.Len(), "length after adding")
}

func TestPriorityQueue_Pop(t *testing.T) {
	now := time.Now()
	vertex1 := &vertex{weight: now.Add(7 * time.Minute)}
	vertex2 := &vertex{weight: now.Add(5 * time.Minute)}
	vertex3 := &vertex{weight: now.Add(-1 * time.Minute)}
	vertex4 := &vertex{weight: now}
	queue := priorityQueue{}
	heap.Push(&queue, vertex1)
	heap.Push(&queue, vertex2)
	heap.Push(&queue, vertex3)
	heap.Push(&queue, vertex4)

	result := make([]*vertex, 0, 4)
	for queue.Len() != 0 {
		result = append(result, heap.Pop(&queue).(*vertex))
	}
	assert.Equal(t, []*vertex{vertex3, vertex4, vertex2, vertex1}, result, "pop order is not correct")
}

func TestPriorityQueue_update(t *testing.T) {
	now := time.Now()
	vertex1 := &vertex{weight: now.Add(7 * time.Minute)}
	vertex2 := &vertex{weight: now.Add(5 * time.Minute)}
	vertex3 := &vertex{weight: now.Add(-1 * time.Minute)}
	vertex4 := &vertex{weight: now}
	queue := priorityQueue{}
	heap.Push(&queue, vertex1)
	heap.Push(&queue, vertex2)
	heap.Push(&queue, vertex3)
	heap.Push(&queue, vertex4)
	vertex2.weight = now.Add(-3 * time.Hour)
	queue.update(vertex2)
	result := make([]*vertex, 0, 4)
	for queue.Len() != 0 {
		result = append(result, heap.Pop(&queue).(*vertex))
	}
	assert.Equal(t, []*vertex{vertex2, vertex3, vertex4, vertex1}, result, "pop order is not correct")
}
