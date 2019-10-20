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
	switch answer := state.filledSquares[row][col].PeekMax().Value.(type) {
	case string:
		return answer
	default:
		return ""
	}
}
