# HangOn

A macOS menu-bar app built with Go and a thin native AppKit bridge. A coffee
cup with steam appears while HangOn prevents idle system sleep.

## Build and run

Requires macOS, Go 1.24 or later, and Xcode Command Line Tools (`xcode-select --install`).

```sh
make run
```

This builds `build/HangOn.app` and opens it without a Dock icon or terminal
window. You can also copy the app to `/Applications`. For foreground development,
use `go run .`.

Click the coffee cup and choose **Stop HangOn** to quit and restore normal idle
sleep. After the Mac wakes from system sleep, a popup offers **Keep Running** or
**Stop HangOn**.

HangOn uses an IOKit power assertion equivalent to `caffeinate -i`: no child
process, administrator permission, or permanent power-setting changes are
needed. The assertion lasts only for the lifetime of the app. Display sleep,
locking the screen, closing the lid, and deliberate sleep remain available.
The reminder is triggered by system wake, not merely unlocking the screen or
waking the display.

## Implementation

Go owns application state, wake/stop decisions, duplicate-prompt suppression,
UI text, PNG icon rendering, signal handling, and the IOKit assertion lifecycle.
Objective-C only adapts AppKit: it displays the menu and alerts, forwards native
events to Go, and runs the macOS event loop. There are no third-party dependencies.

Run the Go behavior tests with `go test -race ./...`.
