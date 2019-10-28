package session

import (
	"encoding/json"
	"github.com/TenjouUtena/onedown/backend/src/onedown/cassandra"
	"github.com/TenjouUtena/onedown/backend/src/onedown/configuration"
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/wangjia184/sortedset"
	"time"
)

var nobody = uuid.MustParse("00000000-0000-0000-0000-000000000000")

type session struct {
	puzz        *puzzle.Puzzle
	channel     chan MessageForSession
	state       *PuzzleState
	solvers     map[uuid.UUID]*Solver
	initialized bool
}

func doSession(sessionId uuid.UUID, sesh *session) {
	// on session end, always write state
	defer sesh.writeSessionStateToCassandra(sessionId)
	lastWrite := time.Now() // variable to track time of last write. see end of for loop for usage

	// main session loop
	for msg := range sesh.channel {
		switch typedMsg := msg.Message.(type) {
		case JoinSession:
			sesh.broadcastSolverMessage(SolverJoined{
				Solver: typedMsg.Solver.Id,
			})
			sesh.solvers[typedMsg.Solver.Id] = typedMsg.Solver
			typedMsg.Solver.Tell(CurrentPuzzleState{
				Solvers:     sesh.getSolverIds(),
				Puzzle:      sesh.puzz,
				PuzzleState: sesh.state,
			})
		case LeaveSession:
			oldSolver := sesh.solvers[msg.solver]
			if oldSolver != nil {
				close(oldSolver.responseChannel) // unmarshal solver
				delete(sesh.solvers, msg.solver)
				sesh.broadcastSolverMessage(SolverLeft{
					Solver: msg.solver,
				})
			}
		case WriteSquare:
			sesh.state.putAnswer(msg.solver, typedMsg.Row, typedMsg.Col, typedMsg.Answer)
			sesh.broadcastSolverMessage(SquareUpdated{
				Row:      typedMsg.Row,
				Col:      typedMsg.Col,
				NewValue: typedMsg.Answer,
				FilledBy: msg.solver,
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
								Row:      rowIndex,
								Col:      colIndex,
								NewValue: square,
								FilledBy: nobody,
							})
						}
					}
				}
				sesh.broadcastSolverMessage(SquaresUpdated{Updates: updates})
			})
		default:
			log.Error().Msg("Invalid session message sent.")
		}
		if time.Now().After(lastWrite.Add(configuration.Get().PuzzleSessionWriteDelay)) {
			lastWrite = time.Now()
			go func() {
				time.Sleep(configuration.Get().PuzzleSessionWriteDelay)
				sesh.writeSessionStateToCassandra(sessionId)
			}()
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

func createSession(puzz *puzzle.Puzzle) (*session, uuid.UUID) {
	channel := make(chan MessageForSession)
	blankState := make([][]*sortedset.SortedSet, puzz.GetRowCount())
	for row := range blankState {
		blankState[row] = make([]*sortedset.SortedSet, puzz.GetColCount())
		for col := range blankState[row] {
			blankState[row][col] = sortedset.New()
		}
	}
	sessionObj := session{
		puzz:        puzz,
		channel:     channel,
		state:       &PuzzleState{
			filledSquares: blankState,
		},
		solvers:     make(map[uuid.UUID]*Solver),
		initialized: true,
	}
	sessionId := uuid.New()
	go doSession(sessionId, &sessionObj)
	return &sessionObj, sessionId
}

func (session* session) writeSessionStateToCassandra(sessionId uuid.UUID) {
	log.Debug().Str("sessionId", sessionId.String()).Msg(
		"Serializing session data.")
	jsonBlob, err := json.Marshal(*session)
	if err != nil {
		log.Error().Err(err).Str("sessionId", sessionId.String()).Msg(
			"Failed to marshal session data for serialization.")
	}
	cassandra.Session.Query("INSERT INTO puzzle_sessions VALUES (?, ?)",
		sessionId,
		jsonBlob,
	)
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
