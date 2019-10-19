import React from 'react';
import './App.scss';

var sqsize=40;
var curborder=2;

var dirs = {
  ACROSS: 'across',
  DOWN: 'down'
}

class Clue extends React.Component {

  render() {
    const value = this.props.value;

    let text = String(value.number) + ". " + value.text

    return <div className="Clue" selstyle={value.selected ? 'selected' : 'not-selected'} onClick={(e) => this.props.onClick(e,value.number, value.dir)}>{text}</div>
  }
}


class ClueList extends React.Component {
  dir = dirs.ACROSS;
  cname = "AcrossList";

  render() {
    const value = this.props.value;


    if (value) {
            return (
        <div className={this.cname} style={this.props.style}>        {
          Object.keys(value.values).map((k) => {
            let vv = {
              number: k,
              text: value.values[k],
              dir: this.dir,
              selected: value.selected == k
            }
          return (<Clue value={vv} key={k} onClick={this.props.onClick}/>);
          })
                        }
        </div>
      );}
    else {
      return <div />
    }
  }
}


class AcrossClueList extends ClueList {}

class DownClueList extends ClueList {
  dir = dirs.DOWN
  cname = "DownList"
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


class Square extends React.Component {
  render() {
    const value = this.props.value;
    var style = {
       top: value.row*sqsize,
       left: value.col*sqsize
    };
    let clue;
    if (value.clueNum) {
      clue = <span className="SquareClue" >{value.clueNum}</span>
    }

    let black;
    if (value.isBlack) {
      black = "black"
    } else {
      black = "white"
    }
    return (
      <div className="Square" style={style} id={black} onClick={(e) => this.props.onClick(e, this.props.value.row, this.props.value.col)}
                              selstyle={this.props.value.selected ? 'selected' : 'notSelected'}>
        {clue}
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
    height: 15
  }

  constructor() {
    super();

    this.handlePuzzleInputChange = this.handlePuzzleInputChange.bind(this)
    this.loadPuzzle = this.loadPuzzle.bind(this)
    this.handlePuzzleLoad = this.handlePuzzleLoad.bind(this)
    this.handleSquareClick = this.handleSquareClick.bind(this)
    this.handleClueClick = this.handleClueClick.bind(this)
    this.calcClueNums = this.calcClueNums.bind(this)

  }
  
  findSquareFromArray (sqs,row,col) {
    return sqs.filter(e => (e.col === col && e.row === row))[0]
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

  loadPuzzle () {
    fetch("http://localhost:8080/puzzle/" + this.state.puzzleInput + "/get")
    .then(res => {
      if(!res.ok) throw(res);
      return(res);
    })
    .then(res => res.json())
    .then(jj => {
      //this.resetSelection(jj.squares)
      jj.squares = this.calcClueNums(jj.squares)
      this.setState({ squares: jj.squares,
                      acrossClues: jj.acrossClues,
                      downClues: jj.downClues})
    })
    .catch(res => {
      this.setState({ gameMessage: "Error Loading Puzzle..." });
      console.log(res)
    });
  }

  componentDidMount() {
    this.loadPuzzle();
    this.selectSquare(0,0);
  }

  handlePuzzleInputChange (event) {
    this.setState({puzzleInput: event.target.value})
  }

  handlePuzzleLoad (event) {
    this.setState({gameMessage: ""})
    this.loadPuzzle();
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
        <div>
          <span className="GameMessage">{this.state.gameMessage}</span>
        </div>
        <div>
          <span>What x do you want?  </span>
          <input type='text' value={this.state.puzzleInput} onChange={this.handlePuzzleInputChange}/>
          <input type='button' value="load puzzle" onClick={this.handlePuzzleLoad}/>
        </div>
       <div className="Grid">
       {squares.map( (t, i) => {
          return (
          <Square value={this.state.squares[i]} key={(t.row*100)+t.col} onClick={this.handleSquareClick}/> 
          );
       }
       )}
        <Selector value={this.state.selectorPos}></Selector>
        </div>

        <AcrossClueList className="AcrossClueList" value={aval} style={astyle} onClick={this.handleClueClick} />
        <DownClueList className="DownClueList" value={dval} style={dstyle} onClick={this.handleClueClick}/>
       </div>

      );}
}


function App() {
  return (
      <Game />
  );
}

export default App;
