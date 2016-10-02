all: win linux mac

win:
	GOOS=windows go build

linux:
	go build -o hass

mac:
	GOOS=darwin go build -o hass_mac
