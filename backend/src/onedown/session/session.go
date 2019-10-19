package session

import (
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/TenjouUtena/onedown/backend/src/onedown/solver"
	"github.com/google/uuid"
)

type session struct {
	puzz        *puzzle.Puzzle
	channel     chan SessionMessage
	state       PuzzleState
	solvers     map[uuid.UUID]*solver.Solver
	initialized bool
}

func doSession(sesh *session) {
	for msg := range sesh.channel {
		switch typedMsg := msg.(type) {
		case JoinSession:
			sesh.broadcastSolverMessage(solver.SolverJoined{
				Solver: typedMsg.solver.Id,
			})
			sesh.solvers[typedMsg.solver.Id] = typedMsg.solver
			typedMsg.solver.Tell(solver.PuzzleState{
				Solvers:     sesh.getSolverIds(),
				Puzzle:      sesh.puzz,
				PuzzleState: &sesh.state,
			})
		case LeaveSession:
			delete(sesh.solvers, typedMsg.solver)
			sesh.broadcastSolverMessage(solver.SolverLeft{
				Solver: typedMsg.solver,
			})
			// TODO any unmarshalling of solver?
		case WriteSquare:
			sesh.state.putAnswer(typedMsg.solver, typedMsg.row, typedMsg.col, typedMsg.answer)
			sesh.broadcastSolverMessage(solver.SquareUpdated{
				Row:      typedMsg.row,
				Col:      typedMsg.col,
				NewValue: typedMsg.answer,
				FilledBy: typedMsg.solver,
			})
		case CheckSquares:
			if typedMsg.rowIndices[1] >= typedMsg.rowIndices[0] && typedMsg.colIndices[1] >= typedMsg.colIndices[0] {
				slice := make([][]string, typedMsg.rowIndices[1]-typedMsg.rowIndices[0]+1)
				for rowIndex := typedMsg.rowIndices[0]; rowIndex <= typedMsg.rowIndices[1]; rowIndex++ {
					slice[rowIndex] = make([]string, typedMsg.colIndices[0]-typedMsg.colIndices[1]+1)
					for colIndex := typedMsg.colIndices[0]; colIndex <= typedMsg.colIndices[1]; colIndex++ {
						slice[rowIndex][colIndex] = sesh.state.getSquare(rowIndex, colIndex)
					}
				}
				result := sesh.puzz.CheckSection(typedMsg.rowIndices[0], typedMsg.colIndices[0], slice)
				sesh.solvers[typedMsg.solver].Tell(solver.CheckResult{
					StartRow: typedMsg.rowIndices[0],
					StartCol: typedMsg.colIndices[0],
					Result:   result,
				})
			} else {
				// log.error("Bad indices on check message.")
			}
		default:
			// log.error("Invalid session message sent.")
		}
	}
}

func createSession(puzz *puzzle.Puzzle) *session {
	channel := make(chan SessionMessage)
	sessionObj := session{
		puzz:        puzz,
		channel:     channel,
		state:       PuzzleState{},
		initialized: true,
	}
	go doSession(&sessionObj)
	return &sessionObj
}

func (sesh *session) broadcastSolverMessage(message solver.SolverMessage) {
	for _, slv := range sesh.solvers {
		slv.Tell(message)
	}
}

func (sesh *session) getSolverIds() []uuid.UUID {
	result := make([]uuid.UUID, 0)
	for slvId, _ := range sesh.solvers {
		result = append(result, slvId)
	}
	return result
}
