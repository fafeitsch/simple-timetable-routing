package routing

import (
	"fmt"
	"sort"
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
	for _, stop := range t.Stops {
		stop.computeEdges(vertexMap)
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
		edge := edge{target: vertices[event[0].NextStop.Name], weight: event.weightFunction()}
		result = append(result, edge)
	}
	return result
}

func (s *Stop) groupEvents() map[string]EventGroup {
	result := make(map[string]EventGroup)
	for _, event := range s.Events {
		list, ok := result[event.NextStop.Id]
		if !ok {
			list = make([]Event, 0, 0)
		}
		list = append(list, event)
		result[event.NextStop.Id] = list
	}
	return result
}

type EventGroup []Event

func (e EventGroup) weightFunction() edgeWeight {
	return func(t time.Time, currentLine *Line) (time.Duration, *Line, bool) {
		arrivalMap := make(map[time.Time]Event)
		arrivals := make([]time.Time, 0, len(e))
		for _, event := range e {
			switchTime := 0 * time.Minute
			if event.Line != currentLine {
				switchTime = 5 * time.Minute
			}
			if event.Departure.After(t.Add(switchTime)) {
				arrivalMap[event.ArrivalAtNextStop] = event
				arrivals = append(arrivals, event.ArrivalAtNextStop)
			}
		}
		if len(arrivals) == 0 {
			return 0 * time.Minute, nil, false
		}
		sort.Slice(arrivals, func(i, j int) bool {
			return arrivals[i].Before(arrivals[j])
		})
		event := arrivalMap[arrivals[0]]
		return event.ArrivalAtNextStop.Sub(t), event.Line, true
	}
}

type Line struct {
	Id   string
	Name string
}

type Event struct {
	ArrivalAtNextStop time.Time
	Departure         time.Time
	Line              *Line
	NextStop          *Stop
}

type Connection struct {
	duration time.Duration
}
