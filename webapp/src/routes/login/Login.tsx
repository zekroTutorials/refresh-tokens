import "./Login.css";

import { Component } from "react";
import RestAPI from "../../services/api";
import { RouteComponentProps, withRouter } from "react-router-dom";

export class LoginRoute extends Component<RouteComponentProps> {
  state = {
    username: "",
    password: "",
  };

  render() {
    return (
      <div className="login-container">
        <div>
          <input
            placeholder="Username"
            value={this.state.username}
            onChange={(e) => this.setState({ username: e.target.value })}
          ></input>
          <input
            type="password"
            placeholder="Password"
            value={this.state.password}
            onChange={(e) => this.setState({ password: e.target.value })}
          ></input>
          <button onClick={() => this.onLogin()}>Login</button>
        </div>
      </div>
    );
  }

  private async onLogin() {
    try {
      await RestAPI.login(this.state.username, this.state.password);
      this.props.history.replace("/");
    } catch {}
  }
}

export default withRouter(LoginRoute);
