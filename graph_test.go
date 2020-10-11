package routing

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGraph_shortestPath(t *testing.T) {
	a := &vertex{}
	b := &vertex{}
	c := &vertex{}
	d := &vertex{}
	e := &vertex{}
	f := &vertex{}
	g := &vertex{}
	a.data = &StopGroup{Stops: []Stop{{Name: "A"}}}
	a.neighbors = []edge{
		{target: f, weight: constantWeight(100)},
		{target: b, weight: constantWeight(10)},
	}
	b.data = &StopGroup{Stops: []Stop{{Name: "B"}}}
	b.neighbors = []edge{
		{target: e, weight: constantWeight(30)},
		{target: d, weight: constantWeight(10)},
	}
	c.data = &StopGroup{Stops: []Stop{{Name: "C"}}}
	c.neighbors = []edge{
		{target: g, weight: constantWeight(40)},
	}
	d.data = &StopGroup{Stops: []Stop{{Name: "D"}}}
	d.neighbors = []edge{
		{target: c, weight: constantWeight(5)},
		{target: f, weight: constantWeight(45)},
		{target: e, weight: constantWeight(10)},
	}
	e.data = &StopGroup{Stops: []Stop{{Name: "E"}}}
	e.neighbors = []edge{
		{target: f, weight: constantWeight(10)},
	}
	f.data = &StopGroup{Stops: []Stop{{Name: "F"}}}
	f.neighbors = []edge{
		{target: c, weight: constantWeight(40)},
		{target: b, weight: constantWeight(25)},
		{target: d, weight: constantWeight(80)},
	}
	g.data = &StopGroup{Stops: []Stop{{Name: "G"}}}
	g.neighbors = []edge{
		{target: f, weight: constantWeight(20)},
	}
	graph := graph{vertices: []*vertex{a, b, c, d, e, f, g}}
	t.Run("success", func(t *testing.T) {
		start, _ := time.Parse(time.RFC3339, "2020-10-11T18:00:00Z")
		path := graph.shortestPath(a, f, start)
		assert.Equal(t, []*vertex{a, b, d, e, f}, path, "path not computed correctly")
		assert.Equal(t, "2020-10-11T18:40:00Z", f.weight.Format(time.RFC3339), "arrival time not computed correctly")
	})
	t.Run("panic s is missing", func(t *testing.T) {
		assert.PanicsWithValue(t, "the provided start vertex StopGroup-Vertex(\"A\") was not found in the graph", func() { graph.shortestPath(&vertex{data: a.data}, f, time.Now()) }, "panic message wrong")
	})
	t.Run("panic t is missing", func(t *testing.T) {
		assert.PanicsWithValue(t, "the provided target vertex StopGroup-Vertex(\"F\") was not found in the graph", func() { graph.shortestPath(a, &vertex{data: f.data}, time.Now()) }, "panic message wrong")
	})
}

func constantWeight(weight int) func(time.Time) time.Duration {
	return func(t time.Time) time.Duration {
		return time.Duration(weight) * time.Minute
	}
}
