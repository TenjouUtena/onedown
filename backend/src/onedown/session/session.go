package session

import (
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/google/uuid"
	"github.com/wangjia184/sortedset"
	"time"
)

type session struct {
	puzz        *puzzle.Puzzle
	channel     chan SessionMessage
	state       puzzleState
	solvers     []*uuid.UUID
	initialized bool
}

func doSession(sesh *session) {
	for msg := range sesh.channel {
		switch typedMsg := msg.(type) {
		case WriteSquare:
			sesh.state.putAnswer(typedMsg.player, typedMsg.row, typedMsg.col, typedMsg.answer)
		}
	}
}

func createSession(puzz *puzzle.Puzzle) *session {
	channel := make(chan SessionMessage)
	sessionObj := session{
		puzz:        puzz,
		channel:     channel,
		state:       puzzleState{},
		initialized: true,
	}
	go doSession(&sessionObj)
	return &sessionObj
}

type puzzleState struct {
	filledSquares [][]*sortedset.SortedSet
}

func (state *puzzleState) putAnswer(solver uuid.UUID, row int, col int, answer string) {
	now := sortedset.SCORE(time.Now().UnixNano())
	state.filledSquares[row][col].AddOrUpdate(solver.String(), now, answer)
}