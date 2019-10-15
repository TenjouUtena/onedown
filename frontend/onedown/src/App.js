import React from 'react';
import './App.scss';

var sqsize=40;
var curborder=2;

class Selector extends React.Component {

  render() {
    const value = this.props.value;

    var style = {
        top: value.Y*sqsize + curborder,
        left: value.X*sqsize + curborder
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
       top: value.Y*sqsize,
       left: value.X*sqsize
    };
    let clue;
    if (value.DrawDown) {
      clue = <span className="SquareClue" >{value.DownClue}</span>
    }
    if (value.DrawAcross) {
      clue = <span className="SquareClue" >{value.AcrossClue}</span>
    }

    let black;
    if (value.Black) {
      black = "black"
    } else {
      black = "white"
    }
    return (
      <div className="Square" style={style} id={black} onClick={(e) => this.props.onClick(e, this.props.value.X, this.props.value.Y)}
                              selstyle={this.props.value.selected ? 'selected' : 'notSelected'}>
        {clue}
      </div>
    );
  }
}

class Game extends React.Component {
  state = {
    puzzle: { squares: []},
    gameMessage: "",
    puzzleInput: "Oct1219",
    selectorPos: {"X":0, "Y":0},
    selectingAcross: true
  }

  constructor() {
    super();

    this.handlePuzzleInputChange = this.handlePuzzleInputChange.bind(this)
    this.loadPuzzle = this.loadPuzzle.bind(this)
    this.handlePuzzleLoad = this.handlePuzzleLoad.bind(this)
    this.handleSquareClick = this.handleSquareClick.bind(this)

  }

  findSquare (x,y) {
    return this.state.puzzle.squares.reduce((t,c) => {
      if (c.X === x && c.Y === y)
        return c
      return t
    },false)
  }

  selectSquare (x,y) {
    this.setState({selectorPos: {"X":x, "Y":y}});
    let sqs = this.state.puzzle.squares;

    // TODO: pls fix for candidate squares.
    sqs.forEach(e => {
      if(e.X === x && e.Y === y) {
        e.selected = true;
      }
    })

    this.setState({puzzle: {squares: sqs}})

  }

  handleSquareClick (event, x, y) {
    let sqs = this.state.puzzle.squares;
   let s = this.findSquare(x,y);

    if (!s.Black) {
      sqs = this.resetSelection(sqs);
      this.setState({puzzle: {squares: sqs}})
      this.selectSquare(x,y);
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
      jj = this.resetSelection(jj)
      this.setState({ puzzle: { squares: jj}})

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
      const {puzzle: {squares}} = this.state;
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
          <Square value={this.state.puzzle.squares[i]} key={(t.Y*100)+t.X} onClick={this.handleSquareClick}/> 
          );
       }
       )}
       <Selector value={this.state.selectorPos}></Selector>
       </div>
      </div>

      );}
}


function App() {
  return (
      <Game />
  );
}

export default App;
