package puzzle

import (
	"bytes"
	"encoding/gob"
	"github.com/TenjouUtena/onedown/backend/cassandra"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

var initialized bool = false

func init_puzzleserialization() {
	if !initialized {
		gob.Register(Puzzle{})
		initialized = true
	}
}

func (puzzle *Puzzle) WriteToCassandra() error {
	init_puzzleserialization()
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(*puzzle); err != nil {
		return err
	}
	err := cassandra.GetSession().Query("INSERT INTO puzzle (ID, PuzzleData) VALUES (?, ?);",
		puzzle.Id, buffer.Bytes()).Exec()
	return err
}

func GetPuzzleFromCassandra(puzzleID uuid.UUID) (*Puzzle, error) {
	init_puzzleserialization()
	puzzle := Puzzle{}
	m := map[string]interface{}{}
	iter := cassandra.GetSession().Query("SELECT PuzzleData FROM puzzle WHERE ID = ?;", puzzleID).
		Consistency(gocql.One).Iter()
	for iter.MapScan(m) {
		puzzBytes := m["PuzzleData"].([]byte)
		buffer := bytes.Buffer{}
		buffer.Write(puzzBytes)
		decoder := gob.NewDecoder(&buffer)
		err := decoder.Decode(puzzle)
		return &puzzle, err
	}
	return nil, gocql.ErrNotFound
}
