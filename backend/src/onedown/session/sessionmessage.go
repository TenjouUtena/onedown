package session

import (
	"github.com/TenjouUtena/onedown/backend/src/onedown/solver"
	"github.com/google/uuid"
)

type SessionMessage interface{}

type JoinSession struct {
	SessionMessage
	solver *solver.Solver
}

type LeaveSession struct {
	SessionMessage
	solver uuid.UUID
}

type WriteSquare struct {
	SessionMessage
	solver uuid.UUID
	row    int
	col    int
	answer string
}

type CheckSquares struct {
	SessionMessage
	solver     uuid.UUID
	rowIndices [2]int
	colIndices [2]int
}
