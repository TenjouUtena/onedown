package puzzle

type Puzzle struct {
	squares     [][]square
	AcrossClues map[int]Clue
	DownClues   map[int]Clue
	Metadata puzzleMeta
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
