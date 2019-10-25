import React from 'react';
import { sqsize } from './App';
export class Square extends React.Component {
  render() {
    const value = this.props.value;
    var style = {
      top: value.row * sqsize,
      left: value.col * sqsize
    };
    let clue;
    if (value.clueNum) {
      clue = <span className="SquareClue">{value.clueNum}</span>;
    }
    let black;
    if (value.isBlack) {
      black = "black";
    }
    else {
      black = "white";
    }

    let answer
    if (value.answer) {
        answer = <span className="answer">{value.answer.guess}</span>;
    }

    return (<div className="Square" style={style} id={black} onClick={(e) => this.props.onClick(e, this.props.value.row, this.props.value.col)} onKeyPress={(e) => this.props.onKeyPress(e)} selstyle={this.props.value.selected ? 'selected' : 'notSelected'} tabIndex="0">
      {clue}
      {answer}
    </div>);
  }
}
