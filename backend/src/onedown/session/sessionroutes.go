package session

import (
	puzzle2 "github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func InitSessionRoutes(r *gin.Engine) {
	// create session with puzzle file
	r.GET("/session", func(context *gin.Context) {
		channel := make(chan []uuid.UUID)
		SessionDaemon <- GetSessions{
			ResponseChannel: channel,
		}
		uuids := <- channel
		context.JSON(200, uuids)
	})
	r.POST("/session", func(context *gin.Context) {
		// get puzzle file from request TODO: is there a better way to do this w/ gin?
		puzzFileHeader, err := context.FormFile("puzFile")
		if err != nil {
			context.JSON(500, gin.H{"error": err})
			return
		}
		puzzFile, err := puzzFileHeader.Open()
		if err != nil {
			context.JSON(500, gin.H{"error": err})
			return
		}
		defer puzzFile.Close()
		tmpFile, err := ioutil.TempFile("", "puzz")
		if err != nil {
			context.JSON(500, gin.H{"error": err})
			return
		}
		_, err = io.Copy(tmpFile, puzzFile)
		if err != nil {
			context.JSON(500, gin.H{"error": err})
			return
		}
		puzzle, err := puzzle2.ReadPuzfile(tmpFile)
		if err != nil {
			context.JSON(500, gin.H{"error": err})
			return
		}

		finalPuzz := puzzle.ToPuzzle()

		request := context.Request
		responseWriter := context.Writer
		sessionId := uuid.New()
		socket, err := upgrader.Upgrade(responseWriter, request, nil)
		if err != nil {
			context.JSON(500, gin.H{"error": err})
			return
		}
		solver := InitSolver(socket, sessionId)
		SessionDaemon <- SpawnSessionWithSolver{
			Puzzle: &finalPuzz,
			Solver: solver,
		}
	})
	// join session as solver
	r.POST("/session/:sessionId", func(context *gin.Context) {
		// build socket
		request := context.Request
		responseWriter := context.Writer
		sessionId, err := uuid.Parse(context.Param("sessionId"))
		if err != nil {
			context.JSON(500, gin.H{"error": err})
		} else {
			socket, err := upgrader.Upgrade(responseWriter, request, nil)
			if err != nil  {
				context.JSON(500, gin.H{"error": err})
			} else {
				InitSolver(socket, sessionId)
			}
		}
	})
}
