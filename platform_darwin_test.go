//go:build cgo

package main

import "testing"

func TestNativeStringsDeferredFree(t *testing.T) {
	var strings nativeStrings
	func() {
		defer strings.free()
		strings.add("first")
		strings.add("second")
		if got := len(strings); got != 2 {
			t.Fatalf("native strings before cleanup: %d; want 2", got)
		}
	}()

	if got := len(strings); got != 0 {
		t.Fatalf("native strings remaining after deferred cleanup: %d; want 0", got)
	}

	strings.free()
	strings.add("reused")
	strings.free()
	if got := len(strings); got != 0 {
		t.Fatalf("native strings remaining after reuse: %d; want 0", got)
	}
	strings.free()
}
