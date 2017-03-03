package mrmiddle

import (
	"fmt"
	"time"
)

// write byte data to designated line
func (mm *MrMiddle) writeByte(y int, v byte) (err error) {
	addr, gpio := y2AddrAndGpio(y - 1)

	if err = mm.e.I2cStart(addr); checkError(err) {
		return
	}

	if gpio == GPIOA {
		v = byte2Row(v).reversed().toByte()
	}

	data := []byte{byte(gpio), v}

	return mm.e.I2cWrite(addr, data)
}

func (mm *MrMiddle) writeAllLow() (err error) {
	for y := 0; y < 8; y++ {
		if err = mm.writeByte(y, 0x00); checkError(err) {
			return
		}
	}

	return
}

// driveCoil drives coils as given pole direction
func (mm *MrMiddle) driveCoil(pd Pole) (err error) {
	switch pd {
	case N:
		if err = writeDuty(IN1, PWMLEVEL); checkError(err) {
			return wrapError(err)
		}

		if err = pwmEnable(IN1, "1"); checkError(err) {
			return wrapError(err)
		}

	case S:
		if err = writeDuty(IN2, PWMLEVEL); checkError(err) {
			return wrapError(err)
		}

		if err = pwmEnable(IN2, "1"); checkError(err) {
			return wrapError(err)
		}
	}

	return
}

// driveCoil drives coils as given pole direction
func (mm *MrMiddle) ochitukeCoil(pd Pole) (err error) {
	fmt.Println("落ち着け")
	switch pd {
	case N:

		if err = writeDuty(IN1, 0); checkError(err) {
			return wrapError(err)
		}

		if err = pwmEnable(IN1, "0"); checkError(err) {
			return wrapError(err)
		}

	case S:
		if err = writeDuty(IN2, 0); checkError(err) {
			return wrapError(err)
		}

		if err = pwmEnable(IN2, "0"); checkError(err) {
			return wrapError(err)
		}
	}

	return
}

// releaseCoil releases coils
func (mm *MrMiddle) releaseCoil() (err error) {

	fmt.Println("3ばん")
	if err = writeDuty(IN1, 0); checkError(err) {
		return wrapError(err)
	}

	fmt.Println("4ばん")
	if err = writeDuty(IN2, 0); checkError(err) {
		return wrapError(err)
	}
	fmt.Println("1ばん")
	if err = pwmEnable(IN1, "0"); checkError(err) {
		return wrapError(err)
	}

	fmt.Println("2ばん")
	if err = pwmEnable(IN2, "0"); checkError(err) {
		return wrapError(err)
	}

	fmt.Println("5ばん")

	return
}

// HighWhile make (x, y) to High while ms[msec]
func (mm *MrMiddle) HighWhile(x int, y int, ms time.Duration, pd Pole) (err error) {
	bits := byte(0x01 << uint(x-1))

	if err = mm.writeByte(y, bits); checkError(err) {
		return
	}

	if err = mm.driveCoil(pd); checkError(err) {
		return
	}

	time.Sleep(ms)

	if err = mm.ochitukeCoil(pd); checkError(err) {
		return
	}

	return mm.writeByte(y, 0x00)
}

// Flip flips a stone at (x, y)
func (mm *MrMiddle) Flip(x int, y int, pd Pole) (err error) {
	err = mm.HighWhile(x, y, FLIPTIME, pd)

	if checkError(err) {
		return wrapError(err)
	}

	return
}
