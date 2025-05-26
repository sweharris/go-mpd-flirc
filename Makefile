SRC=$(wildcard *.go)

mpd-flirc: $(SRC)
	go build
