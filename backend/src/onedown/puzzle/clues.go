package puzzle

// Clue type represents a clue in the puzzle. It references the puzzle its parented to allow it to check to see if
// it has been filled correctly. The ClueText variable represents the text to be printed to aid solving the clue, and
// Length represents the number of squares. Note that the clue is not cognizant of its answer here, and depends on the
// puzzle to determine if it is correctly answered.
type Clue struct {
	ClueText      string
	ClueDirection Direction
	Length        int
	puzzle        *Puzzle
}

// Direction enum indicates whether a clue goes down or across.
type Direction int

const (
	Across Direction = 0
	Down   Direction = 1
)
