package puzzle

type Puzzle struct {
	squares     [][]square
	AcrossClues map[int]*Clue
	DownClues   map[int]*Clue
	Metadata puzzleMeta
}

// CheckSection checks a set of squares with top left at fromRow, fromCol coordinates against the correct solution.
// This function assumes fromRow and fromCol are 0 indexed.
func (puzz *Puzzle) CheckSection(fromRow int, fromCol int, answers [][]string) [][]bool {
	result := make([][]bool, len(answers))
	for checkRow := 0; checkRow < len(answers); checkRow++ {
		result[checkRow] = make([]bool, len(answers[checkRow]))
		for checkCol := 0; checkCol < len(answers[checkRow]); checkCol++ {
			answerToCheck := answers[checkRow][checkCol]
			squareToCheck := puzz.squares[fromRow + checkRow][fromCol + checkCol]
			result[checkRow][checkCol] = squareToCheck.correctValue == answerToCheck
		}
	}
	return result
}

type puzzleMeta struct {
	Title string
	Author string
	Copyright string
	Notes string
}

type square struct {
	number       int
	correctValue string
}

func (sq square) isBlack() bool {
	return sq.correctValue == ""
}