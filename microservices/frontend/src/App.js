import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import Users from './pages/Users';
import Products from './pages/Products';
import Orders from './pages/Orders';
import './App.css';

const App = () => {
  return (
    <Router>
      <div className="App">
        <Navbar /> {/* This ensures the Navbar is rendered on all pages */}
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/users" component={Users} />
          <Route path="/products" component={Products} />
          <Route path="/orders" component={Orders} />
        </Switch>
      </div>
    </Router>
  );
};

export default App;
