import React from 'react';
import './App.scss';

var sqsize=40;
var curborder=2;

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
    puzzle: { squares: []},
    gameMessage: "",
    puzzleInput: "Apr0914",
    selectorPos: {"row":0, "col":0},
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
    return this.state.puzzle.squares.filter(e => (e.row === x && e.col === y))
  }

  selectSquare (row,col) {
    this.setState({selectorPos: {"row":row, "col":col}});
    let sqs = this.state.puzzle.squares;

    sqs.forEach(e => {
      if(e.row === row && e.col === col)
        e.selected = true;
    })

    this.setState({puzzle: {squares: sqs}})


  }

  handleSquareClick (event, row, col) {
    let sqs = this.state.puzzle.squares;
    let s = this.findSquare(row,col);

    if (!s.isBlack) {
      sqs = this.resetSelection(sqs);
      this.setState({puzzle: {squares: sqs}})
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
      this.setState({ puzzle: jj})

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
          <Square value={this.state.puzzle.squares[i]} key={(t.col*100)+t.row} onClick={this.handleSquareClick}/> 
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
