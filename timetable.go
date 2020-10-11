package routing

import (
	"strings"
	"time"
)

type Timetable struct {
	Stops []Stop
}

type Stop struct {
	Id     string
	Name   string
	Events []Event
}

type StopGroup struct {
	Stops []Stop
}

func (s *StopGroup) String() string {
	result := make([]string, 0, len(s.Stops))
	for _, stop := range s.Stops {
		result = append(result, stop.Name)
	}
	return strings.Join(result, ",")
}

type Line struct {
	Id   string
	Name string
}

type Event struct {
	Arrival   *time.Time
	Departure *time.Time
	Line      *Line
}

type Connection struct {
	duration time.Duration
}
