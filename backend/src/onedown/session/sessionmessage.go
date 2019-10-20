package session

import (
	"github.com/google/uuid"
)

type SessionMessage interface{}

type JoinSession struct {
	SessionMessage
	Solver *Solver
}

type LeaveSession struct {
	SessionMessage
	Solver uuid.UUID
}

type WriteSquare struct {
	SessionMessage
	Solver uuid.UUID
	Row    int
	Col    int
	Answer string
}

type CheckSquares struct {
	SessionMessage
	Solver     uuid.UUID
	RowIndices [2]int
	ColIndices [2]int
}

type GetSessions struct {
	SessionMessage
	ResponseChannel chan []uuid.UUID
}
