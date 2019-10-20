package session

import (
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type session struct {
	puzz        *puzzle.Puzzle
	channel     chan SessionMessage
	state       PuzzleState
	solvers     map[uuid.UUID]*Solver
	initialized bool
}

func doSession(sesh *session) {
	for msg := range sesh.channel {
		switch typedMsg := msg.(type) {
		case JoinSession:
			sesh.broadcastSolverMessage(SolverJoined{
				Solver: typedMsg.Solver.Id,
			})
			sesh.solvers[typedMsg.Solver.Id] = typedMsg.Solver
			typedMsg.Solver.Tell(CurrentPuzzleState{
				Solvers:     sesh.getSolverIds(),
				Puzzle:      sesh.puzz,
				PuzzleState: &sesh.state,
			})
		case LeaveSession:
			oldSolver := sesh.solvers[typedMsg.Solver]
			if oldSolver != nil {
				close(oldSolver.responseChannel) // unmarshal solver
				delete(sesh.solvers, typedMsg.Solver)
				sesh.broadcastSolverMessage(SolverLeft{
					Solver: typedMsg.Solver,
				})
			}
		case WriteSquare:
			sesh.state.putAnswer(typedMsg.Solver, typedMsg.Row, typedMsg.Col, typedMsg.Answer)
			sesh.broadcastSolverMessage(SquareUpdated{
				Row:      typedMsg.Row,
				Col:      typedMsg.Col,
				NewValue: typedMsg.Answer,
				FilledBy: typedMsg.Solver,
			})
		case CheckSquares:
			if typedMsg.RowIndices[1] >= typedMsg.RowIndices[0] && typedMsg.ColIndices[1] >= typedMsg.ColIndices[0] {
				slice := make([][]string, typedMsg.RowIndices[1]-typedMsg.RowIndices[0]+1)
				for rowIndex := typedMsg.RowIndices[0]; rowIndex <= typedMsg.RowIndices[1]; rowIndex++ {
					slice[rowIndex] = make([]string, typedMsg.ColIndices[0]-typedMsg.ColIndices[1]+1)
					for colIndex := typedMsg.ColIndices[0]; colIndex <= typedMsg.ColIndices[1]; colIndex++ {
						slice[rowIndex][colIndex] = sesh.state.getSquare(rowIndex, colIndex)
					}
				}
				result := sesh.puzz.CheckSection(typedMsg.RowIndices[0], typedMsg.ColIndices[0], slice)
				sesh.solvers[typedMsg.Solver].Tell(CheckResult{
					StartRow: typedMsg.RowIndices[0],
					StartCol: typedMsg.ColIndices[0],
					Result:   result,
				})
			} else {
				log.Error().
					Int("rowLow", typedMsg.RowIndices[0]).
					Int("rowHigh", typedMsg.RowIndices[1]).
					Int("colLow", typedMsg.ColIndices[0]).
					Int("colHigh", typedMsg.ColIndices[1]).
					Msg("Bad indices on check message.")
			}
		default:
			log.Error().Msg("Invalid session message sent.")
		}
	}
}

func createSession(puzz *puzzle.Puzzle) *session {
	channel := make(chan SessionMessage)
	sessionObj := session{
		puzz:        puzz,
		channel:     channel,
		state:       PuzzleState{},
		solvers:     make(map[uuid.UUID]*Solver),
		initialized: true,
	}
	go doSession(&sessionObj)
	return &sessionObj
}

func (sesh *session) broadcastSolverMessage(message SolverMessage) {
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
