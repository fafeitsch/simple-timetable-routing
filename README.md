Simple Timetable Routing
===

This package provides a simple algorithm to compute
fastest routes in public transportation networks. This project
was mainly meant as a study project to understand
the algorithm and to experiment with different data models
for public transportation networks in Go.

This implementation is **not** suited for production
usage, it may still have bugs and is not at all optimized. As of now,
the implementation assumes a minimum change time at stations of five minutes.
It is planed to make this constant customizable.

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
   
Implementation Details
---

The package implements the time-dependent variant of [Dijkstra's Algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm).
The time-dependent model is a classical solution for finding routes in public transportation networks.
For a thorough explanation see for example [this paper ("Time-Dependent Route Planning" by Daniel Delling and Dorothea Wagner)](https://i11www.iti.kit.edu/extra/publications/dw-tdrp-09.pdf).

The advantage of the time-dependent approach is that the original algorithm by Dijstrka must hardly be changed.
The first step is to build a graph that resembles the real network graph: Every station is represented by a vertex and
if there is at least one line connecting two stations, then there is an edge between the stations in the graph.

In normal Dijstrka the weight of the edge *w(e)* is constant. In the time-dependent model, the weight
changes over time, thus we have a weight function of *w(e,t)* where t is the time. Given two stations
*u* and *v* which are connected by lines *A* and *B*. If a train of line *A* departures at 14:30 and a train of line *B*
departures at 14:40 and the travel time to *v* are three minutes for both lines, then
*w((u,v), 14:25)* is eight minutes and *w((u,v), 14:39)* is four minutes.

My implementation is (at the moment) not focussed on performance, but rather was meant to be a
working method for route finding in public networks. Thus, my implementation does **neither** make use
of fibonnaci heaps in the plain Dijsktra algorithm, **nor** is the implementation of *w(e,t)* very
sophisticated. 

License
---
See LICENSE file.
