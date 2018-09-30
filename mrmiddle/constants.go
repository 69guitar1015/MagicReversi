package mrmiddle

import "time"

// Timing
const (
	// POLLTIME is board polling interval time
	TIMING_POLL = 200 * time.Millisecond

	// FLIPTIME is the time of output for flip
	TIMING_FLIP = 500 * time.Millisecond
)

// Pole direction
const (
	// N pole
	N Pole = 1
	// S pole
	S Pole = -1
)

// Mapping list to map board cell to expander pin
// format: [expander id][input port][input pin][output port][output pin]
var ExpanderMap = [8][8]string{
	[8]string{"0A0A1", "0A2A3", "0A4A5", "0A6A7", "1A0A1", "1A2A3", "1A4A5", "1A6A7"},
	[8]string{"0B0B1", "0B2B3", "0B4B5", "0B6B7", "1B0B1", "1B2B3", "1B4B5", "1B6B7"},
	[8]string{"2A0A1", "2A2A3", "2A4A5", "2A6A7", "3A0A1", "3A2A3", "3A4A5", "3A6A7"},
	[8]string{"2B0B1", "2B2B3", "2B4B5", "2B6B7", "3B0B1", "3B2B3", "3B4B5", "3B6B7"},
	[8]string{"4A0A1", "4A2A3", "4A4A5", "4A6A7", "5A0A1", "5A2A3", "5A4A5", "5A6A7"},
	[8]string{"4B0B1", "4B2B3", "4B4B5", "4B6B7", "5B0B1", "5B2B3", "5B4B5", "5B6B7"},
	[8]string{"6A0A1", "6A2A3", "6A4A5", "6A6A7", "7A0A1", "7A2A3", "7A4A5", "7A6A7"},
	[8]string{"6B0B1", "6B2B3", "6B4B5", "6B6B7", "7B0B1", "7B2B3", "7B4B5", "7B6B7"},
}

// For debug
// var ExpanderMap = [][]string{
// 	[]string{"0A0A1", "0A2A3", "0A4A5", "0A6A7"},
// 	[]string{"0B0B1", "0B2B3", "0B4B5", "0B6B7"},
// }

// Motor pin assignment
const (
	MOTOR_PIN1 = "11"
	MOTOR_PIN2 = "12"
)

// Motor driver function state
var MOTOR_STOP = MotorState{0, 0}
var MOTOR_BRAKE = MotorState{1, 1}
var MOTOR_DRIVE = map[Pole]MotorState{
	N: MotorState{1, 0},
	S: MotorState{0, 1},
}
