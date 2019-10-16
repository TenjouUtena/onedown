package puzzle

import "encoding/json"

func (puzz Puzzle) MarshalJSON() ([]byte, error) {
	acrossClues := make(map[int]string)
	for number, clue := range puzz.AcrossClues {
		acrossClues[number] = clue.ClueText
	}

	downClues := make(map[int]string)
	for number, clue := range puzz.DownClues {
		downClues[number] = clue.ClueText
	}

	squares := make([]jsonSquare, 0)
	for row := 0; row < len(puzz.squares); row++ {
		for col := 0; col < len(puzz.squares[row]); col++ {
			thisSquare := puzz.squares[row][col]
			newSquare := jsonSquare{
				Row:     row,
				Col:     col,
				ClueNum: thisSquare.number,
				Black:   thisSquare.isBlack(),
			}
			squares = append(squares, newSquare)
		}
	}

	jsonPuz := jsonPuzzle{
		AcrossClues: acrossClues,
		DownClues:   downClues,
		Squares:     squares,
	}
	return json.Marshal(jsonPuz)
}

/* Used to marshall a puzzle to the following format:
{
	acrossClues: {
		1: "Lorem ipsum",
		...
	},
	downClues: {
		2: "Dolor sit amet",
		...
	},
	squares: [
		{
			row: 0,
			col: 0,
			acrossClue: 1
			downClue: 1
			isBlack: false
		},
		{
			row: 0,
			col: 1,
			isBlack: false
		},
		{
			row: 0,
			col: 2,
			isBlack: true
		}
		...
	]
}
*/
type jsonPuzzle struct {
	AcrossClues map[int]string `json:"acrossClues"`
	DownClues   map[int]string `json:"downClues"`
	Squares     []jsonSquare   `json:"squares"`
}

type jsonSquare struct {
	Row     int  `json:"row"`
	Col     int  `json:"col"`
	ClueNum int  `json:"clueNum"`
	Black   bool `json:"isBlack"`
}
