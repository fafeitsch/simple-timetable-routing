package routing

import (
	"fmt"
	"time"
)

type Timetable struct {
	Stops []Stop
}

func (t *Timetable) buildGraph() graph {
	vertices := make([]*vertex, 0, len(t.Stops))
	vertexMap := make(map[string]*vertex)
	for _, stop := range t.Stops {
		vertex := &vertex{data: &stop}
		vertexMap[stop.Id] = vertex

	}
	return graph{vertices: vertices}
}

type Stop struct {
	Id     string
	Name   string
	Events []Event
}

func (s *Stop) String() string {
	return fmt.Sprintf("Stop[Id=\"%s\", Name=\"%s\"]", s.Id, s.Name)
}

func (s *Stop) computeEdges(vertices map[string]*vertex) []edge {
	eventGroups := s.groupEvents()
	result := make([]edge, 0, 0)
	for _, event := range eventGroups {
		edge := edge{target: vertices[event[0].NextStop.Name]}
		edge = edge
	}
	return result
}

func (s *Stop) groupEvents() map[string][]Event {
	result := make(map[string][]Event)
	for _, event := range s.Events {
		if event.NextStop == nil || event.Departure == nil {
			continue
		}
		list, ok := result[event.NextStop.Id]
		if !ok {
			list = make([]Event, 0, 0)
		}
		list = append(list, event)
	}
	return result
}

type Line struct {
	Id   string
	Name string
}

type Event struct {
	Arrival    *time.Time
	Departure  *time.Time
	Line       *Line
	NextStop   *Stop
	TravelTime *time.Duration
}

type Connection struct {
	duration time.Duration
}
