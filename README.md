# Scrape Websites for numbers and turn them into diagrams

Main usecase: To scrape WLAN-connected electical sockets and aggregate/show their usage.

Written purely for personal use, so don't expect high quality code or reliability.

## Configuration
Example configuration for a [Shelly Plug S](https://shelly-api-docs.shelly.cloud/gen1/#shelly-plug-plugs-status), just replace the credentials and IP.

```toml
[[sensor]]
 id=1
name="Current power being drawn (W)"
url="http://user:pass@127.0.0.1:80/status"
jsonPath="meters.0.power"

[[sensor]]
id=2
name="Cummulative power since start (Wm)"
url="http://user:pass@127.0.0.1:80/status"
jsonPath="meters.0.total"
```

## Cross-compiling
I could not get the CGO sqlite library to cross-compile (compiling on the target worked fine).
But after switching to `go-sqlite`, building for ARM on Mac is no problem anymore.
See the [Makefile](./Makefile) for details.

## Local testing
A simple `go run .` should start the server in a development mode (aka reloading http files).
Use `bash random_data.sh` to load some random data into the generated sqlite db, if there are no actual sensors configured.
