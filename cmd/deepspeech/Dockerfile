FROM kai5263499/diy-jarvis-deepspeech-models as models
FROM kai5263499/diy-jarvis-builder as builder

COPY / /go/src/github.com/kai5263499/diy-jarvis

RUN cd /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech && \
    go build && \
    ldd deepspeech | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

RUN cd /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech && \
    mv deps/deepspeech/libdeepspeech.so deps/lib/libdeepspeech.so && \
    rm -rf deps/deepspeech

FROM ubuntu:18.04

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

ENV LOG_LEVEL=info
ENV MQTT_BROKER=
ENV MODEL="/deepspeech_models/deepspeech.pbmm"

COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech/deps /
COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/deepspeech/deepspeech /deepspeech
COPY --from=models /deepspeech_models /deepspeech_models

EXPOSE 6000

ENTRYPOINT [ "/deepspeech" ]