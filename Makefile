# Mac
ARCH = amd64
OS = darwin
# Raspberry Pi
#ARCH = arm
#OS = linux

.PHONY: install build clean
build: sensor-scraper

install: build systemd.service
	sed "s_PWD_$(PWD)_g; s/USER/$(USER)/g" systemd.service | sudo tee /etc/systemd/system/sensor-scraper.service
	sudo systemctl daemon-reload

sensor-scraper: $(shell find . -name ' *.go') $(wildcard templates/* static/*)
	GOARCH=$(ARCH) GOOS=$(OS) go build -tags release .
sensor-scraper-pi: $(shell find . -name ' *.go') $(wildcard templates/* static/*)
	GOARCH=arm GOOS=linux go build -tags release .

clean:
	rm sensor-scraper
