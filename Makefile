.PHONY: all run build clean

all: run

run:
	air

build:
	go build -o bin/app cmd/server/main.go

clean:
	rm -rf bin/
	rm -rf tmp/