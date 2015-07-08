var cluster = require('cluster');
var http = require('http');
var numCPUs = require('os').cpus().length;
var path = require('path');
var fs = require('fs');

var React = require('react');
var Router = require('react-router');
var uni = require('./build/universal');

if (cluster.isMaster) {
  for (var i = 0; i < numCPUs; i++) {
    cluster.fork();
  }

  cluster.on('exit', function(worker, code, signal) {
    console.log('worker ' + worker.process.pid + ' died');
  });
} else {
  http.createServer(function(req, res) {
    var rPath = require('url').parse(req.url).path;

    if (rPath === '/favicon.ico') {
      res.writeHead(204);
      return;
    }

    if (rPath === '/universal.js') {
      res.writeHead(200);
      res.end(fs.readFileSync(path.join(__dirname, './build/bundle.js'), 'utf8'));
      return;
    }

    res.writeHead(200);
    Router.run(uni.routes, rPath, function(Root, state) {
      var html = React.createFactory(uni.Html)({
        markup: React.renderToString(React.createElement(Root, { params: state.params }))
      });
      
      res.end(React.renderToStaticMarkup(html));
    });
  }).listen(8080);
}