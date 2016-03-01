package selfjs

import "github.com/ry/v8worker"

type workerPool struct {
	ch chan *Worker
}

func (o *workerPool) get() *Worker {
	return <-o.ch
}

func (o *workerPool) put(ot *Worker) {
	o.ch <- ot
}

func newPool(size int, fn func(*Worker)) *workerPool {
	pool := &workerPool{
		ch: make(chan *Worker, size),
	}
loop:
	for {
		select {
		case pool.ch <- newWorker(pool, fn):
		default:
			break loop
		}

	}
	return pool
}

func discardSendSync(msg string) string { return "" }

func newWorker(pool *workerPool, fn func(*Worker)) *Worker {
	w := new(Worker)

	v8w := v8worker.New(func(msg string) {
		w.ch <- msg
	}, discardSendSync)

	w.Worker = v8w

	fn(w)

	return w
}
