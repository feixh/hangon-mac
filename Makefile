.PHONY: build run

build:
	mkdir -p build/HangOn.app/Contents/MacOS
	CGO_ENABLED=1 go build -o build/HangOn.app/Contents/MacOS/hangon .
	cp Info.plist build/HangOn.app/Contents/Info.plist

run: build
	open build/HangOn.app
