package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc -Wall
#cgo LDFLAGS: -framework Cocoa -framework IOKit
#include <stdlib.h>
#include <IOKit/pwr_mgt/IOPMLib.h>
#include "app.h"
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

func preventSleep() (func(), error) {
	name := C.CString("HangOn is keeping this Mac awake")
	defer C.free(unsafe.Pointer(name))
	description := C.CFStringCreateWithCString(0, name, C.kCFStringEncodingUTF8)
	if description == 0 {
		return nil, fmt.Errorf("could not allocate sleep-prevention description")
	}
	defer C.CFRelease(C.CFTypeRef(description))
	var assertion C.IOPMAssertionID
	result := C.IOPMAssertionCreateWithName(
		C.kIOPMAssertionTypePreventUserIdleSystemSleep,
		C.kIOPMAssertionLevelOn, description, &assertion)
	if result != C.kIOReturnSuccess {
		return nil, fmt.Errorf("could not prevent idle sleep (IOKit error 0x%x)", uint32(result))
	}
	return func() {
		if result := C.IOPMAssertionRelease(assertion); result != C.kIOReturnSuccess {
			fmt.Fprintf(os.Stderr, "HangOn: release sleep assertion: IOKit error 0x%x\n", uint32(result))
		}
	}, nil
}

type nativeStrings []*C.char

func (s *nativeStrings) add(text string) *C.char {
	value := C.CString(text)
	*s = append(*s, value)
	return value
}

// A pointer receiver lets deferred cleanup see strings added later.
func (s *nativeStrings) free() {
	for _, value := range *s {
		C.free(unsafe.Pointer(value))
	}
	*s = nil
}

func runNativeApp(icon []byte) {
	var strings nativeStrings
	defer strings.free()
	config := C.AppConfig{
		title:         strings.add("HangOn"),
		status:        strings.add("Keeping your Mac awake"),
		stop:          strings.add("Stop HangOn"),
		tooltip:       strings.add("HangOn - keeping your Mac awake"),
		accessibility: strings.add("HangOn, keeping your Mac awake"),
	}
	data := C.CBytes(icon)
	defer C.free(data)
	C.runApp(config, data, C.size_t(len(icon)))
}

func showWakePrompt() bool {
	var strings nativeStrings
	defer strings.free()
	return C.showPrompt(
		strings.add("HangOn is still keeping your Mac awake"),
		strings.add("This app is running in the background and preventing idle sleep. Would you like to stop it?"),
		strings.add("Keep Running"),
		strings.add("Stop HangOn")) == 1
}

func showError(title, detail string) {
	var strings nativeStrings
	defer strings.free()
	C.showPrompt(strings.add(title), strings.add(detail), strings.add("OK"), nil)
}

func requestStop() { C.requestStop() }

//export goDidWake
func goDidWake() { activeApp.wake() }

//export goStop
func goStop() { activeApp.stop() }
