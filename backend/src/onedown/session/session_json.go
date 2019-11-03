package session

import (
	"encoding/json"
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/google/uuid"
)

func (sesh session) MarshalJSON() ([]byte, error) {
	solvers := make([]uuid.UUID, 0)
	for solverId, _ := range sesh.solvers {
		solvers = append(solvers, solverId)
	}
	jsonMsg := jsonSession{
		PuzzleId: sesh.puzz.Id,
		State:    sesh.state,
		Solvers:  solvers,
	}
	return json.Marshal(jsonMsg)
}

func (sesh session) UnmarshalJSON(data []byte) error {
	jsonSesh := jsonSession{}
	if err := json.Unmarshal(data, &jsonSesh); err != nil {
		return err
	}

	if puzz, err := puzzle.GetPuzzleFromCassandra(sesh.puzz.Id); err != nil {
		return err
	} else {
		sesh.puzz = puzz
	}
	solvers := make(map[uuid.UUID]*Solver)
	for _, solver := range jsonSesh.Solvers {
		// nil solvers on load, will be added as they connect
		solvers[solver] = nil
	}
	sesh.channel = make(chan MessageForSession)
	sesh.state = jsonSesh.State
	sesh.solvers = solvers
	return nil
}

type jsonSession struct {
	PuzzleId uuid.UUID    `json:"puzzleId"`
	State    *PuzzleState `json:"state"`
	Solvers  []uuid.UUID  `json:"solvers"`
}
