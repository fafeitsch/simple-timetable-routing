package routing

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTimetable_Query(t *testing.T) {
	mainStation := &Stop{Id: "MS", Name: "Main Station", Events: make([]Event, 0, 0)}
	docksAE := &Stop{Id: "DAE", Name: "Docks A–E", Events: make([]Event, 0, 0)}
	docksFG := &Stop{Id: "DFG", Name: "Docks F and G", Events: make([]Event, 0, 0)}
	historicMall := &Stop{Id: "HM", Name: "Historic Mall", Events: make([]Event, 0, 0)}
	schusterStreet := &Stop{Id: "SS", Name: "Schuster Street", Events: make([]Event, 0, 0)}
	marketPlace := &Stop{Id: "MP", Name: "Market Place", Events: make([]Event, 0, 0)}
	airport := &Stop{Id: "AR", Name: "Airport", Events: make([]Event, 0, 0)}
	northAvenue := &Stop{Id: "NA", Name: "North Avenue", Events: make([]Event, 0, 0)}
	chalet := &Stop{Id: "CH", Name: "Chalet", Events: make([]Event, 0, 0)}
	northEnd := &Stop{Id: "NE", Name: "North End", Events: make([]Event, 0, 0)}

	blueLine := &Line{Id: "#0000FF", Name: "Blue Line", startStop: mainStation, Segments: []Segment{
		{TravelTime: 2 * time.Minute, NextStop: northAvenue},
		{TravelTime: 3 * time.Minute, NextStop: historicMall},
		{TravelTime: 1 * time.Minute, NextStop: schusterStreet},
		{TravelTime: 2 * time.Minute, NextStop: chalet}},
	}
	redLine := &Line{Id: "#FF0000", Name: "Red Line", startStop: northEnd, Segments: []Segment{
		{TravelTime: 2 * time.Minute, NextStop: northAvenue},
		{TravelTime: 2 * time.Minute, NextStop: mainStation},
		{TravelTime: 3 * time.Minute, NextStop: docksAE},
		{TravelTime: 5 * time.Minute, NextStop: airport}},
	}
	greenLine := &Line{Id: "#00FF00", Name: "Green Line", startStop: docksAE, Segments: []Segment{
		{TravelTime: 2 * time.Minute, NextStop: docksFG},
	}}
	whiteLine := &Line{Id: "#FFFFFF", Name: "White Line", startStop: marketPlace, Segments: []Segment{
		{TravelTime: 2 * time.Minute, NextStop: chalet},
		{TravelTime: 2 * time.Minute, NextStop: schusterStreet},
		{TravelTime: 1 * time.Minute, NextStop: historicMall},
		{TravelTime: 3 * time.Minute, NextStop: northAvenue},
		{TravelTime: 2 * time.Minute, NextStop: mainStation},
	}}
	_, _ = greenLine, whiteLine

	for hour := 8; hour < 20; hour++ {
		// blue line (every twenty minutes)
		for j := 0; j < 3; j++ {
			mainStation.Events = append(mainStation.Events, Event{Departure: Time{Hour: hour, Minute: j*20 + 5}, Segment: &blueLine.Segments[0], Line: blueLine})
			northAvenue.Events = append(northAvenue.Events, Event{Departure: Time{Hour: hour, Minute: j*20 + 7}, Segment: &blueLine.Segments[1], Line: blueLine})
			historicMall.Events = append(historicMall.Events, Event{Departure: Time{Hour: hour, Minute: j*20 + 10}, Segment: &blueLine.Segments[2], Line: blueLine})
			schusterStreet.Events = append(schusterStreet.Events, Event{Departure: Time{Hour: hour, Minute: j*20 + 11}, Segment: &blueLine.Segments[3], Line: blueLine})
		}
	}

	for hour := 10; hour < 20; hour++ {
		// red line (from 10 to 20 o'clock) every five minutes
		for j := 0; j < 12; j++ {
			northEnd.Events = append(northEnd.Events, Event{Departure: Time{Hour: hour, Minute: j * 5}, Segment: &redLine.Segments[0], Line: redLine})
			northAvenue.Events = append(northAvenue.Events, Event{Departure: Time{Hour: hour, Minute: j*5 + 2}, Segment: &redLine.Segments[1], Line: redLine})
			mainStation.Events = append(mainStation.Events, Event{Departure: Time{Hour: hour, Minute: j*5 + 4}, Segment: &redLine.Segments[2], Line: redLine})
			if j == 11 {
				docksAE.Events = append(docksAE.Events, Event{Departure: Time{Hour: hour + 1, Minute: 2}, Segment: &redLine.Segments[3], Line: redLine})
			} else {
				docksAE.Events = append(docksAE.Events, Event{Departure: Time{Hour: hour + 1, Minute: j*5 + 7}, Segment: &redLine.Segments[3], Line: redLine})
			}
		}
	}

	timetable := NewTimetable([]*Stop{mainStation, docksAE, docksFG, historicMall, schusterStreet, marketPlace, airport, northAvenue, chalet, northEnd})

	t.Run("single line", func(j *testing.T) {
		connection := timetable.Query(northAvenue, schusterStreet, date("14:34"))
		assert.Equal(t, 1, len(connection.Legs), "number of legs in the connection")
		assert.Equal(t, northAvenue, connection.Legs[0].FirstStop, "first stop wrong")
		assert.Equal(t, schusterStreet, connection.Legs[0].LastStop, "last stop wrong")
		assert.Equal(t, blueLine, connection.Legs[0].Line, "line is wrong")
		assert.Equal(t, date("14:51"), connection.Arrival, "time is wrong")
	})
	t.Run("single line(to late)", func(t *testing.T) {
		connection := timetable.Query(schusterStreet, chalet, date("21:23"))
		assert.Nil(t, connection, "there is no connection any more")
	})
	t.Run("single line(no connection)", func(t *testing.T) {
		connection := timetable.Query(mainStation, northEnd, date("10:00"))
		assert.Nil(t, connection, "there is no connection to the target")
	})
	t.Run("single line", func(j *testing.T) {
		connection := timetable.Query(mainStation, northAvenue, date("8:04"))
		assert.Equal(t, 1, len(connection.Legs), "number of legs in the connection")
		assert.Equal(t, mainStation, connection.Legs[0].FirstStop, "first stop wrong")
		assert.Equal(t, northAvenue, connection.Legs[0].LastStop, "last stop wrong")
		assert.Equal(t, blueLine, connection.Legs[0].Line, "line is wrong")
		assert.Equal(t, date("8:07"), connection.Arrival, "time is wrong")
	})
	t.Run("blue line/red line", func(j *testing.T) {
		connection := timetable.Query(northEnd, chalet, date("9:30"))
		assert.Equal(t, 2, len(connection.Legs), "number of legs in the connection")
		assert.Equal(t, northEnd, connection.Legs[0].FirstStop, "first stop wrong")
		assert.Equal(t, northAvenue, connection.Legs[0].LastStop, "last stop wrong")
		assert.Equal(t, redLine, connection.Legs[0].Line, "line is wrong")
		assert.Equal(t, northAvenue, connection.Legs[1].FirstStop, "first stop wrong")
		assert.Equal(t, chalet, connection.Legs[1].LastStop, "last stop wrong")
		assert.Equal(t, blueLine, connection.Legs[1].Line, "line is wrong")
		assert.Equal(t, date("10:13"), connection.Arrival, "time is wrong")
	})
	t.Run("blue line/red line(switch time)", func(j *testing.T) {
		connection := timetable.Query(northEnd, chalet, date("10:25"))
		assert.Equal(t, 2, len(connection.Legs), "number of legs in the connection")
		assert.Equal(t, northEnd, connection.Legs[0].FirstStop, "first stop wrong")
		assert.Equal(t, northAvenue, connection.Legs[0].LastStop, "last stop wrong")
		assert.Equal(t, redLine, connection.Legs[0].Line, "line is wrong")
		assert.Equal(t, northAvenue, connection.Legs[1].FirstStop, "first stop wrong")
		assert.Equal(t, chalet, connection.Legs[1].LastStop, "last stop wrong")
		assert.Equal(t, blueLine, connection.Legs[1].Line, "line is wrong")
		assert.Equal(t, date("10:53"), connection.Arrival, "time is wrong")
	})
}

