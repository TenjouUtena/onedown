package puzzle_test

import (
	"encoding/json"
	"os"
	"testing"
)

func getTestPuzzle(t *testing.T) Puzzle {
	file, err := os.Open("testdata/test_puzzle_1.puz")
	if err != nil {
		t.Fatal(err)
	}
	puzFile, err := ReadPuzfile(file)
	if err != nil {
		t.Fatal(err)
	}

	return puzFile.ToPuzzle()
}

// Test that parsing a file to Puzzle correctly collects all of the clues.
func TestPuzzlefile_ToPuzzle_parsesClues(t *testing.T) {
	testPuzzle := getTestPuzzle(t)
	if testPuzzle.AcrossClues[1].ClueText != "1 Across" || testPuzzle.AcrossClues[4].ClueText != "4 Across" ||
		testPuzzle.AcrossClues[5].ClueText != "5 Across" || testPuzzle.DownClues[1].ClueText != "1 Down" ||
		testPuzzle.DownClues[2].ClueText != "2 Down" || testPuzzle.DownClues[3].ClueText != "3 Down" {
		t.Fail()
	}
}

func TestPuzzlefile_ToPuzzleP_parsesCorrectAnswer(t *testing.T) {
	testPuzzle := getTestPuzzle(t)
	checkAnswers := [][]string{{"A", "A", "A"}, {"A", "A", "A"}, {"A", "A", "A"}}
	checkResults := testPuzzle.CheckSection(0, 0, checkAnswers)

	for _, row := range checkResults {
		for _, col := range row {
			if !col {
				t.Fail()
			}
		}
	}

}

func TestPuzzle_MarshalJSON(t *testing.T) {
	testPuzzle := getTestPuzzle(t)
	puzzJson, err := json.Marshal(testPuzzle)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"width":3,"height":3,"acrossClues":{"1":"1 Across","4":"4 Across","5":"5 Across"},"downClues":{"1":"1 Down","2":"2 Down","3":"3 Down"},"squares":[{"row":0,"col":0,"clueNum":1,"isBlack":false},{"row":0,"col":1,"clueNum":2,"isBlack":false},{"row":0,"col":2,"clueNum":3,"isBlack":false},{"row":1,"col":0,"clueNum":4,"isBlack":false},{"row":1,"col":1,"clueNum":0,"isBlack":false},{"row":1,"col":2,"clueNum":0,"isBlack":false},{"row":2,"col":0,"clueNum":5,"isBlack":false},{"row":2,"col":1,"clueNum":0,"isBlack":false},{"row":2,"col":2,"clueNum":0,"isBlack":false}]}`

	if string(puzzJson) != expected {
		t.Errorf("Generated JSON was incorrect:\n Expected: %s\n Actual:  %s", expected, string(puzzJson))
	}
}
