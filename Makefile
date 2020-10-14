.PHONY: build,

build:
	go build -o mrkt

run:
	./mrkt

dev: dev-build dev-run

dev-build:
	go build -race -o mrkt

dev-run:
	./mrkt -cpuprofile=cpu.prof -memprofile=mem.prof

cpu-profile:
	go tool pprof mrkt cpu.prof

mem-profile:
	go tool pprof mrkt mem.prof
