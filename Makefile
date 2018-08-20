all: win linux mac
target = ./cmd/hass/*.go

win:
	GOOS=windows GOARCH=amd64 go build -o hass.exe $(target)

linux:
	GOOS=linux GOARCH=amd64 go build -o hass $(target)

mac:
	GOOS=darwin GOARCH=amd64 go build -o hass_mac $(target)

test:
	go test -v ./pkg/...
