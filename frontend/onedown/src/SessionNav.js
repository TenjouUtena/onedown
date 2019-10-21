import React from 'react';



export class SessionNav extends React.Component {
  state = {
    sessions: [],
    uploadFile: {}
  };

  constructor(props) {
      super(props)

      this.selector = React.createRef();
  }

  handleChange(event) {
    this.setState({ uploadFile: event.target.files[0] });
  }

  handleCreateSession(event) {
    var data = new FormData()
    data.append('puzFile',this.state.uploadFile)
    fetch("http://localhost:8080/session" , {method:'POST', body: data})
    .then(res => {
        if(!res.ok) throw(res);
        return(res);
      })
      .then(e => {
          this.loadSessions();
      })
      .catch(res => {
        console.log(res)
      });

  }

  loadSessions() {
      fetch("http://localhost:8080/session")
      .then(res => {
        if(!res.ok) throw(res);
        return(res);
      })
      .then(res => res.json())
      .then(jj => {
          this.setState({sessions: jj})
      })
      .catch(res => {
        console.log(res)
      });
  }

  connectSession(event) {
      this.props.buildws("ws://localhost:8080/session/" + this.selector.current.value)

      document.getElementsByClassName('SessionNav')[0].style.height=0;
      document.getElementsByClassName('SessionNav')[0].style.borderStyle='none';
  }

  componentDidMount() {
      this.loadSessions()
  }

  render() {
    return (<div className="SessionNav">
      <h1>Session Maintenance</h1>
      <p />
      <input type='file' name='file' onChange={e => { this.handleChange(e); }} />
      <input type='button' value="Create Session from Puzzle" onClick={(e) => this.handleCreateSession(e)} />
      <p />
      <select ref={this.selector}>
          {this.state.sessions.map((s) => {
              return <option value={s}>{s}</option>
          })}
      </select>
      <input type='button' value="Join Session" onClick={e => this.connectSession(e)} />

    </div>);
  }
}
