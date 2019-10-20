import React from 'react';
import { dirs } from './App';


class Clue extends React.Component {
  constructor() {
    super();
    this.ref = React.createRef();
  }
  componentDidUpdate() {
    if (this.props.value.selected) {
      this.ref.current.scrollIntoView();
    }
  }
  render() {
    const value = this.props.value;
    let text = String(value.number) + ". " + value.text;
    return <div ref={this.ref} className="Clue" selstyle={value.selected ? 'selected' : 'not-selected'} onClick={(e) => this.props.onClick(e, value.number, value.dir)}>{text}</div>;
  }
}

class ClueList extends React.Component {
  dir = dirs.ACROSS;
  cname = "AcrossList";
  render() {
    const value = this.props.value;
    if (value) {
      return (<div className={this.cname} style={this.props.style}>        {Object.keys(value.values).map((k) => {
        let vv = {
          number: k,
          text: value.values[k],
          dir: this.dir,
          selected: value.selected == k
        };
        return (<Clue value={vv} key={k} onClick={this.props.onClick} />);
      })}
      </div>);
    }
    else {
      return <div />;
    }
  }
}

export class AcrossClueList extends ClueList {
}

export class DownClueList extends ClueList {
  dir = dirs.DOWN;
  cname = "DownList";
}
