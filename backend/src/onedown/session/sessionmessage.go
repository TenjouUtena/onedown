package session

import (
	"github.com/TenjouUtena/onedown/backend/src/onedown/solver"
	"github.com/google/uuid"
)

type SessionMessage interface{}

type JoinSession struct {
	SessionMessage
	Solver *solver.Solver
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
