// Package routing offers an algorithm to find the fastest route in a public transportation network.
//
// The cornerstone of the package is the Timetable type which contains all information
// of the public transportation network.
//
// Users of this package have to define Stops – physical locations where vehicles stop. A stop
// must have an Id and a list of departure events. All events consist of a departure time,
// a reference to the line they belong to, and information about the travel time to the next stop
// and which stop this is. The list of stops is used to create a timetable. This timetable can then
// be queried for shortest routes.
//
// The departure times are given in the format "15:04" (resp. "HH:MM"). Queries always contain a
// real date and time (time.Time type). The departure times are then interpreted to take place at the certain date.
// In order to simulate timetables spanning more than one day, departure times can also be given
// four hours bigger than 23 o'clock, e.g. 26:34 means 02:34 on the second day.
package routing
