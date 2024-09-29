package workers

import "scraper/store"

type Worker struct {
	input  chan<- store.Url
	done   <-chan struct{}
	isIdle bool
}

func NewWorker() (*Worker, chan store.Url, chan struct{}) {
	done := make(chan struct{})
	input := make(chan store.Url)

	return &Worker{
		isIdle: true,
		done:   done,
		input:  input,
	}, input, done
}

func (w *Worker) AssignJob(url store.Url) {
	w.isIdle = false
	w.input <- url
}

func (w *Worker) SetIdle() {
	w.isIdle = true
}

func (w *Worker) IsIdle() bool {
	return w.isIdle
}

func (w *Worker) Done() <-chan struct{} {
	return w.done
}

type Pool struct {
	workers     []Worker
	Size        uint64
	idleCounter uint64
}

func NewPool(size uint64) Pool {
	return Pool{
		workers:     make([]Worker, size),
		Size:        size,
		idleCounter: size,
	}
}

func (p *Pool) InitWorkers(job func(chan store.Url, chan struct{})) {
	for i := uint64(0); i < p.Size; i++ {
		worker, input, done := NewWorker()
		p.workers[i] = *worker

		go job(input, done)
	}
}

func (p *Pool) GetWorkerById(workerId uint64) *Worker {
	return &p.workers[workerId]
}

func (p *Pool) AllWorkersIdle() bool {
	for _, worker := range p.workers {
		if !worker.isIdle {
			return false
		}
	}

	return true
}
