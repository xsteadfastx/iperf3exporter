<h1 align="center">🚄 IPERF3EXPORTER 💨</h1>
<div align="center">

A iperf3 speedtest exporter for prometheus

[![Build Status](https://ci.xsfx.dev/api/badges/xsteadfastx/iperf3exporter/status.svg?ref=refs/heads/main)](https://ci.xsfx.dev/xsteadfastx/iperf3exporter)
[![Go Reference](https://pkg.go.dev/badge/go.xsfx.dev/iperf3exporter.svg)](https://pkg.go.dev/go.xsfx.dev/iperf3exporter)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/xsteadfastx/iperf3exporter.svg)](https://github.com/xsteadfastx/iperf3exporter)
[![Go Report Card](https://goreportcard.com/badge/go.xsfx.dev/iperf3exporter)](https://goreportcard.com/report/go.xsfx.dev/iperf3exporter)

![readme](./README.gif)

</div>

It runs the `iperf3` command as client. Once as server sends/client receives and once as client sends/server receives. It parses the JSON output and exports them as prometheus metrics.

## Installation

### via docker

```shell
docker run -d --name iperf3exporter ghcr.io/xsteadfastx/iperf3exporter:latest
```

**Notice**: Please use a fixed version for productive use!

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

```yaml
scrape_configs:
  - job_name: speedtest-myfunnybox
    scrape_interval: 2m # maybe a even higher interval would be useful. not fill the whole traffic just with speedtests ;-)
    scrape_timeout: 1m # a little higher timeout. because the scrape can take a while
    metrics_path: /probe
    static_configs:
      - targets:
          - speedtest.wobcom.de # default port 5201 is used
          - footest.bar.tld:1234 # target with defined port
    relabel_configs:
      # takes the address from the targets and uses it as url parameter key `target`
      - source_labels: [__address__]
        target_label: __param_target

      # takes that address and stores it in the label `instance`
      - source_labels: [__param_target]
        target_label: instance

      # replaces the scrape address with the real hostname:port of the exporter.
      # so it can use the targets for defining the iperf3 servers to use.
      - target_label: __address__
        replacement: 192.168.39.191:9119
```

In this example it replaces the targets with the real exporter adress and adds a label `host` that can be used to identify the scrape boxes and not just the iperf3 servers to test against.

You can specify a port for the iperf3 server target. If its not set, it will use the default port `5201`.

## Exposed metrics

| name                                     | type  |
| ---------------------------------------- | ----- |
| iperf3_download_sent_bits_per_second     | gauge |
| iperf3_download_sent_seconds             | gauge |
| iperf3_download_sent_bytes               | gauge |
| iperf3_download_sent_retransmits         | gauge |
| iperf3_download_received_bits_per_second | gauge |
| iperf3_download_received_seconds         | gauge |
| iperf3_download_received_bytes           | gauge |
| iperf3_upload_sent_bits_per_second       | gauge |
| iperf3_upload_sent_seconds               | gauge |
| iperf3_upload_sent_bytes                 | gauge |
| iperf3_upload_sent_retransmits           | gauge |
| iperf3_upload_received_bits_per_second   | gauge |
| iperf3_upload_received_seconds           | gauge |
| iperf3_upload_received_bytes             | gauge |
