# Mac
ARCH = amd64
OS = linux
# Raspberry Pi
#ARCH = arm64
#OS = linux

.PHONY: install build
build: scraper

install: scraper systemd.service
	sed "s/PATH/$(PWD)/g; s/USER/$(USER)/g" systemd.service > /etc/systemd/system/scraper.service
	systemd daemon-reload

scraper: $(shell find . -name ' *.go') $(wildcard templates/* static/*)
	GOARCH=$(ARCH) GOOS=$(OS) go build -tags release .
