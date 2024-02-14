# Scrape Websites for numbers and turn them into diagrams

Main usecase: To scrape WLAN-connected electical sockets and aggregate/show their usage.

Written purely for personal use, so don't expect high quality code or reliability.

## Cross-compiling
Cool in theory, could not get it to work in practice yet.
Setting the right `ARCH` helps, but I could not get CGO to work yet.
Just install golang>=1.16 instead and run `make build`.

## Local testing
A simple `go run .` should start the server in a development mode (aka reloading http files).
Use `bash random_data.sh` to load some random data into the generated sqlite db, if there are no actual sensors configured.
