all: win linux mac
target = ./cmd/hass/*.go

win:
	GOOS=windows go build -o hass.exe $(target)

linux:
	go build -o hass $(target)

mac:
	GOOS=darwin go build -o hass_mac $(target)

test:
	go test -v ./pkg/...
