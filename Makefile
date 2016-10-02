all: win linux

win:
	GOOS=windows go build

linux:
	go build
