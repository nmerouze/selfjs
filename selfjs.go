package selfjs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"runtime"

	"github.com/ry/v8worker"
)

type Worker struct {
	*v8worker.Worker
	ch chan string
}

type message struct {
	Fn   string        `json:"fn"`
	Args []interface{} `json:"args"`
}

func readFile(filePath string) string {
	_, filename, _, _ := runtime.Caller(1)
	file, _ := ioutil.ReadFile(path.Join(path.Dir(filename), filePath))
	return string(file)
}

func New(poolSize int, script string) http.Handler {
	bundle := bytes.NewBufferString(selfjs)
	bundle.WriteString(script)

	pool := newPool(poolSize, func(w *Worker) {
		if err := w.Load("bundle.js", bundle.String()); err != nil {
			log.Panicf("error while loading js: %#v", err)
		}
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := make(chan string)
		worker := pool.get()
		worker.ch = c
		req := map[string]interface{}{"path": r.URL.Path}
		msg := message{Fn: "beforeHandleRequest", Args: []interface{}{req}}
		sMsg, _ := json.Marshal(msg)
		go worker.Send(string(sMsg))
		res := <-c
		worker.ch = nil
		pool.put(worker)
		w.Write([]byte(res + "\n"))
	})
}

const selfjs = `
var selfjs = {};

$recv(function(msg) {
  var pMsg = JSON.parse(msg);
  this[pMsg.fn].apply(null, pMsg.args);
});

function beforeHandleRequest(req) {
  var res = {
    write: function(str) {
      $send(str);
    }
  };

  selfjs.handleRequest(req, res);
}
`
