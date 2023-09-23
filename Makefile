# Mac
ARCH = amd64
OS = darwin
# Raspberry Pi
#ARCH = arm64
#OS = linux

.PHONY: install build clean
build: sensor-scraper

install: build systemd.service
	sed "s/PATH/$(PWD)/g; s/USER/$(USER)/g" systemd.service > /etc/systemd/system/sensor-scraper.service
	systemd daemon-reload

scraper: $(shell find . -name ' *.go') $(wildcard templates/* static/*)
	GOARCH=$(ARCH) GOOS=$(OS) go build -tags release .

clean:
	rm sensor-scraper
