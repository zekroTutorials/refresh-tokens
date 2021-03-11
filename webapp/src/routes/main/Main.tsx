import "./Main.css";

import { Component } from "react";
import { RouteComponentProps, withRouter } from "react-router";
import RestAPI from "../../services/api";

export class MainRoute extends Component<RouteComponentProps> {
  state = {
    id: "",
  };

  async componentDidMount() {
    try {
      const id = await RestAPI.me();
      this.setState({ id });
    } catch (e) {}
  }

  render() {
    return <div>ID: {this.state.id}</div>;
  }
}

export default withRouter(MainRoute);
