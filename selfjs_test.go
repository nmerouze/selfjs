package selfjs_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/nmerouze/selfjs"
)

func eq(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

var script = `
selfjs.handleRequest = function(req, res) {
	if (req.path === '/') {
		res.write('Hello World!');
	}
};
`

func TestHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	selfjs.New(1, script).ServeHTTP(w, r)
	eq(t, "Hello World!\n", w.Body.String())
}
