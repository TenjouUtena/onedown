package session

import "github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"

type session struct {
	puzz        *puzzle.Puzzle
	channel     chan SessionMessage
	state       puzzleState
	solvers     []*solver
	initialized bool
}

func doSession(sesh *session) {
	// TODO: impl
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
}

type solver struct {
}

type filledSquare struct {
}
