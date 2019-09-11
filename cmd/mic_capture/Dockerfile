FROM kai5263499/diy-jarvis-builder as builder

COPY / /go/src/github.com/kai5263499/diy-jarvis

RUN cd /go/src/github.com/kai5263499/diy-jarvis/cmd/mic_capture && \
    go build && \
    ldd mic_capture | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

FROM ubuntu:18.04

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

ENV DURATION=10
ENV AUDIO_PROCESSOR_ADDRESS=""

COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/mic_capture/deps /
COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/mic_capture/mic_capture /mic_capture

RUN apt-get update && \
    apt-get install -y \
    alsa-utils \
	pulseaudio \
    pulseaudio-utils 

ENTRYPOINT [ "/mic_capture" ]