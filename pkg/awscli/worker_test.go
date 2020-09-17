package awscli

import (
	"testing"
)

func TestStart_Error(t *testing.T) {
	e, s, w := make(chan int), make(chan int), worker{
		consumer: fakeErrQueue{},
	}

	go Start(w, e, s)

	if got := <-e; got != 1 {
		t.Errorf("Got: %d; want: 1", got)
	}
}

func TestStart(t *testing.T) {
	e, s, w := make(chan int), make(chan int), worker{
		consumer: fakeOkQueue{},
		poller:   fakePoller{},
	}

	go Start(w, e, s)

	if got := <-s; got != 1 {
		t.Errorf("Got: %d; want: 1", got)
	}
}
