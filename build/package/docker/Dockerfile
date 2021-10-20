FROM alpine:3.14

RUN set -ex \
 && apk add --no-cache \
	iperf3
COPY iperf3exporter /

EXPOSE 9119

USER nobody

ENV IPERF3EXPORTER_EXPORTER_LISTEN=0.0.0.0:9119

ENTRYPOINT ["/iperf3exporter"]
