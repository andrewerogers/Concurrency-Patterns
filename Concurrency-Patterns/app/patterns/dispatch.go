package patterns

import "fmt"

const (
	MAX_QUEUE = 10
)

type Job struct {
	Payload string
	Handler func() error
}

var JobQueue chan Job

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
}

// Dispatcher will start MAX_WORKERS number of job channels
// with registered workers to execute jobs
func NewDispatcher() *Dispatcher {
	pool := make(chan chan Job, MAX_WORKERS)
	return &Dispatcher{WorkerPool: pool}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < MAX_WORKERS; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

// dispatch will listen to the job queue and
// pass the job to the next available worker
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool  chan chan Job
	JobChannel  chan Job
	quit    	chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				err := job.Handler()
				if err != nil {
					fmt.Println(err)
				}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}