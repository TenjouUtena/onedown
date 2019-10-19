package session

import "github.com/google/uuid"

type SessionMessage interface {}

type WriteSquare struct {
	SessionMessage
	player uuid.UUID
	row int
	col int
	answer string
}