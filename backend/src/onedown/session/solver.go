package session

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
		if err != nil {
			log.Error().Err(err).Str("solverId", solver.Id.String()).Msg("Error when writing message to socket.")
		}
	}
	log.Debug().Str("solverId", solver.Id.String()).Msg("Solver channel closed, closing socket.")
	err := solver.socket.Close()
	if err != nil {
		log.Error().Err(err).Str("solverId", solver.Id.String()).Msg("Error closing socket.")
	}
}

func doSolverSocket(solver *Solver) {
	for {
		request := MessageForSession{}
		err := solver.socket.ReadJSON(request)
		if err != nil {
			// TODO At the moment this will error when a user leaves a session. need to figure that one out
			log.Error().Err(err).Str("solverId", solver.Id.String()).Msg("Error when reading message from socket.")
			break
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
