# IPERF3EXPORTER

[![Build Status](https://ci.xsfx.dev/api/badges/xsteadfastx/iperf3exporter/status.svg?ref=refs/heads/main)](https://ci.xsfx.dev/xsteadfastx/iperf3exporter)
[![Go Reference](https://pkg.go.dev/badge/go.xsfx.dev/iperf3exporter.svg)](https://pkg.go.dev/go.xsfx.dev/iperf3exporter)

A iperf3 speedtest exporter for prometheus

![readme](./README.gif)

It runs the `iperf3` command as client. Once as server sends/client receives and once as client sends/server receives. It parses the JSON output and exports them as prometheus metrics.

## Installation

### via docker

```shell
docker run -d --name iperf3exporter ghcr.io/xsteadfastx/iperf3exporter:0.1.1
```

### via package

You can get `apt`, `rpm` and `apk` packages on the [release page](https://github.com/xsteadfastx/iperf3exporter/releases). They also include an init file.

### via archive

For easy testing you can download the `tar.gz`-archive from the [release page](https://github.com/xsteadfastx/iperf3exporter/releases), extract it and run it.

## Usage

```shell
Usage:
  iperf3exporter [flags]

Flags:
  -c, --config string      config file
  -h, --help               help for iperf3exporter
      --listen string      listen string (default "127.0.0.1:9119")
      --log-colors         colorful log output (default true)
      --log-json           JSON log output
      --process-metrics    exporter process metrics (default true)
      --time int           time in seconds to transmit for (default 5)
      --timeout duration   scraping timeout (default 1m0s)
  -v, --version            print version
```

### Configuration

#### File

All flags can also be set through a config file. Here is an example:

```toml
[exporter] # everything related to the exporter itself
listen = "0.0.0.0:9119" # connection string for the webserver
timeout = "1m" # timeout of the iperf3 command to run
process_metrics = true # export go process metrics

[log]
json = true # enables json log output
colors = false # disable colors. this is only usable if log.json is set to false

[iperf3] # straight up iperf3 command line flag options
time = 10 # this sets the --time flag of iperf3 to 10
```

#### Environment variables

Its also possible to set this settings through environment variables. The environment prefix is `IPERF3EXPORTER`.

```shell
# this will disable colorful logs
IPERF3EXPORTER_LOG_COLORS=false /usr/local/bin/iperf3exporter

# this sets the iperf3 time flag to 10 seconds
IPERF3EXPORTER_IPERF3_TIME=10 /usr/local/bin/iperf3exporter
```

## Example prometheus config

You can find a scrape config example [here](./test/prometheus.yml). This is the config that gets spun up while testing things for me locally. It replaces the targets with the real exporter adress and adds a label `host` that can be used to identify the scrape boxes and not just the iperf3 servers to test against.

## Exposed metrics

| name                                     | type  |
| ---------------------------------------- | ----- |
| iperf3_download_sent_bits_per_second     | gauge |
| iperf3_download_sent_seconds             | gauge |
| iperf3_download_sent_bytes               | gauge |
| iperf3_download_received_bits_per_second | gauge |
| iperf3_download_received_seconds         | gauge |
| iperf3_download_received_bytes           | gauge |
| iperf3_upload_sent_bits_per_second       | gauge |
| iperf3_upload_sent_seconds               | gauge |
| iperf3_upload_sent_bytes                 | gauge |
| iperf3_upload_received_bits_per_second   | gauge |
| iperf3_upload_received_seconds           | gauge |
| iperf3_upload_received_bytes             | gauge |
