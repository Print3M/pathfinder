package workers

import "pathfinder/src/store"

type Worker struct {
	input    chan<- store.Url
	done     <-chan struct{}
	isIdle   bool
	Assigned store.Url
}

func NewWorker() (*Worker, chan store.Url, chan struct{}) {
	done := make(chan struct{})
	input := make(chan store.Url)

	return &Worker{
		isIdle: true,
		input:  input,
		done:   done,
	}, input, done
}

func (w *Worker) AssignJob(url store.Url) {
	w.isIdle = false
	w.Assigned = url
	w.input <- url
}

func (w *Worker) isJobDone() bool {
	select {
	case <-w.done:
		return true
	default:
		return false
	}
}

func (w *Worker) Update() {
	if !w.isIdle {
		isDone := w.isJobDone()
		w.isIdle = isDone
	}
}

func (w *Worker) IsIdle() bool {
	return w.isIdle
}

func (w *Worker) Shutdown() {
	close(w.input)
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

func (p *Pool) AreAllWorkersIdle() bool {
	for _, worker := range p.workers {
		if !worker.IsIdle() {
			return false
		}
	}

	return true
}

func (p *Pool) ShutdownAllWorkers() {
	for _, worker := range p.workers {
		worker.Shutdown()
	}
}
