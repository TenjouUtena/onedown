import React from 'react';
import './App.scss';
import { SessionNav } from './SessionNav';
import { AcrossClueList, DownClueList } from './Clue';
import { w3cwebsocket as W3CWebSocket } from "websocket";
import { Square } from './Square';

export var sqsize=40;
var curborder=2;

export var dirs = {
  ACROSS: 'across',
  DOWN: 'down'
}

class Selector extends React.Component {

  render() {
    const value = this.props.value;

    var style = {
        top: value.row*sqsize + curborder,
        left: value.col*sqsize + curborder
    }

    return (
      <div className="Selector" style={style} >

      </div>
    );
  }

}


class Game extends React.Component {
  
  state = {
    squares: [],
    acrossClues: {},
    downClues: {},
    gameMessage: "",
    puzzleInput: "Apr0914",
    selectorPos: {"row":0, "col":0},
    selectingAcross: true,
    acrossSelected: 1,
    downSelected: 1,
    width: 15,
    height: 15,
    session: ""
  }

  constructor() {
    super();

    this.client = null;
    this.sessnav = React.createRef();

    this.handleSquareClick = this.handleSquareClick.bind(this)
    this.handleClueClick = this.handleClueClick.bind(this)
    this.calcClueNums = this.calcClueNums.bind(this)

  }

  putGuess(row, col, guess) {
    let sqs = this.state.squares
    sqs.forEach((s) => {
      if(s.row === row && s.col === col) {
        s.answer = {guess: guess}
      }
    })

    this.setState({squares: sqs})
  }

  onClientMessage(message) {
    const m = JSON.parse(message.data)
    const p = m.payload
    if(m.name === "CurrentPuzzleState") {
      const puz = p.puzzle;
      const pw = this.calcClueNums(puz.squares)
      this.setState({width: puz.width,
                  height: puz.height,
                  acrossClues: puz.acrossClues,
                  downClues: puz.downClues,
                  squares: pw
      })
    }
  }

  buildws (url) {
    this.client = new W3CWebSocket(url)
    this.client.onmessage = (mess) => this.onClientMessage(mess);
    this.client.onopen = () => console.log("Connected to Session.")
    this.setState({session:this.sessnav.current.state.session})


  }
  
  findSquareFromArray (sqs,row,col) {
    return sqs.filter(e => (e.col === col && e.row === row))[0]
  }

  handleKeys (event) {
    //console.log(event.key)
    let guess = event.key.toUpperCase()[0]
    if(this.client) {
      var mess = {name:"WriteSquare",
                  session: this.state.session,
                  payload: JSON.stringify({
                    row: this.state.selectorPos.row,
                    col: this.state.selectorPos.col,
                    answer: guess
                  })}
      this.client.send(JSON.stringify(mess))
    }

    this.putGuess(this.state.selectorPos.row, this.state.selectorPos.col, guess)
  }
  
  calcClueNums(sqs) {
    sqs.forEach((e) => {
      if(e.clueNum > 0) {
        let aclue = e.clueNum;
        let a = e
        if(!a.acrossClue) {
          do {
            a.acrossClue = aclue;
            if(a.col+1 < this.state.width)
              a = this.findSquareFromArray(sqs,a.row,a.col+1)
          } while(!a.isBlack && a.col+1 < this.state.width)
          a.acrossClue = aclue;
        }
        let dclue = e.clueNum
        a=e
        if(!a.downClue) {
          do {
            a.downClue = dclue;
            if(a.row+1 < this.state.height)
              a = this.findSquareFromArray(sqs, a.row+1, a.col)
          } while(!a.isBlack && a.row+1 < this.state.height)
          a.downClue = dclue;
        }
      }
    })
    return sqs
  }

  handleClueClick (event, number, direction) {
    let s = this.findSquareFromClue(number)
    this.selectSquare(s.row, s.col)

    if(direction === dirs.ACROSS)
       this.setState({acrossSelected: number})
    if(direction === dirs.DOWN)
       this.setState({downSelected: number})
  }

  findSquareFromClue (clue) {
    return this.state.squares.filter(e => (e.clueNum == clue))[0]
  }

  findSquare (row,col) {
    return this.state.squares.filter(e => (e.col === col && e.row === row))[0]
  }

  selectSquare (row,col) {
    let sqs = this.state.squares;
    sqs = this.resetSelection(sqs);
    this.setState({squares: sqs})
    this.setState({selectorPos: {"row":row, "col":col}});

    let dsel = 0
    let asel = 0

    sqs.forEach(e => {
      if(e.row === row && e.col === col) {
        e.selected = true;
        asel = e.acrossClue;
        dsel = e.downClue;
      }
    })

    this.setState({squares: sqs,
                   acrossSelected: asel,
                   downSelected: dsel
    })


  }

  handleSquareClick (event, row, col) {

    let s = this.findSquare(row,col);

    if (!s.isBlack) {

      this.selectSquare(row,col);
    }
   
  }

  resetSelection (sqs) {
    sqs.forEach(element => {
      element.selected = false;
    });
    return sqs;
  }

  showSessionNav (event) {
    document.getElementsByClassName('SessionNav')[0].style.borderStyle='solid';
    document.getElementsByClassName('SessionNav')[0].style.height="800px";

  }


  componentDidMount() {
    this.selectSquare(0,0);
  }

  render () {
      const squares = this.state.squares;

      var astyle = {
        top: 50
      }

      var dstyle = {
        top: (40*6.5)+50,
      }

      var aval = {
        values: this.state.acrossClues,
        selected: this.state.acrossSelected
      }

      var dval = {
        values: this.state.downClues,
        selected: this.state.downSelected
      }

      return (
      <div className="Game">

        <SessionNav buildws={(e) => this.buildws(e)} session={this.state.session} ref={this.sessnav}/>

        <div>
          <span className="GameMessage">{this.state.gameMessage}</span>
        </div>
       <div className="Grid">
       {squares.map( (t, i) => {
          return (
          <Square value={this.state.squares[i]} key={(t.row*100)+t.col} onClick={this.handleSquareClick} onKeyPress={(e) => (this.handleKeys(e))}/> 
          );
       }
       )}
        <Selector value={this.state.selectorPos}></Selector>
        </div>

        <AcrossClueList className="AcrossClueList" value={aval} style={astyle} onClick={this.handleClueClick} />
        <DownClueList className="DownClueList" value={dval} style={dstyle} onClick={this.handleClueClick}/>
        <button className="SessButt" onClick={(e) => this.showSessionNav(e)}>Show Session</button>
       </div>

      );}
}


function App() {
  return (
      <Game />
  );
}

export default App;
