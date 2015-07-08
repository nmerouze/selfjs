# Self.js [![GoDoc](https://godoc.org/github.com/nmerouze/selfjs?status.png)](https://godoc.org/github.com/nmerouze/selfjs)

Wrapper around [v8worker](https://github.com/ry/v8worker) and `net/http`. Run Universal React applications faster than Node.js.

# Setup

First you need to install [v8worker](https://github.com/ry/v8worker). Then:

```
go get -u github.com/nmerouze/selfjs
```

# Usage

``` go
import (
  "io/ioutil"
  "net/http"
  "runtime"

  "github.com/nmerouze/selfjs"
)

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  js, _ := ioutil.ReadFile("path/to/file.js")
  http.Handle("/", selfjs.New(runtime.NumCPU(), `
    selfjs.handleRequest = function(req, res) {
      res.write('Hello World!');
    }
  `))
  http.ListenAndServe(":8080", nil)
}
```

# Example

```
cd example
npm install
npm run build
go run server.go
open http://localhost:8080
```

The application is rendered on the server. Both the client and the server share the same file (universal.js) with a few conditions to select the right rendering function. The code also runs on Node.js, just run `node server.js`. Self.js is 50% than Node.js while Node.js consumes 50% memory.