func TestStop_groupEvents(t *testing.T) {
	zoo := &Stop{Name: "Zoo", Id: "ZO"}
	mall := &Stop{Name: "Mall", Id: "MA"}
	court := &Stop{Name: "Court", Id: "CO"}
	mainStreet := &Stop{Name: "Main Street", Id: "MS"}

	e1 := Event{Segment: &Segment{NextStop: zoo}}
	e2 := Event{Segment: &Segment{NextStop: mall}}
	e3 := Event{Segment: &Segment{NextStop: zoo}}
	e4 := Event{Segment: &Segment{NextStop: court}}
	e5 := Event{Segment: &Segment{NextStop: zoo}}
	e6 := Event{Segment: &Segment{NextStop: mainStreet}}
	e7 := Event{Segment: &Segment{NextStop: mainStreet}}

	centralStation := &Stop{Name: "Central Station", Id: "CS", Events: []Event{e1, e2, e3, e4, e5, e6, e7}}
	groups := centralStation.groupEvents()

	assert.Equal(t, 4, len(groups), "number of groups")
	assert.Equal(t, EventGroup{e1, e3, e5}, groups[zoo.Id], "zoo group members")
	assert.Equal(t, EventGroup{e2}, groups[mall.Id], "mall group members")
	assert.Equal(t, EventGroup{e4}, groups[court.Id], "court group members")
	assert.Equal(t, EventGroup{e6, e7}, groups[mainStreet.Id], "mainStreet group members")
}

