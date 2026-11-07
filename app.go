package main

import "sync"

type app struct {
	mu        sync.Mutex
	prompting bool
	stopped   bool
	prompt    func() bool
	quit      func()
}

func (a *app) wake() {
	a.mu.Lock()
	if a.stopped || a.prompting {
		a.mu.Unlock()
		return
	}
	a.prompting = true
	a.mu.Unlock()

	defer func() {
		a.mu.Lock()
		a.prompting = false
		a.mu.Unlock()
	}()
	if a.prompt() {
		a.stop()
	}
}

func (a *app) stop() {
	a.mu.Lock()
	if a.stopped {
		a.mu.Unlock()
		return
	}
	a.stopped = true
	a.mu.Unlock()
	a.quit()
}
