package solver

import (
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/TenjouUtena/onedown/backend/src/onedown/session"
	"github.com/google/uuid"
)

type SolverMessage interface{}

type CheckResult struct {
	SolverMessage
	StartRow int      `json:"startRow"`
	StartCol int      `json:"startCol"`
	Result   [][]bool `json:"result"`
}

type SquareUpdated struct {
	SolverMessage
	Row      int       `json:"row"`
	Col      int       `json:"col"`
	NewValue string    `json:"newValue"`
	FilledBy uuid.UUID `json:"filledBy"`
}

// todo: add user information when informing of a solver, instead of just UUID
type PuzzleState struct {
	SolverMessage
	Solvers     []uuid.UUID          `json:"solvers"`
	Puzzle      *puzzle.Puzzle       `json:"puzzle"`
	PuzzleState *session.PuzzleState `json:"puzzleState"`
}

type SolverJoined struct {
	SolverMessage
	Solver uuid.UUID `json:"solver"`
}

type SolverLeft struct {
	SolverMessage
	Solver uuid.UUID `json:"solver"`
}
