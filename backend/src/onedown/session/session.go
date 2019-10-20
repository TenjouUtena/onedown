package session

import (
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var nobody = uuid.MustParse("00000000-0000-0000-0000-000000000000")

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
			ifValidIndices(typedMsg.RowIndices, typedMsg.ColIndices, func() {
				slice := make([][]string, typedMsg.RowIndices[1]-typedMsg.RowIndices[0]+1)
				for rowIndex := typedMsg.RowIndices[0]; rowIndex <= typedMsg.RowIndices[1]; rowIndex++ {
					slice[rowIndex] = make([]string, typedMsg.ColIndices[0]-typedMsg.ColIndices[1]+1)
					for colIndex := typedMsg.ColIndices[0]; colIndex <= typedMsg.ColIndices[1]; colIndex++ {
						slice[rowIndex][colIndex] = sesh.state.getSquare(rowIndex, colIndex)
					}
				}
				result := sesh.puzz.CheckSection(typedMsg.RowIndices[0], typedMsg.ColIndices[0], slice)
				sesh.broadcastSolverMessage(CheckResult{
					StartRow: typedMsg.RowIndices[0],
					StartCol: typedMsg.ColIndices[0],
					Result:   result,
				})
			})
		case RevealSquares:
			ifValidIndices(typedMsg.RowIndices, typedMsg.ColIndices, func() {
				updates := make([]SquareUpdated, 0)
				solutions := sesh.puzz.GetSolutions(typedMsg.RowIndices, typedMsg.ColIndices)
				for rowIndex, row := range solutions {
					for colIndex, square := range row {
						if sesh.state.getSquare(rowIndex, colIndex) != square {
							sesh.state.putAnswer(nobody, rowIndex, colIndex, square)
							updates = append(updates, SquareUpdated{
								Row:           rowIndex,
								Col:           colIndex,
								NewValue:      square,
								FilledBy:      nobody,
							})
						}
					}
				}
				sesh.broadcastSolverMessage(SquaresUpdated{Updates: updates})
			})
		default:
			log.Error().Msg("Invalid session message sent.")
		}
	}
}

func ifValidIndices(rowIndices [2]int, colIndices [2]int, thenDo func()) {
	if rowIndices[1] >= rowIndices[0] && colIndices[1] >= colIndices[0] {
		thenDo()
	} else {
		log.Error().
			Int("rowLow", rowIndices[0]).
			Int("rowHigh", rowIndices[1]).
			Int("colLow", colIndices[0]).
			Int("colHigh", colIndices[1]).
			Msg("Bad indices on message.")

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
