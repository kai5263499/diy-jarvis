FROM kai5263499/diy-jarvis-builder as builder

COPY / /go/src/github.com/kai5263499/diy-jarvis

RUN cd /go/src/github.com/kai5263499/diy-jarvis/cmd/mic_capture && \
    go mod vendor && \
    go build && \
    ldd mic_capture | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

FROM ubuntu:22.04

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

ENV LOG_LEVEL=info
ENV MQTT_BROKER=
ENV PULSE_DURATION=10s
ENV AUDIO_CAPTURE_DURATION=3s

COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/mic_capture/deps /
COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/mic_capture/mic_capture /mic_capture

# The || true here is to swallow an error caused when setting up systemd
# systemd attempts to symlink /etc/resolv.conf which is being actively managed by docker
RUN apt-get update && \
    apt-get install -y \
    alsa-utils \
	pulseaudio \
    pulseaudio-utils || true

ENTRYPOINT [ "/mic_capture" ]