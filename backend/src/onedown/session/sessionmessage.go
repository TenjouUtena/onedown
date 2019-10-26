package session

import (
	"encoding/json"
	"github.com/google/uuid"
)

type SessionMessage interface{}

type ClientSessionMessage interface {
	SessionMessage
	unmarshalClientPayload(payload []byte) (SessionMessage, error)
}

type JoinSession struct {
	Solver *Solver
}

func (msg JoinSession) unmarshalClientPayload(payload []byte) (SessionMessage, error) {
	unmarshaledMessage := JoinSession{}
	err := json.Unmarshal(payload, &unmarshaledMessage)
	return unmarshaledMessage, err
}

type LeaveSession struct { }

func (msg LeaveSession) unmarshalClientPayload(payload []byte) (SessionMessage, error) {
	unmarshaledMessage := LeaveSession{}
	err := json.Unmarshal(payload, &unmarshaledMessage)
	return unmarshaledMessage, err
}

type WriteSquare struct {
	Row    int
	Col    int
	Answer string
}

func (msg WriteSquare) unmarshalClientPayload(payload []byte) (SessionMessage, error) {
	unmarshaledMessage := WriteSquare{}
	err := json.Unmarshal(payload, &unmarshaledMessage)
	return unmarshaledMessage, err
}

type CheckSquares struct {
	RowIndices [2]int
	ColIndices [2]int
}

func (msg CheckSquares) unmarshalClientPayload(payload []byte) (SessionMessage, error) {
	unmarshaledMessage := CheckSquares{}
	err := json.Unmarshal(payload, &unmarshaledMessage)
	return unmarshaledMessage, err
}

type RevealSquares struct {
	RowIndices [2]int
	ColIndices [2]int
}

func (msg RevealSquares) unmarshalClientPayload(payload []byte) (SessionMessage, error) {
	unmarshaledMessage := RevealSquares{}
	err := json.Unmarshal(payload, &unmarshaledMessage)
	return unmarshaledMessage, err
}

var sessionClientMessages = map[string]ClientSessionMessage{
	"JoinSession":   JoinSession{},
	"LeaveSession":  LeaveSession{},
	"WriteSquare":   WriteSquare{},
	"CheckSquares":  CheckSquares{},
	"RevealSquares": RevealSquares{},
}

type GetSessions struct {
	SessionMessage
	ResponseChannel chan []uuid.UUID
}
