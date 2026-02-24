GOX ?= go

.PHONY: build

build: examples/bin/demo

examples/bin/demo: *.go examples/demo.go
	$(GOX) build -o examples/bin/demo examples/demo.go

clean:
	rm examples/bin/*
