package parallelizer

import (
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	calls := make(chan struct{}, 10)
	done := make(chan struct{}, 1)
	prlzr := New(1)

	go prlzr.Execute(func() {
		<-done
	})

	go prlzr.Execute(func() {
		calls <- struct{}{}
	})

	if len(calls) != 0 {
		t.Error("Next executions shouldn't be executed until first one ends")
	}

	done <- struct{}{}

	time.Sleep(time.Millisecond)

	if len(calls) != 1 {
		t.Error("Second execution should have been executed")
	}
}

func TestWait(t *testing.T) {
	concurrency := 5
	executions := 4
	calls := make(chan struct{}, executions)
	prlzr := New(concurrency)

	for x := 0; x < executions; x++ {
		prlzr.Execute(func() {
			time.Sleep(time.Millisecond)
			calls <- struct{}{}
		})
	}

	if len(calls) != 0 {
		t.Error("Executions shouldn't be executed yet due a stablished delay")
	}

	prlzr.Wait()

	if len(calls) != executions {
		t.Error("Wait should be wait until all executions are finished")
	}
}
