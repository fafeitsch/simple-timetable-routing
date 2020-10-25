package routing

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"
)

// Time is a string data type that can be interpreted as simple time
// (without date). The time string should always match the TimeRegex,
// otherwise a panic may be risen.
//
// Examples: 12:04, 14:34, 28:23 are all valid times
type Time string

// CreateTime creates a time string from a given hour and minute.
func CreateTime(hour, minute int) Time {
	return Time(fmt.Sprintf("%02d:%02d", hour, minute))
}

// TimeRegex is used to validate time strings.
var TimeRegex = regexp.MustCompile("^([0-9]+):?([0-5][0-9])$")

func (t Time) interpret(date time.Time) time.Time {
	submatch := TimeRegex.FindStringSubmatch(string(t))
	if submatch == nil {
		panic(fmt.Sprintf("the string \"%s\" does not match the required format", t))
	}
	hour, _ := strconv.Atoi(submatch[1])
	minute, _ := strconv.Atoi(submatch[2])
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())
}

// Timetable contains all routing information in a public transport network.
// Timetables should be created with the NewTimetable function.
type Timetable struct {
	stops map[string]*vertex
	graph graph
}

// NewTimetable creates a new timetable containing the passed stops. The stops
// contain all relevant information about the transport network (arrivals, departures, and lines).
func NewTimetable(stops []*Stop) Timetable {
	vertices := make([]*vertex, 0, len(stops))
	vertexMap := make(map[string]*vertex)
	for _, stop := range stops {
		vertex := &vertex{data: stop}
		vertexMap[stop.Id] = vertex
		vertices = append(vertices, vertex)
	}
	t := Timetable{}
	t.graph = graph{vertices: vertices}
	return Timetable{graph: graph{vertices: vertices}, stops: vertexMap}
}

// Query computes the fastest route between source and target with the specified start time.
// If there is no connection, then nil is returned.
func (t *Timetable) Query(source *Stop, target *Stop, start time.Time) *Connection {
	for _, stop := range t.stops {
		edges := stop.data.computeEdges(start, t.stops)
		stop.neighbors = edges
	}
	s, ok := t.stops[source.Id]
	if !ok {
		panic(fmt.Sprintf("source \"%s\" not found in the timetable", source.Id))
	}
	ta, ok := t.stops[target.Id]
	if !ok {
		panic(fmt.Sprintf("target \"%s\" not found in the timetable", target.Id))
	}
	path := t.graph.shortestPath(s, ta, start)
	return createConnection(path)
}

// Stop is a physical stop where a public transport vehicle stops and lets
// passengers enter and exit. The Id of the stop must be unique.
type Stop struct {
	Id     string
	Name   string
	Events []Event
}

// NewStop creates a new stop with the given id and name and an empty events slice.
func NewStop(id, name string) *Stop {
	return &Stop{Id: id, Name: name, Events: make([]Event, 0, 0)}
}

func (s *Stop) computeEdges(date time.Time, vertices map[string]*vertex) []edge {
	eventGroups := s.groupEvents()
	result := make([]edge, 0, 0)
	for _, event := range eventGroups {
		edge := edge{target: vertices[event[0].nextStop().Id], weight: event.weightFunction(date)}
		result = append(result, edge)
	}
	return result
}

func (s *Stop) groupEvents() map[string]eventGroup {
	result := make(map[string]eventGroup)
	for _, event := range s.Events {
		list, ok := result[event.nextStop().Id]
		if !ok {
			list = make([]Event, 0, 0)
		}
		list = append(list, event)
		result[event.nextStop().Id] = list
	}
	return result
}

type eventGroup []Event

func (e eventGroup) weightFunction(date time.Time) edgeWeight {
	return func(t time.Time, currentLine *Line) (time.Duration, *Line, bool) {
		arrivalMap := make(map[time.Time]Event)
		arrivals := make([]time.Time, 0, len(e))
		for _, event := range e {
			switchTime := 0 * time.Minute
			// if currentLine == nil, we are at the source station
			if currentLine != nil && event.Line != currentLine {
				switchTime = 5 * time.Minute
			}
			departure := event.Departure.interpret(date)
			switchFinished := t.Add(switchTime)
			if departure.Equal(switchFinished) || departure.After(switchFinished) {
				arrival := departure.Add(event.durationToNextStop())
				arrivalMap[arrival] = event
				arrivals = append(arrivals, arrival)
			}
		}
		if len(arrivals) == 0 {
			return 0 * time.Minute, nil, false
		}
		sort.Slice(arrivals, func(i, j int) bool {
			return arrivals[i].Before(arrivals[j])
		})
		event := arrivalMap[arrivals[0]]
		arrival := arrivals[0]
		return arrival.Sub(t), event.Line, true
	}
}

// Line represents a line in a public transportation network. It consists
// of an Id, which should be unique among all lines, as well as stops and segments.
// Variants of lines (e.g. additional stops in the rush hour) must be modeled with
// additional lines.
type Line struct {
	Id        string
	Name      string
	startStop *Stop
	Segments  []Segment
}

// Segment describes the travel time from the last stop in a line to the next stop
// in the line.
type Segment struct {
	TravelTime time.Duration
	NextStop   *Stop
}

// Event describes the departure of a certain line's vehicle at a station. The segment property
// points to the line segment that describes the journey to the next stop after the event's stop.
type Event struct {
	Departure Time
	Line      *Line
	Segment   *Segment
}

func (e *Event) nextStop() *Stop {
	return e.Segment.NextStop
}

func (e *Event) durationToNextStop() time.Duration {
	return e.Segment.TravelTime
}

// Connection is the result of a route computation. It contains the arrival time
// as well a slice of legs which describe the lines that have to be taken at certain stations
// in order to reach the target station.
type Connection struct {
	Arrival time.Time
	Legs    []Leg
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
	return &Connection{Legs: legs, Arrival: path[len(path)-1].weight}
}

// Leg is a part of a journey during which there is no change of lines. A leg
// has the first stop, a last stop and a line.
type Leg struct {
	Line      *Line
	FirstStop *Stop
	LastStop  *Stop
}
