package mrmiddle

import "time"

// Timing
const (
	// POLLTIME is board polling interval time
	TIMING_POLL = 200 * time.Millisecond

	// FLIPTIME is the time of output for flip
	TIMING_FLIP = 1000 * time.Millisecond
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
	[8]string{"0A1A0", "0A3A2", "0A5A4", "0A7A6", "1A1A0", "1A3A2", "1A5A4", "1A7A6"},
	[8]string{"0B6B7", "0B4B5", "0B2B3", "0B0B1", "1B0B1", "1B2B3", "1B4B5", "1B6B7"},
	[8]string{"2A1A0", "2A3A2", "2A5A4", "2A7A6", "3A1A0", "3A3A2", "3A5A4", "3A7A6"},
	[8]string{"2B6B7", "2B4B5", "2B2B3", "2B0B1", "3B0B1", "3B2B3", "3B4B5", "3B6B7"},
	[8]string{"4A1A0", "4A3A2", "4A5A4", "4A7A6", "5A1A0", "5A3A2", "5A5A4", "5A7A6"},
	[8]string{"4B6B7", "4B4B5", "4B2B3", "4B0B1", "5B0B1", "5B2B3", "5B4B5", "5B6B7"},
	[8]string{"6A1A0", "6A3A2", "6A5A4", "6A7A6", "7A1A0", "7A3A2", "7A5A4", "7A7A6"},
	[8]string{"6B6B7", "6B4B5", "6B2B3", "6B0B1", "7B0B1", "7B2B3", "7B4B5", "7B6B7"},
}

// For debug
// var ExpanderMap = [][]string{
// 	[]string{"0A0A1", "0A2A3", "0A4A5", "0A6A7"},
// 	[]string{"0B0B1", "0B2B3", "0B4B5", "0B6B7"},
// }

// Motor pin assignment
const (
	MOTOR_PIN1 = "8"
	MOTOR_PIN2 = "10"
)

// Motor driver function state
var MOTOR_STOP = MotorState{0, 0}
var MOTOR_BRAKE = MotorState{1, 1}
var MOTOR_DRIVE = map[Pole]MotorState{
	N: MotorState{1, 0},
	S: MotorState{0, 1},
}
