package puzzle

import (
	"os"
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
	// read checksum data
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
	stringOffset := int64(0x34 + solutionLength)
	var stringLength int64 = 0
	stat, statErr := puzFile.Stat()
	if statErr == nil {
		stringLength = stat.Size() - stringOffset
	} else {
		err = statErr
	}

	stringBytes := make([]byte, stringLength)
	err = readTo(puzFile, stringBytes, stringOffset, err, func() {
		puzzleStrings := strings.Split(string(stringBytes), "\000")
		puzzfile.title = puzzleStrings[0]
		puzzfile.author = puzzleStrings[1]
		puzzfile.copyright = puzzleStrings[2]

		// clues
		clueCountBytes := make([]byte, 2)
		err = readTo(puzFile, clueCountBytes, 0x2E, err, func() {
			clueCount := clueCountBytes[0] + clueCountBytes[1]*0xFF
			puzzfile.clues = puzzleStrings[3 : 3+clueCount]
		})
	})

	// TODO: don't discard the rest of the data
	return puzzfile, err
}

func (puzzfile *Puzzlefile) WriteToPuzFile(puzFile *os.File) {
	// TODO: write
}

func (puzzfile *Puzzlefile) ToPuzzle() Puzzle {
	puzzle := Puzzle{}
	puzzle.squares = make([][]square, puzzfile.height)
	for index := 0; index < len(puzzfile.solution); index++ {
		currCol := index % len(puzzfile.solution)
		currRow := index / int(puzzfile.width)
		if currCol == 0 { // we are on a new row
			puzzle.squares[currRow] = make([]square, puzzfile.width)
		}
		if puzzfile.solution[index] != '.' {
			puzzle.squares[currRow][currCol] = square{
				correctValue: string(puzzfile.solution[index]),
			}
		}
	}

	puzzle.Clues = make([]Clue, len(puzzfile.clues))
	for index := 0; index < len(puzzfile.clues); index++ {

		// TODO: figure out the direction/length and correctly implement them.
		puzzle.Clues[index] = Clue{
			ClueText:      puzzfile.clues[index],
			ClueDirection: Across,
			Length:        0,
			puzzle:        &puzzle,
		}
	}

	return puzzle
}
