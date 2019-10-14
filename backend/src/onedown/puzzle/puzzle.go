package puzzle

type Puzzle struct {
	squares [][]square
	Clues   []Clue
}

type square struct {
	correctValue string
}