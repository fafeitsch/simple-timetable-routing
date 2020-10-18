package routing

import (
	"fmt"
	"sort"
	"time"
)

type Timetable struct {
	stops map[string]*vertex
	graph graph
}

func NewTimetable(stops []Stop) Timetable {
	vertices := make([]*vertex, 0, len(stops))
	vertexMap := make(map[string]*vertex)
	for _, stop := range stops {
		vertex := &vertex{data: &stop}
		vertexMap[stop.Id] = vertex
	}
	for _, stop := range stops {
		stop.computeEdges(vertexMap)
	}
	t := Timetable{}
	t.graph = graph{vertices: vertices}
	return Timetable{graph: graph{vertices: vertices}, stops: vertexMap}
}

func (t *Timetable) Query(source *Stop, target *Stop, start time.Time) *Connection {
	s, ok := t.stops[source.Id]
	if !ok {
		panic(fmt.Sprintf("source \"%s\" not found in the timetable", source))
	}
	ta, ok := t.stops[target.Id]
	if !ok {
		panic(fmt.Sprintf("target \"%v\" not found in the timetable", target))
	}
	path := t.graph.shortestPath(s, ta, start)
	return createConnection(path)
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
	Duration time.Duration
	Legs     []Leg
}

func createConnection(path []*vertex) *Connection {
	if len(path) < 2 {
		return nil
	}
	legs := make([]Leg, 0, 0)
	firstStop := path[0]
	currentLine := path[1].currentLine
	remaining := path[1:]
	for i, v := range remaining {
		if v.currentLine != currentLine {
			legs = append(legs, Leg{Line: currentLine, FirstStop: firstStop.data, LastStop: remaining[i-1].data})
			currentLine = v.currentLine
			firstStop = remaining[i-1]
		}
	}
	legs = append(legs, Leg{Line: currentLine, FirstStop: firstStop.data, LastStop: path[len(path)-1].data})
	return &Connection{Legs: legs}
}

type Leg struct {
	Line      *Line
	FirstStop *Stop
	LastStop  *Stop
}
