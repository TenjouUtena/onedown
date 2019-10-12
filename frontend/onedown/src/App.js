import React from 'react';
import './App.css';

var sqsize=40;

class Square extends React.Component {
  state = {
  value: {}
  }

  componentDidMount() {
    this.setState({value: this.props.value});
  }


  render() {

    const value = this.state.value;
    var style = {
       top: value.Y*sqsize,
       left: value.X*sqsize,
       width: sqsize-1,
       height: sqsize-1
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
      <div className="Square" style={style} id={black}>
        {clue}
      </div>
    );
  }
}

class Game extends React.Component {
  state = {
    puzzle: { squares: []},
    gameMessage: "",
    puzzleInput: "Oct1219"
  }

  constructor() {
    super();

    this.handlePuzzleInputChange = this.handlePuzzleInputChange.bind(this)
    this.loadPuzzle = this.loadPuzzle.bind(this)
    this.handlePuzzleLoad = this.handlePuzzleLoad.bind(this)

  }

  loadPuzzle () {
    fetch("http://localhost:8080/puzzle/" + this.state.puzzleInput + "/get")
    .then(res => {
      if(!res.ok) throw(res);
      return(res);
    })
    .then(res => res.json())
    .then(jj => {
      this.setState({ puzzle: { squares: jj}})

    })
    .catch(res => {
      this.setState({ gameMessage: "Error Loading Puzzle..." });
      console.log(res)
    });
  }

  componentDidMount() {
    this.loadPuzzle();
  }

  handlePuzzleInputChange (event) {
    this.setState({puzzleInput: event.target.value})
  }

  handlePuzzleLoad (event) {
    this.loadPuzzle();
    this.forceUpdate();
  }

  render () {
      const {puzzle: {squares}} = this.state;
      return (
      <div className="Game">
        <div>
          <span className="GameMessage" id="GameMessage" ref="GameMessage">{this.state.gameMessage}</span>
        </div>
        <div>
          <span>What x do you want?  </span>
          <input type='text' value={this.state.puzzleInput} onChange={this.handlePuzzleInputChange}/>
          <input type='button' value="load puzzle" onClick={this.handlePuzzleLoad}/>
        </div>
       <div className="Grid">
       {squares.map( (t) =>
         <Square value={t} key={(t.Y*100)+t.X}/>
       )}
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
