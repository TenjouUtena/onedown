package session

import (
	"encoding/json"
	"github.com/google/uuid"
)

func (puzzState PuzzleState) MarshalJSON() ([]byte, error) {
	squares := make([]jsonPuzzleStateSquare, 0)
	for rowIndex, row := range puzzState.filledSquares {
		for colIndex, col := range row {
			maybeString := (*col).PeekMax()
			if maybeString == nil {
				// no value has been entered
				continue
			}
			switch value := (*col).PeekMax().Value.(type) {
			case string:
				if value == "" {
					solver, err := uuid.Parse(col.PeekMax().Key())
					if err == nil {
						squares = append(squares, jsonPuzzleStateSquare{
							Row:      rowIndex,
							Col:      colIndex,
							Value:    value,
							FilledBy: solver,
						})
					}
				}
			}
		}
	}
	return json.Marshal(jsonPuzzleState{Squares: squares})
}

type jsonPuzzleState struct {
	Squares []jsonPuzzleStateSquare `json:"squares"`
}

type jsonPuzzleStateSquare struct {
	Row      int       `json:"row"`
	Col      int       `json:"col"`
	Value    string    `json:"value"`
	FilledBy uuid.UUID `json:"filledBy"`
}
