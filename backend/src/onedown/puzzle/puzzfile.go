package puzzle

import (
	"os"
	"sort"
	"strings"
)

type Puzzlefile struct {
	checksum          [2]byte
	cibChecksum       [2]byte
	maskLowChecksum   [4]byte
	maskHiChecksum    [4]byte
	version           string
	reserved1c        [2]byte
	scrambledChecksum [2]byte
	reserved20        [0xC]byte
	width             uint8
	height            uint8
	bitmask           [2]byte
	scrambledTag      [2]byte
	solution          string
	title             string
	author            string
	copyright         string
	clues             []string
	notes             string
}

func readTo(file *os.File, targetArray []byte, offset int64, lastError error, andThen func()) (error) {
	if lastError != nil {
		return lastError
	}
	_, err := file.ReadAt(targetArray, offset)
	if err == nil {
		andThen()
	}
	return err
}

func ReadPuzfile(puzFile *os.File) (Puzzlefile, error) {
	puzzfile := Puzzlefile{}
	var err error
	// read checksum testdata
	err = readTo(puzFile, puzzfile.checksum[0:2], 0x00, err, func() {})
	err = readTo(puzFile, puzzfile.cibChecksum[0:2], 0x0E, err, func() {})
	err = readTo(puzFile, puzzfile.maskLowChecksum[0:4], 0x10, err, func() {})
	err = readTo(puzFile, puzzfile.maskHiChecksum[0:4], 0x14, err, func() {})
	err = readTo(puzFile, puzzfile.scrambledChecksum[0:2], 0x1E, err, func() {})

	// read puzzle metadata
	versionBytes := make([]byte, 2)
	err = readTo(puzFile, versionBytes, 0x10, err, func() {
		puzzfile.version = string(versionBytes)
	})

	widthBytes := make([]byte, 1)
	err = readTo(puzFile, widthBytes, 0x2C, err, func() {
		puzzfile.width = widthBytes[0]
	})

	heightBytes := make([]byte, 1)
	err = readTo(puzFile, heightBytes, 0x2C, err, func() {
		puzzfile.height = heightBytes[0]
	})

	// read reserved space
	err = readTo(puzFile, puzzfile.reserved1c[0:2], 0x1C, err, func() {})
	err = readTo(puzFile, puzzfile.reserved20[0:0xC], 0x1E, err, func() {})

	// read remaining non-strings
	err = readTo(puzFile, puzzfile.bitmask[0:2], 0x30, err, func() {})
	err = readTo(puzFile, puzzfile.scrambledTag[0:2], 0x32, err, func() {})

	// read solution
	solutionLength := puzzfile.width * puzzfile.height
	solnBytes := make([]byte, solutionLength)
	err = readTo(puzFile, solnBytes, 0x34, err, func() {
		puzzfile.solution = string(solnBytes)
	})

	// read puzzle strings
	stringOffset := 0x34 + int64(solutionLength) + int64(solutionLength)
	var stringLength int64 = 0
	stat, statErr := puzFile.Stat()
	if statErr == nil {
		stringLength = int64(stat.Size()) - stringOffset
	} else {
		err = statErr
	}

	stringBytes := make([]byte, stringLength)
	err = readTo(puzFile, stringBytes, int64(stringOffset), err, func() {
		puzzleStrings := strings.Split(string(stringBytes), "\000")

		// clues
		clueCountBytes := make([]byte, 2)
		err = readTo(puzFile, clueCountBytes, 0x2E, err, func() {
			clueCount := clueCountBytes[0] + clueCountBytes[1]*0xFF
			puzzfile.clues = puzzleStrings[3 : 3+clueCount]
		})
	})

	// TODO: don't discard the rest of the testdata
	return puzzfile, err
}

func (puzzfile *Puzzlefile) WriteToPuzFile(puzFile *os.File) {
	// TODO: write
}

func (puzzfile *Puzzlefile) ToPuzzle() Puzzle {
	puzzle := Puzzle{
		Metadata: puzzleMeta{
			Title:     puzzfile.title,
			Author:    puzzfile.author,
			Copyright: puzzfile.copyright,
			Notes:     puzzfile.notes,
		},
	}
	puzzle.squares = make([][]square, puzzfile.height)
	puzzle.AcrossClues = make(map[int]*Clue)
	puzzle.DownClues = make(map[int]*Clue)
	var squareNumber int = 1

	// first off, we will build the array of squares in the puzzle. we have the solved puzzle in the file, and
	// we will use this to construct the square objects. at the same time, we need to identify where clues fit in the
	// puzzle so we can build them later.
	for index := 0; index < len(puzzfile.solution); index++ {
		currCol := index % int(puzzfile.width)
		currRow := index / int(puzzfile.width)
		if currCol == 0 { // we are on a new row
			puzzle.squares[currRow] = make([]square, puzzfile.width)
		}
		if puzzfile.solution[index] != '.' {
			puzzle.squares[currRow][currCol] = square{
				correctValue: string(puzzfile.solution[index]),
			}
			// if either the above square or the square to the left is null, or off the bounds. we need to place a number
			// and note a clue
			addedClue := false
			if currCol == 0 || puzzle.squares[currRow][currCol-1].correctValue == "" {
				addedClue = true
				puzzle.AcrossClues[squareNumber] = &Clue{
					puzzle: &puzzle,
				}
			}
			if currRow == 0 || puzzle.squares[currRow-1][currCol].correctValue == "" {
				addedClue = true
				puzzle.DownClues[squareNumber] = &Clue{
					puzzle: &puzzle,
				}
			}

			// bump clue number if we added something
			if addedClue {
				puzzle.squares[currRow][currCol].number = squareNumber
				squareNumber++
			}
		}
	}

	// get array indices of all of the clues we generated so we can add their strings below
	acrossClueIndices := make([]int, 0)
	for k, _ := range puzzle.AcrossClues {
		acrossClueIndices = append(acrossClueIndices, k)
	}
	sort.Ints(acrossClueIndices)

	downClueIndices := make([]int, 0)
	for k, _ := range puzzle.DownClues {
		downClueIndices = append(downClueIndices, k)
	}
	sort.Ints(downClueIndices)

	for index := 0; index < len(puzzfile.clues); index++ {
		if index < len(acrossClueIndices) {
			puzzle.AcrossClues[acrossClueIndices[index]].setText(puzzfile.clues[index])
		} else {
			downIndex := downClueIndices[index-len(acrossClueIndices)]
			puzzle.DownClues[downIndex].setText(puzzfile.clues[index])
		}
	}

	return puzzle
}
