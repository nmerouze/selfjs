var React = require('react');
var Router = require('react-router');

var Route = Router.Route;
var DefaultRoute = Router.DefaultRoute;
var RouteHandler = Router.RouteHandler;
var Link = Router.Link;

var Index = React.createClass({
  render: function() {
    return (
      <div>Hello World! Go to <Link to="/about">About page</Link>.</div>
    );
  }
});

var About = React.createClass({
  render: function() {
    return (
      <div>About page! <Link to="/">Return Home</Link>.</div>
    );
  }
});

var Layout = React.createClass({
  render () {
    return (
      <div>
        <h1>App</h1>
        <RouteHandler/>
      </div>
    )
  }
});

var Html = React.createClass({
  render: function() {
    return (
      <html>
        <head>
          <title>Universal App</title>
          <meta charSet="utf-8" />
        </head>
        <body dangerouslySetInnerHTML={{__html: this.props.markup}} />
        <script src="/universal.js" async></script>
      </html>
    );
  }
});

var routes = (
  <Route handler={Layout}>
    <DefaultRoute handler={Index}/>
    <Route path="about" handler={About}/>
  </Route>
);

module.exports = {
  Html: Html,
  routes: routes
};

if (typeof selfjs !== 'undefined') {
  selfjs.handleRequest = function(req, res) {
    Router.run(routes, req.path, function(Root, state) {
      var html = React.createFactory(Html)({
        markup: React.renderToString(<Root params={state.params}/>)
      });
      
      res.write(React.renderToStaticMarkup(html));
    });
  }
}

if (typeof window !== 'undefined') {
  function render() {
    Router.run(routes, Router.HistoryLocation, function(Root, state) {
      var body = React.createFactory(Root)({
        params: state.params
      });

      React.render(body, document.body);
    });
  }

  window.onload = render;
}