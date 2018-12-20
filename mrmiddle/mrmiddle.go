package mrmiddle

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func wrapError(err error) error {
	prefix := "Middleware Error"
	if strings.HasPrefix(err.Error(), prefix) {
		return err
	}

	return fmt.Errorf("%s: %s", prefix, err)
}

// Pole is a representation of a polar direction N and S
// constant variable N and S are defined in constants.go
type Pole int

// MotorState is motor driver state
type MotorState [2]int

// position represents a expander pin belonging to MrMiddle
type position struct {
	i    uint8
	port string
	pin  uint8
}

// Board represents board status
type Board [8][8]int8

// Print prints board status
func (b *Board) Print() {
	fmt.Println("-----------------------------------------")
	for i := range *b {
		fmt.Printf("|")
		for j := range (*b)[i] {
			fmt.Printf("\t%v", (*b)[i][j])
		}
		fmt.Printf("\t|\n")
	}
	fmt.Println("-----------------------------------------")
}

// MrMiddle is Magic Reversi's middle ware object
type MrMiddle struct {
	master    *raspi.Adaptor
	expanders [8]*i2c.MCP23017Driver
	motorPin  [2]*gpio.DirectPinDriver
	readMap   [8][8]position
	writeMap  [8][8]position

	// For debug
	// expanders [1]*i2c.MCP23017Driver
	// motorPin  [2]*gpio.DirectPinDriver
	// readMap   [2][4]position
	// writeMap  [2][4]position
}

// NewMrMiddle returns MrMiddle instance
func NewMrMiddle() (mm *MrMiddle, err error) {
	mm = &MrMiddle{}
	mm.master = raspi.NewAdaptor()

	// IO expanders
	for i := range mm.expanders {
		mm.expanders[i] = i2c.NewMCP23017Driver(
			mm.master,
			i2c.WithBus(1),
			i2c.WithAddress(0x20+i),
			i2c.WithMCP23017Bank(0),
			i2c.WithMCP23017Mirror(0),
			i2c.WithMCP23017Seqop(0),
			i2c.WithMCP23017Disslw(0),
			i2c.WithMCP23017Haen(0),
			i2c.WithMCP23017Odr(0),
			i2c.WithMCP23017Intpol(0),
		)
	}

	for i := range mm.readMap {
		for j := range mm.readMap[i] {
			enc := ExpanderMap[i][j]
			id, _ := strconv.Atoi(string(enc[0]))
			readPort := string(enc[1])
			readPin, _ := strconv.Atoi(string(enc[2]))
			writePort := string(enc[3])
			writePin, _ := strconv.Atoi(string(enc[4]))

			mm.readMap[i][j] = position{uint8(id), readPort, uint8(readPin)}
			mm.writeMap[i][j] = position{uint8(id), writePort, uint8(writePin)}
		}
	}

	// Motor Driver
	mm.motorPin[0] = gpio.NewDirectPinDriver(mm.master, MOTOR_PIN1)
	mm.motorPin[1] = gpio.NewDirectPinDriver(mm.master, MOTOR_PIN2)

	return
}

// Init is initialization function of MrMiddle
func (mm *MrMiddle) Init() error {
	log.Println("Initialize circuit...")

	// Raspberry pi zero W
	if err := mm.master.Connect(); err != nil {
		return wrapError(err)
	}

	// Motor driver
	if err := mm.setMotorState(MOTOR_STOP); err != nil {
		return wrapError(err)
	}

	// I/O Expander
	for _, exp := range mm.expanders {
		exp.Start()
	}

	if err := mm.writeAllLow(); err != nil {
		return wrapError(err)
	}

	return nil
}

// Finalize execute finalizing process
func (mm *MrMiddle) Finalize() (err error) {
	fmt.Println("Finalize...")

	// Stop motor
	if err := mm.setMotorState(MOTOR_STOP); err != nil {
		return err
	}

	// Shut all gates
	if err := mm.writeAllLow(); err != nil {
		return err
	}

	// Finalize master
	if err := mm.master.Finalize(); err != nil {
		return err
	}

	return
}

/*
	Output function
*/

// set motor driver state
func (mm *MrMiddle) setMotorState(state MotorState) error {
	for i := 0; i < 2; i++ {
		err := mm.motorPin[i].DigitalWrite(byte(state[i]))
		if err != nil {
			return wrapError(err)
		}
	}
	return nil
}

// write `val` at `pos`
func (mm *MrMiddle) writeAt(pos position, val uint8) error {
	err := mm.expanders[pos.i].WriteGPIO(pos.pin, val, pos.port)
	if err != nil {
		// return wrapError(err)
		fmt.Println(pos, err)
	} else {
		fmt.Println(pos)
	}
	return nil
}

func (mm *MrMiddle) writeAllLow() error {
	for i := range mm.writeMap {
		for j := range mm.writeMap[i] {
			pos := mm.writeMap[i][j]
			err := mm.writeAt(pos, 0)
			if err != nil {
				return wrapError(err)
			}
		}
	}
	return nil
}

// Flip outputs at (i, j) cell in `TIMING_FLIP` time
func (mm *MrMiddle) Flip(i uint8, j uint8, pole Pole) error {
	// Open the gate at (i, j)
	if err := mm.writeAt(mm.writeMap[i][j], 1); err != nil {
		return wrapError(err)
	}

	// Turn on the motor driver, direction is `Pole`
	if err := mm.setMotorState(MOTOR_DRIVE[pole]); err != nil {
		return wrapError(err)
	}

	// Wait
	time.Sleep(TIMING_FLIP)

	// Turn off
	if err := mm.setMotorState(MOTOR_STOP); err != nil {
		return wrapError(err)
	}

	// Close the gate at (i, j)
	if err := mm.writeAt(mm.writeMap[i][j], 0); err != nil {
		return wrapError(err)
	}

	return nil
}

/*
  Input function
*/

func (mm *MrMiddle) readAt(pos position) (uint8, error) {
	val, err := mm.expanders[pos.i].ReadGPIO(pos.pin, pos.port)
	if err != nil {
		return 0, wrapError(err)
	}
	return val, nil
}

// ReadBoard reads current board status
func (mm *MrMiddle) ReadBoard() (board Board, err error) {
	for i := range mm.readMap {
		for j := range mm.readMap[i] {
			val, err := mm.readAt(mm.readMap[i][j])

			if err != nil {
				board[i][j] = -1
			} else {
				board[i][j] = int8(val)
			}
		}
	}
	return
}
