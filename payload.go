package main

import (
	"encoding/json"
	"fmt"

	"github.com/69guitar1015/MagicReversi/mrmiddle"
)

type GetBoardPayload struct {
	board mrmiddle.Board
}

func (p *GetBoardPayload) Compose() []byte {
	b := []byte{}
	for i := range (*p).board {
		for j := range (*p).board[i] {
			b = append(b, byte((*p).board[i][j]))
		}
	}
	return b
}

type NotifyPayload struct {
	x     uint8
	y     uint8
	state uint8
}

func (p *NotifyPayload) Compose() []byte {
	return []byte{(*p).x, (*p).y, (*p).state}
}

type ErrorPayload struct {
	err string
}

func (p *ErrorPayload) Compose() []byte {
	jsonBytes, err := json.Marshal(map[string]interface{}{
		"error": (*p).err,
	})

	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return nil
	}
	return jsonBytes
}
