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
	DisplayName		string
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
	log.Debug().Str("solverId", solver.Id.String()).Msg("Solver channel closed.")
}

func doSolverSocket(solver *Solver) {
	defer func() {
		err := solver.socket.Close()
		if err != nil {
			log.Error().Err(err).Str("solverId", solver.Id.String()).Msg("Error closing socket.")
		}
	}()

	socketLive := true
	for socketLive {
		messageType, messageBytes, err := solver.socket.ReadMessage()
		if err != nil {
			switch errType := err.(type) {
			case *websocket.CloseError:
				// Per https://github.com/Luka967/websocket-close-codes these are expected closes.
				if errType.Code == 1000 || errType.Code == 1001 {
					log.Info().Str("solverId", solver.Id.String()).Msg("Solver socket closed, doing cleanup.")
				} else {
					log.Error().Err(err).Str("solverId", solver.Id.String()).Msg("Unexpected socket close.")
				}
			default:
				log.Error().Err(err).Str("solverId", solver.Id.String()).Msg("Error when reading message from socket.")
			}
			break
		} else if messageType != websocket.TextMessage {
			log.Error().Str("solverId", solver.Id.String()).Msg("Message sent as binary, not currently supported.")
		} else if string(messageBytes) == "PING" {
			err := solver.socket.WriteMessage(websocket.TextMessage, []byte("PONG"))
			if err != nil {
				log.Error().Err(err).Msg("Error on pong!")
			}
		} else {
			// Write it to daemon to be delegated to the appropriate session
			unmarshalledMessage, err := solver.unmarshallSocketMessage(messageBytes)
			if err != nil {
				log.Error().
					Err(err).
					Str("solverId", solver.Id.String()).
					Bytes("payload", messageBytes).
					Msg("Bad message payload.")
			} else {
				switch unmarshalledMessage.Message.(type) {
				case LeaveSession:
					log.Info().Str("solverId", solver.Id.String()).Msg("Solver is leaving.")
					socketLive = false // break out of loop
				default:
					SessionDaemon <- unmarshalledMessage
				}
			}
		}
	}

	// Inform Daemon user has left
	SessionDaemon <- UserDisconnected{Solver: solver.Id}
}

// We are expecting messages in this JSON format, with payload being the *fields* of the struct of type `name`:
/*
{
	"name":"SomeMessageType",
	"session":"00000000-0000-0000-0000-000000000000",
	"payload":"..."
}
*/
func (solver *Solver) unmarshallSocketMessage(messagePayload []byte) (MessageForSession, error) {
	stringMapPayload := make(map[string]string)
	err := json.Unmarshal(messagePayload, &stringMapPayload)
	if err != nil {
		return MessageForSession{}, err
	}
	typeName := stringMapPayload["name"]
	session, err := uuid.Parse(stringMapPayload["session"])
	if err != nil {
		return MessageForSession{}, err
	}
	payload := []byte(stringMapPayload["payload"])
	messageSkeleton := sessionClientMessages[typeName]
	msg, err := messageSkeleton.unmarshalClientPayload(payload)
	if err != nil {
		return MessageForSession{}, err
	}
	return NewMessageForSession(solver.Id, session, msg), nil
}

func InitSolver(socket *websocket.Conn, session uuid.UUID) *Solver {
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
		Session: session,
		Message: JoinSession{
			Solver: solver,
		},
	}
	return solver
}

func (solver *Solver) Tell(message SolverMessage) {
	solver.responseChannel <- message
}
