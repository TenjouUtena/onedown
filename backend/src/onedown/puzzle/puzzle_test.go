package puzzle_test

import (
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"os"
	"testing"
)

// Test that parsing a file to Puzzle correctly collects all of the clues.
func TestPuzzlefile_ToPuzzle_parsesClues(t *testing.T) {
	file, err := os.Open("testdata/test_puzzle_1.puz")
	if err != nil {
		t.Fatal(err)
	}
	puzFile, err := puzzle.ReadPuzfile(file)
	if err != nil {
		t.Fatal(err)
	}

	testPuzzle := puzFile.ToPuzzle()
	if testPuzzle.AcrossClues[1].ClueText != "1 Across" || testPuzzle.AcrossClues[4].ClueText != "4 Across" ||
		testPuzzle.AcrossClues[5].ClueText != "5 Across" || testPuzzle.DownClues[1].ClueText != "1 Down" ||
		testPuzzle.DownClues[2].ClueText != "2 Down" || testPuzzle.DownClues[3].ClueText != "3 Down" {
		t.Fail()
	}
}

func TestPuzzlefile_ToPuzzleP_parsesCorrectAnswer(t *testing.T) {
	file, err := os.Open("testdata/test_puzzle_1.puz")
	if err != nil {
		t.Fatal(err)
	}
	puzFile, err := puzzle.ReadPuzfile(file)
	if err != nil {
		t.Fatal(err)
	}

	testPuzzle := puzFile.ToPuzzle()
	checkAnswers := [][]string{{"A", "A", "A"}, {"A", "A", "A"}, {"A", "A", "A"}}
	checkResults := testPuzzle.CheckSection(0,0, checkAnswers)

	for _, row := range checkResults {
		for _, col := range row {
			if !col {
				t.Fail()
			}
		}
	}

}