FROM alpine:3.14

RUN set -ex \
 && apk add --no-cache \
	iperf3
COPY iperf3exporter /

EXPOSE 9119

USER nobody

ENTRYPOINT ["/iperf3exporter"]
