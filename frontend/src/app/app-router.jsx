import React from "react";
import { BrowserRouter as Router, Route, Switch, Redirect } from "react-router-dom";

import Home from './Home/Home'
import SearchNotFound from './SearchNotFound/SearchNotFound'
import BlockDetails from './BlockDetails/BlockDetails';
import TransactionDetails from './TransactionDetails/TransactionDetails';
import Address from "./Address/Address";
import Blocks from "./Blocks/Blocks";
import Stakes from "./Stakes/Stakes";
import Stats from "./Stats/Stats";
import Transactions from "./Transactions/Transactions";

export default function () {
  return (
    <Switch>
      <Route path="/" exact component={Home} />
      <Route path="/search-not-found" exact component={SearchNotFound} />
      <Route path="/blocks" exact component={Blocks} />
      <Route path="/stakes" exact component={Stakes} />
      <Route path="/stats" exact component={Stats} />
      <Route path="/transactions" exact component={Transactions} />
      <Route path="/blocks/:height" exact component={BlockDetails} />
      <Route path="/transactions/:hashId" exact component={TransactionDetails} />
      <Route path="/address/:address" exact component={Address} />
      <Redirect to="/" />
      {/* <Route path="/:network/:type?" exact component={Blockchains} /> */}
    </Switch>
  );
};

