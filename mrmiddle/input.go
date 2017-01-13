package mrmiddle

import "time"

// read given y line
func (mm *MrMiddle) readLine(y int) (r row, err error) {
	addr, gpio := y2AddrAndGpio(y)

	mm.e.I2cStart(addr)

	err = mm.e.I2cWrite(addr, []byte{byte(gpio)})

	if checkError(err) {
		return row{}, wrapError(err)
	}

	data, err := mm.e.I2cRead(addr, 1)

	if checkError(err) {
		return row{}, wrapError(err)
	}

	r = byte2Row(data[0])

	if y%2 == 0 {
		r = r.reversed()
	}

	return
}

// read both A and B bits of given address of the Expander
func (mm MrMiddle) readAB(addr int) (byteSet [2]row, err error) {
	mm.e.I2cStart(addr)

	err = mm.e.I2cWrite(addr, []byte{GPIOA})

	if checkError(err) {
		return [2]row{}, wrapError(err)
	}

	data, err := mm.e.I2cRead(addr, 2)

	if checkError(err) {
		return [2]row{}, wrapError(err)
	}

	byteSet[0] = byte2Row(data[0]).reversed()
	byteSet[1] = byte2Row(data[1])

	return
}

// ReadWholeBoard returns whole board status
func (mm *MrMiddle) readWholeBoard() (byteSet [8]row, err error) {
	for i, addr := range EXIA {
		data, err := mm.readAB(addr)

		if checkError(err) {
			return [8]row{}, wrapError(err)
		}

		byteSet[2*i] = data[0].reversed()
		byteSet[2*i+1] = data[1]
	}

	return
}

// GetInput waits until board state changes and return x, y
func (mm *MrMiddle) GetInput() (int, int, error) {
	old, err := mm.readWholeBoard()

	if checkError(err) {
		return -1, -1, wrapError(err)
	}

	for {
		crr, err := mm.readWholeBoard()

		if checkError(err) {
			return -1, -1, wrapError(err)
		}

		if crr != old {
			for i, r := range crr {
				for j, v := range r {
					// if current status is True and old status is False then reutrn (x, y)
					if v && !old[i][j] {
						return j + 1, i + 1, nil
					}
				}
			}
		}

		time.Sleep(POLLTIME)
	}
}