func TestEventGroup_WeightFunction(t *testing.T) {
	southBound := &Line{Name: "1 SouthBound", Id: "1"}
	harbour := &Line{Name: "2 Harbour", Id: "2"}
	harbourExpress := &Line{Name: "2a Harbour", Id: "2a"}

	e1 := Event{Line: southBound, Departure: ParseTime("14:30"), Segment: &Segment{TravelTime: 5 * time.Minute}}
	e2 := Event{Line: southBound, Departure: ParseTime("14:39"), Segment: &Segment{TravelTime: 5 * time.Minute}}
	e3 := Event{Line: southBound, Departure: ParseTime("14:48"), Segment: &Segment{TravelTime: 5 * time.Minute}}
	e4 := Event{Line: harbour, Departure: ParseTime("14:35"), Segment: &Segment{TravelTime: 13 * time.Minute}}
	e6 := Event{Line: harbourExpress, Departure: ParseTime("14:35"), Segment: &Segment{TravelTime: 12 * time.Minute}}

	group := EventGroup([]Event{e1, e2, e3, e4, e6})
	t.Run("without change", func(t *testing.T) {
		now := date("14:34")
		function := group.weightFunction(now)
		duration, line, b := function(now, harbour)
		assert.Equal(t, 14*time.Minute, duration, "duration is wrong")
		assert.Equal(t, harbour, line, "line after event is wrong")
		assert.True(t, b, "connection should be found")
	})
	t.Run("with change", func(t *testing.T) {
		now := date("14:30")
		function := group.weightFunction(now)
		duration, line, b := function(now, harbour)
		assert.Equal(t, 14*time.Minute, duration, "duration is wrong")
		assert.Equal(t, southBound, line, "line after event is wrong")
		assert.True(t, b, "connection should be found")
	})
	t.Run("no departure found", func(t *testing.T) {
		now := date("16:00")
		function := group.weightFunction(now)
		_, _, b := function(now, harbourExpress)
		assert.False(t, b, "no connection should be found any more")
	})
}

func Test_createConnection(t *testing.T) {
	southBound := &Line{Name: "1 SouthBound", Id: "1"}
	harbour := &Line{Name: "2 Harbour", Id: "2"}

	zoo := &Stop{Name: "Zoo", Id: "ZO"}
	mall := &Stop{Name: "Mall", Id: "MA"}
	court := &Stop{Name: "Court", Id: "CO"}
	mainStreet := &Stop{Name: "Main Street", Id: "MS"}
	centralStation := &Stop{Name: "Central Station", Id: "CS"}

	v1 := &vertex{data: zoo}
	v2 := &vertex{data: mall, currentLine: southBound}
	v3 := &vertex{data: court, currentLine: southBound}
	v4 := &vertex{data: mainStreet, currentLine: southBound}
	v5 := &vertex{data: centralStation, currentLine: harbour}
	path := []*vertex{v1, v2, v3, v4, v5}

	t.Run("test big", func(t *testing.T) {
		got := createConnection(path)
		require.Equal(t, 2, len(got.Legs), "number of legs is wrong")
		assert.Equal(t, southBound, got.Legs[0].Line, "line of leg 0 not correct")
		assert.Equal(t, harbour, got.Legs[1].Line, "line of leg 1 not correct")
		assert.Equal(t, zoo, got.Legs[0].FirstStop, "first stop of leg 0 not correct")
		assert.Equal(t, mainStreet, got.Legs[0].LastStop, "last stop of leg 0 not correct")
		assert.Equal(t, mainStreet, got.Legs[1].FirstStop, "first stop of leg 1 not correct")
		assert.Equal(t, centralStation, got.Legs[1].LastStop, "last stop of leg 1 not correct")
	})
	t.Run("small", func(t *testing.T) {
		got := createConnection(path[3:])
		require.Equal(t, 1, len(got.Legs), "number of legs is wrong")
		assert.Equal(t, harbour, got.Legs[0].Line, "line of leg 0 not correct")
		assert.Equal(t, mainStreet, got.Legs[0].FirstStop, "first stop of leg 1 not correct")
		assert.Equal(t, centralStation, got.Legs[0].LastStop, "last stop of leg 1 not correct")
	})
}

func date(moment string) time.Time {
	compiled, err := time.Parse(time.RFC3339, "2020-10-15T"+moment+":00Z")
	if err != nil {
		panic(err)
	}
	return compiled
}
