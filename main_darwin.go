package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var activeApp *app

func main() {
	runtime.LockOSThread()
	os.Exit(run())
}

func run() int {
	icon, err := coffeeIcon()
	if err != nil {
		fmt.Fprintf(os.Stderr, "HangOn: create menu icon: %v\n", err)
		showError("HangOn could not start", err.Error())
		return 1
	}
	release, err := preventSleep()
	if err != nil {
		fmt.Fprintf(os.Stderr, "HangOn: %v\n", err)
		showError("HangOn could not keep your Mac awake", err.Error())
		return 1
	}
	defer release()
	activeApp = &app{prompt: showWakePrompt, quit: requestStop}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)

	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-signals:
			activeApp.stop()
		case <-done:
		}
	}()

	runNativeApp(icon)
	return 0
}
