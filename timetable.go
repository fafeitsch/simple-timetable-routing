package routing

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type Time struct {
	Minute int
	Hour   int
}

func (t *Time) interpret(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour, t.Minute, 0, 0, date.Location())
}

var TimeRegex = regexp.MustCompile("^([0-9]{2}):?([0-5][0-9])$")

func ParseTime(string string) Time {
	submatch := TimeRegex.FindStringSubmatch(string)
	if submatch == nil {
		panic(fmt.Sprintf("the string \"%s\" does not match the required format", string))
	}
	hour, _ := strconv.Atoi(submatch[1])
	minute, _ := strconv.Atoi(submatch[2])
	return Time{Hour: hour, Minute: minute}
}

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
	t := Timetable{}
	t.graph = graph{vertices: vertices}
	return Timetable{graph: graph{vertices: vertices}, stops: vertexMap}
}

func (t *Timetable) Query(source *Stop, target *Stop, start time.Time) *Connection {
	for _, stop := range t.stops {
		edges := stop.data.computeEdges(start, t.stops)
		stop.neighbors = edges
	}
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

func (s *Stop) computeEdges(date time.Time, vertices map[string]*vertex) []edge {
	eventGroups := s.groupEvents()
	result := make([]edge, 0, 0)
	for _, event := range eventGroups {
		edge := edge{target: vertices[event[0].NextStop.Name], weight: event.weightFunction(date)}
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

func (e EventGroup) weightFunction(date time.Time) edgeWeight {
	return func(t time.Time, currentLine *Line) (time.Duration, *Line, bool) {
		arrivalMap := make(map[time.Time]Event)
		arrivals := make([]time.Time, 0, len(e))
		for _, event := range e {
			switchTime := 0 * time.Minute
			if event.Line != currentLine {
				switchTime = 5 * time.Minute
			}
			if event.Departure.interpret(date).After(t.Add(switchTime)) {
				arrivalMap[event.ArrivalAtNextStop.interpret(date)] = event
				arrivals = append(arrivals, event.ArrivalAtNextStop.interpret(date))
			}
		}
		if len(arrivals) == 0 {
			return 0 * time.Minute, nil, false
		}
		sort.Slice(arrivals, func(i, j int) bool {
			return arrivals[i].Before(arrivals[j])
		})
		event := arrivalMap[arrivals[0]]
		arrival := event.ArrivalAtNextStop.interpret(date)
		return arrival.Sub(t), event.Line, true
	}
}

type Line struct {
	Id   string
	Name string
}

type Event struct {
	ArrivalAtNextStop Time
	Departure         Time
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
