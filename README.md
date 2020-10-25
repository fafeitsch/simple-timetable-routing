Simple Timetable Routing
===

This package provides a simple algorithm to compute
fastest routes in public transportation networks. This project
was mainly meant as a study project to understand
the algorithm and to experiment with different data models
for public transportation networks in Go.

This implementation is **not** suited for production
usage, it may still have bugs and is not at all optimized.

Quick Guide:
---
1. Define some stops with Id and name:
    ```go
    mainStation := NewStop("MS","Main Station")
    docksAE := NewStop("DAE","Docks A–E")
    docksFG := NewStop("DFG", "Docks F and G")
    historicMall := NewStop( "HM","Historic Mall")
    schusterStreet := NewStop("SS","Schuster Street")
    marketPlace := NewStop("MP", "Market Place")
    airport := NewStop("AR","Airport")
    northAvenue := NewStop("NA","North Avenue")
    chalet := NewStop("CH", "Chalet")
    northEnd := NewStop("NE" ,"North End")
    ```
2. Define your lines:
    ```go
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
    ```
3. Define events on the stops:
    ```go
       mainStation.Events = []Event{
         {Departure: "8:02", Line: blueLine, LineSegment: blueLine.Segments[2]},
         {Departure: "8:05", Line: redLine, LineSegment: redLine.Segments[0]}
         ...
       }   
    ```
4. Create a timetable and query it:
    ```go
       timetable := NewTimetable([]*Stop[mainstation, ...])
       connection := timetable.query(historicMall, chalet, time.Now())
    ```