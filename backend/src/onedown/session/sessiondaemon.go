package session

import (
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/google/uuid"
)

var SessionDaemon = make(chan SessionDaemonMessage)

// InitDaemon initiates a puzzle session daemon. Should be invoked as a goroutine.
func InitDaemon(listen chan SessionDaemonMessage) {
	// session map lives for the duration of this goroutine. TODO: initialize w/ serialized DB data
	var sessions = make(map[uuid.UUID]*session)
	for msg := range listen {
		switch typedMsg := msg.(type) {
		case SpawnSessionWithSolver:
			newSession := createSession(typedMsg.Puzzle)
			sessionId := uuid.New()
			sessions[sessionId] = newSession
			newSession.channel <- JoinSession{
				Solver:         typedMsg.Solver,
			}
		case MessageForSession:
			sesh := sessions[typedMsg.SessionId]
			if sesh.initialized {
				sesh.channel <- typedMsg.Message
			} else {
				//log.Error("Message sent to session not on this daemon.")
			}
		case KillDaemon:
			// log.Info("Killing Session Daemon.")
			return
		default:
			// log.Error("Invalid message type sent to Session Daemon.")
		}
	}

}

type SessionDaemonMessage interface{}

type KillDaemon struct {
	SessionDaemonMessage
}

type SpawnSessionWithSolver struct {
	SessionDaemonMessage
	Solver *Solver
	Puzzle *puzzle.Puzzle
}

type MessageForSession struct {
	SessionDaemonMessage
	SessionId uuid.UUID
	Message   SessionMessage
}
