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
    puzzle: { squares: []}
  }

  componentDidMount() {
     fetch("http://localhost:8080/puzzle/blah/get")
      .then(res => res.json())
      .then(
       (result)=>{
       this.setState({ puzzle: { squares: result}})})
  }

  render () {
      const {puzzle: {squares}} = this.state;
      return (
      <div className="Game">
      {squares.map( (t) =>
        <Square value={t} key={(t.Y*100)+t.X}/>
      )}
      </div>
      );}
}


function App() {
  return (
      <Game />
  );
}

export default App;
