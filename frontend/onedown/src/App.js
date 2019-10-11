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
       top: value.x*sqsize,
       left: value.y*sqsize,
       width: sqsize-1,
       height: sqsize-1
    };
    return (
      <div className="Square" style={style}>

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
      <div>
      {squares.map( (t) =>
        <Square value={t} key={(t.y*100)+t.x}/>
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
