package mrmiddle

import "time"

//
//	IO EXPANDER CONFIGS
//
const (
	// IODIRA is 0x00
	IODIRA = 0x00 + iota
	// IODIRB is 0x01
	IODIRB
	// IPOLA is 0x02
	IPOLA
	// IPOLB is 0x03
	IPOLB
	// GPINTENA is 0x04
	GPINTENA
	// GPINTENB is 0x05
	GPINTENB
	// DEFVALA is 0x06
	DEFVALA
	// DEFVALB is 0x07
	DEFVALB
	// INTCONA is 0x08
	INTCONA
	// INTCONB is 0x09
	INTCONB
	// IOCON is 0x0A
	IOCON
	// IOCON2 is 0x0B
	IOCON2
	// GPPUA is 0x0C
	GPPUA
	// GPPUB is 0x0D
	GPPUB
	// INTFA is 0x0E
	INTFA
	// INTFB is 0x0F
	INTFB
	// INTCAPA is 0x10
	INTCAPA
	// INTCAPB is 0x11
	INTCAPB
	// GPIOA is 0x12
	GPIOA
	// GPIOB is 0x13
	GPIOB
	// OLATA is 0x14
	OLATA
	// OLATB is 0x15
	OLATB
)

// EXIA is I/O expander address for read
var EXIA = [4]int{0x20, 0x21, 0x22, 0x23}

// EXOA is I/O expander address for write
var EXOA = [4]int{0x24, 0x25, 0x26, 0x27}

//
//	Driver IC CONFIGS
//

const (
	// VCC is 5V pin
	VCC = "2"
	// IN1 is IN1 pin
	IN1 = "13"
	// IN2 is IN2 pin
	IN2 = "12"
)

//
//	GENERAL CONFIGS
//
const (
	// POLLTIME is board polling interval time
	POLLTIME = 200 * time.Millisecond

	// FLIPTIME is the time of output for flip
	FLIPTIME = 500 * time.Millisecond
)

// Pole represents magnetic poll direction
// N = 1 and S = -1
type Pole int

const (
	// N pole
	N Pole = 1
	// S pole
	S Pole = -1
)
