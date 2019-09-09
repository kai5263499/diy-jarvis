FROM kai5263499/diy-jarvis-builder as builder

COPY / /go/src/github.com/kai5263499/diy-jarvis

RUN cd /go/src/github.com/kai5263499/diy-jarvis/cmd/wav_slicer && \
    go build && \
    ldd wav_slicer | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

FROM scratch

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

ENV FILE=""
ENV AUDIO_PROCESSOR_ADDRESS=""

COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/wav_slicer/deps /
COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/wav_slicer/wav_slicer /wav_slicer

ENTRYPOINT [ "/wav_slicer" ]