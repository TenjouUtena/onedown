package session

import (
	"github.com/google/uuid"
)
import "github.com/gorilla/websocket"

type Solver struct {
	Id              uuid.UUID
	responseChannel chan SolverMessage
	socket          *websocket.Conn
}

func doSolverChannel(solver *Solver) {
	for msg := range solver.responseChannel {
		err := solver.socket.WriteJSON(msg)
		if (err != nil) {
			// log.error(err)
		}
	}
}

func doSolverSocket(solver *Solver) {
	for {
		request := MessageForSession{}
		err := solver.socket.ReadJSON(request)
		if (err != nil) {
			// log.error(err)
		} else {
			// Write it to daemon to be delegated to the appropriate session
			SessionDaemon <- request.Message
		}
	}
}

func InitSolver(socket *websocket.Conn, sessionId uuid.UUID) *Solver {
	id := uuid.New()
	channel := make(chan SolverMessage)
	solver := &Solver{
		Id:              id,
		responseChannel: channel,
		socket:          socket,
	}
	// spawn child goroutines to handle input from socket and output from channel
	go doSolverChannel(solver)
	go doSolverSocket(solver)
	// join the session
	SessionDaemon <- MessageForSession{
		SessionId: sessionId,
		Message: JoinSession{
			Solver: solver,
		},
	}
	return solver
}

func (solver *Solver) Tell(message SolverMessage) {
	solver.responseChannel <- message
}
