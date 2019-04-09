package parallelizer

import "sync"

// Parallelizer manages parallel executions with a limit
type Parallelizer struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

// Execute a function that will wait once the concurrency limit is reached
func (p *Parallelizer) Execute(f func()) {
	p.sem <- struct{}{}
	p.wg.Add(1)

	go func() {
		defer func() {
			<-p.sem
			p.wg.Done()
		}()

		f()
	}()
}

// Wait for all the task to be executed
func (p *Parallelizer) Wait() {
	p.wg.Wait()
}

// New declares a parallelizer handler
func New(threads int) *Parallelizer {
	return &Parallelizer{
		sem: make(chan struct{}, threads),
	}
}
