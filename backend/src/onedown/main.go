
package main

import "github.com/gin-gonic/gin"
import "github.com/gin-contrib/cors"

type Square struct{
  X int
  Y int
  Black bool
  AcrossClue int
  DownClue int
  DrawAcross bool
  DrawDown bool
}


func genpuz() []Square {

  var r []Square
  for i:=0; i<10; i++ {
    for j:=0; j<10; j++ {
      ss := Square{i, j, false, -1, -1, false, false}
      if ss.X==0 && ss.Y==0 {
        ss.DrawAcross = true
        ss.AcrossClue = 1
        ss.DownClue = 1
      }

      if ss.X==1 && ss.Y==0 {
        ss.DrawDown = true
        ss.AcrossClue = 1
        ss.DownClue = 2
      }
      if ss.X==2 && ss.Y==0 {
        ss.DrawDown = true
        ss.AcrossClue = 1
        ss.DownClue = 3
      }
      if ss.X==3 && ss.Y==0 {
        ss.DrawDown = true
        ss.AcrossClue = 1
        ss.DownClue = 4
      }

      if ss.X==4 && ss.Y==0 {
        ss.Black = true
      }


      r = append(r, ss)
    }
  }

  return r
}


func main() {
	r := gin.Default()

	r.Use(cors.Default())  // Needed to allow all API origins.
	r.GET("/puzzle/:puzid/get", func(c *gin.Context) {
	  c.JSON(200, genpuz())
  })
	r.Run() // listen and serve on 0.0.0.0:8080
}




