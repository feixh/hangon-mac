package main

import "testing"

func TestWakeKeepRunning(t *testing.T) {
	prompts, quits := 0, 0
	a := &app{
		prompt: func() bool { prompts++; return false },
		quit:   func() { quits++ },
	}
	a.wake()
	a.wake()
	if prompts != 2 || quits != 0 {
		t.Fatalf("prompts=%d quits=%d; want 2 and 0", prompts, quits)
	}
}

func TestWakeStop(t *testing.T) {
	prompts, quits := 0, 0
	a := &app{
		prompt: func() bool { prompts++; return true },
		quit:   func() { quits++ },
	}
	a.wake()
	a.wake()
	a.stop()
	if prompts != 1 || quits != 1 {
		t.Fatalf("prompts=%d quits=%d; want 1 and 1", prompts, quits)
	}
}

func TestWakeWhilePromptOpen(t *testing.T) {
	prompts := 0
	var a *app
	a = &app{
		prompt: func() bool {
			prompts++
			a.wake()
			return false
		},
		quit: func() { t.Fatal("unexpected quit") },
	}
	a.wake()
	if prompts != 1 {
		t.Fatalf("prompts=%d; want 1", prompts)
	}
}

func TestStopWhilePromptOpen(t *testing.T) {
	opened, dismiss := make(chan struct{}), make(chan struct{})
	finished := make(chan struct{})
	quits := 0
	a := &app{
		prompt: func() bool {
			close(opened)
			<-dismiss
			return true
		},
		quit: func() { quits++ },
	}
	go func() {
		a.wake()
		close(finished)
	}()
	<-opened
	a.stop()
	close(dismiss)
	<-finished
	if quits != 1 {
		t.Fatalf("quits=%d; want 1", quits)
	}
}
