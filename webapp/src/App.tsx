import React, { Component } from "react";
import "./App.css";
import { BrowserRouter as Router, Redirect, Route } from "react-router-dom";
import LoginRoute from "./routes/login/Login";
import RestAPI from "./services/api";
import { MainRoute } from "./routes/main/Main";

export default class App extends Component {
  state = {
    redirect: (null as any) as string,
    id: "",
  };

  async componentDidMount() {
    RestAPI.events.on("authentication-error", () => {
      // this.redirect("/login");
    });

    try {
      await RestAPI.me();
    } catch (e) {}
  }

  render() {
    return (
      <div>
        <Router>
          <Route exact path="/login" component={LoginRoute}></Route>
          <Route exact path="/" component={MainRoute}></Route>

          {this.state.redirect && <Redirect to={this.state.redirect} />}
        </Router>
      </div>
    );
  }

  private redirect(to: string) {
    this.setState({ redirect: to });
    setTimeout(() => this.setState({ redirect: null }), 10);
  }
}
