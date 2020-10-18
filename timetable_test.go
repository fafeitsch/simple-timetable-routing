package routing

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStop_groupEvents(t *testing.T) {
	zoo := &Stop{Name: "Zoo", Id: "ZO"}
	mall := &Stop{Name: "Mall", Id: "MA"}
	court := &Stop{Name: "Court", Id: "CO"}
	mainStreet := &Stop{Name: "Main Street", Id: "MS"}

	e1 := Event{NextStop: zoo}
	e2 := Event{NextStop: mall}
	e3 := Event{NextStop: zoo}
	e4 := Event{NextStop: court}
	e5 := Event{NextStop: zoo}
	e6 := Event{NextStop: mainStreet}
	e7 := Event{NextStop: mainStreet}

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

	e1 := Event{Line: southBound, Departure: date("14:30"), ArrivalAtNextStop: date("14:35")}
	e2 := Event{Line: southBound, Departure: date("14:39"), ArrivalAtNextStop: date("14:44")}
	e3 := Event{Line: southBound, Departure: date("14:48"), ArrivalAtNextStop: date("14:53")}
	e4 := Event{Line: harbour, Departure: date("14:35"), ArrivalAtNextStop: date("14:48")}
	e6 := Event{Line: harbourExpress, Departure: date("14:35"), ArrivalAtNextStop: date("14:47")}

	group := EventGroup([]Event{e1, e2, e3, e4, e6})
	function := group.weightFunction()
	t.Run("without change", func(t *testing.T) {
		now := date("14:34")
		duration, line, b := function(now, harbour)
		assert.Equal(t, 14*time.Minute, duration, "duration is wrong")
		assert.Equal(t, harbour, line, "line after event is wrong")
		assert.True(t, b, "connection should be found")
	})
	t.Run("with change", func(t *testing.T) {
		now := date("14:30")
		duration, line, b := function(now, harbour)
		assert.Equal(t, 14*time.Minute, duration, "duration is wrong")
		assert.Equal(t, southBound, line, "line after event is wrong")
		assert.True(t, b, "connection should be found")
	})
	t.Run("no departure found", func(t *testing.T) {
		now := date("16:00")
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

func date(date string) time.Time {
	compiled, err := time.Parse(time.RFC3339, "2020-10-15T"+date+":00Z")
	if err != nil {
		panic(err)
	}
	return compiled
}
