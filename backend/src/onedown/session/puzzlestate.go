package session

import (
	"github.com/google/uuid"
	"github.com/wangjia184/sortedset"
	"time"
)

type PuzzleState struct {
	filledSquares [][]*sortedset.SortedSet
}

func (state *PuzzleState) putAnswer(solver uuid.UUID, row int, col int, answer string) {
	now := sortedset.SCORE(time.Now().UnixNano())
	state.filledSquares[row][col].AddOrUpdate(solver.String(), now, answer)
}

func (state *PuzzleState) getSquare(row int, col int) string {
	var answer interface {}
	// if this square has been revealed, show the revealed square
	if state.filledSquares[row][col].GetByKey(nobody.String()) != nil {
		answer = state.filledSquares[row][col].GetByKey(nobody.String()).Value
	} else {
		answer = state.filledSquares[row][col].PeekMax().Value
	}
	switch strAnswer := answer.(type) {
	case string:
		return strAnswer
	default:
		return ""
	}
}
