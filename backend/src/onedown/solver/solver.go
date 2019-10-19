package solver

import "github.com/google/uuid"

type Solver struct {
	Id              uuid.UUID
	responseChannel chan SolverMessage
}

func (solver *Solver) Tell(message SolverMessage) {
	solver.responseChannel <- message
}
