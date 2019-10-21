package session

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"reflect"
)
import "github.com/gorilla/websocket"

type Solver struct {
	Id              uuid.UUID
	responseChannel chan SolverMessage
	socket          *websocket.Conn
}

type serverMessage struct {
	MessageName string        `json:"name"`
	Payload     SolverMessage `json:"payload"`
}

func doSolverChannel(solver *Solver) {
	for msg := range solver.responseChannel {
		jsonMsg := serverMessage{
			MessageName: reflect.TypeOf(msg).Name(),
			Payload:     msg,
		}
		err := solver.socket.WriteJSON(jsonMsg)
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
		messageType, messageBytes, err := solver.socket.ReadMessage()
		if err != nil {
			// TODO At the moment this will error when a user leaves a session. need to figure that one out
			log.Error().Err(err).Str("solverId", solver.Id.String()).Msg("Error when reading message from socket.")
			break
		} else if messageType != websocket.TextMessage {
			log.Error().Str("solverId", solver.Id.String()).Msg("Message sent as binary, not currently supported.")
		} else {
			// Write it to daemon to be delegated to the appropriate session
			unmarshalledMessage, err := unmarshallSocketMessage(messageBytes)
			if err != nil {
				log.Error().
					Err(err).
					Str("solverId", solver.Id.String()).
					Bytes("payload", messageBytes).
					Msg("Bad message payload.")
			} else {
				SessionDaemon <- unmarshalledMessage
			}
		}
	}
	// TODO: upon socket closing, we should kick the user out of the puzzle
}

// We are expecting messages in this JSON format, with payload being the *fields* of the struct of type `name`:
/*  {
		"name":"SomeMessageType",
		"session":"00000000-0000-0000-0000-000000000000",
		"payload":"..."
	}
 */
func unmarshallSocketMessage(messagePayload []byte) (MessageForSession, error) {
	stringMapPayload := make(map[string]string)
	err := json.Unmarshal(messagePayload, &stringMapPayload)
	if err != nil {
		return MessageForSession{}, err
	}
	typeName := stringMapPayload["name"]
	sessionId, err := uuid.Parse(stringMapPayload["session"])
	if err != nil {
		return MessageForSession{}, err
	}
	payload := []byte(stringMapPayload["payload"])
	messageSkeleton := sessionClientMessages[typeName]
	msg, err := messageSkeleton.unmarshalClientPayload(payload)
	if err != nil {
		return MessageForSession{}, err
	}
	return MessageForSession{
		SessionId: sessionId,
		Message:   msg,
	}, nil
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
