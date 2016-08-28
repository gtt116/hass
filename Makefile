all: win linux

win:
	GOOS=windows GOARCH=amd64 go build

linux:
	go build
