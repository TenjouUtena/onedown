package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/TenjouUtena/onedown/backend/src/onedown/cassandra"
	"github.com/TenjouUtena/onedown/backend/src/onedown/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Square struct {
	X          int
	Y          int
	Black      bool
	AcrossClue int
	DownClue   int
	DrawAcross bool
	DrawDown   bool
}

type Puzzle struct {
	Sqs    []Square
	Height int
	Width  int
	//Clues?

}

func findsquare(X int, Y int, Sqs []Square) (Square, error) {
	for _, sq := range Sqs {
		if sq.X == X && sq.Y == Y {
			return sq, nil
		}
	}

	return Square{-1, -1, false, -1, -1, false, false}, errors.New("Could not find requested Square")
}

func isblank(X int, Y int, Sqs []Square) bool {
	f, e := findsquare(X, Y, Sqs)
	if e != nil {
		panic(e)
	}

	return f.Black
}

func calcAcrossClues(p Puzzle) Puzzle {
	curclue := 0
	var ts []Square

	for j := 0; j < p.Height; j++ {
		for i := 0; i < p.Width; i++ {
			ss, e := findsquare(i, j, p.Sqs)
			if e != nil {
				panic(e)
			}
			if !ss.Black {
				if i == 0 {
					curclue++
					ss.AcrossClue = curclue
					ss.DrawAcross = true
					//fmt.Printf("ss:%v\n", ss)
				} else {
					tt, ee := findsquare(i-1, j, p.Sqs)
					if ee != nil {
						panic(ee)
					}
					if tt.Black {
						curclue++
						ss.AcrossClue = curclue
						ss.DrawAcross = true
					} else {
						if j == 0 {
							curclue++
							ss.DownClue = curclue
							ss.DrawDown = true
						} else {
							tt, ee := findsquare(i, j-1, p.Sqs)
							if ee != nil {
								panic(ee)
							} else {
								if tt.Black {
									curclue++
									ss.DownClue = curclue
									ss.DrawDown = true
								}
							}
						}
					}
				}
			}
			ts = append(ts, ss)
		}
	}
	p.Sqs = ts
	return p
}

func calcclues(p Puzzle) Puzzle {

	p = calcAcrossClues(p)
	return p
}

func readpuz(file string) (Puzzle, error) {

	var r []Square
	var p Puzzle
	wid := make([]byte, 1)
	hei := make([]byte, 1)

	f, err := os.Open(file)
	if err != nil {

		return p, fmt.Errorf("Error with file: %v", err)
	}

	// Seek to Width
	f.Seek(0x2C, 0)
	f.Read(wid)
	f.Read(hei)

	//fmt.Printf("Width: %d Height: %d", wid[0], hei[0])

	width := int(wid[0])
	height := int(hei[0])

	//Seek to blankpuz
	f.Seek(int64(0x34+(width*height)), 0)
	br := make([]byte, 1)

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			bk := false
			f.Read(br)
			if rune(br[0]) == '.' {
				bk = true
			}
			ss := Square{i, j, bk, -1, -1, false, false}
			r = append(r, ss)
		}
	}

	p.Width = width
	p.Height = height
	p.Sqs = r

	return p, nil

}

func genpuz() []Square {

	var r []Square
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			ss := Square{i, j, false, -1, -1, false, false}
			if ss.X == 0 && ss.Y == 0 {
				ss.DrawAcross = true
				ss.AcrossClue = 1
				ss.DownClue = 1
			}

			if ss.X == 1 && ss.Y == 0 {
				ss.DrawDown = true
				ss.AcrossClue = 1
				ss.DownClue = 2
			}
			if ss.X == 2 && ss.Y == 0 {
				ss.DrawDown = true
				ss.AcrossClue = 1
				ss.DownClue = 3
			}
			if ss.X == 3 && ss.Y == 0 {
				ss.DrawDown = true
				ss.AcrossClue = 1
				ss.DownClue = 4
			}

			if ss.X == 4 && ss.Y == 0 {
				ss.Black = true
			}

			r = append(r, ss)
		}
	}

	return r
}

func main() {
	var cfg Configuration
	cfg.GOPATH = os.Getenv("GOPATH")

	configPath := flag.String("config", path.Join(cfg.GOPATH, "configuration"), "Path to configuration files")
	flag.Parse()
	file, err := os.Open(path.Join(*configPath, "configuration.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cassandraSession := cassandra.Session
	defer cassandraSession.Close()

	router := gin.Default()

	router.Use(cors.Default()) // Needed to allow all API origins.
	router.GET("/puzzle/:puzid/get", func(c *gin.Context) {
		finalPath := path.Join(cfg.PuzzleDirectory, c.Param("puzid")+".puz")
		p, err := readpuz(finalPath)

		if err != nil {
			c.JSON(500, gin.H{"Error": err})
		} else {
			clues := calcclues(p)
			c.JSON(200, clues.Sqs)
		}
	})
	router.POST("/users/new", users.Post)
	router.Run("127.0.0.1:" + strconv.Itoa(cfg.Port))
}
