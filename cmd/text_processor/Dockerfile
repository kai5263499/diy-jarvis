FROM kai5263499/diy-jarvis-builder as builder

COPY / /go/src/github.com/kai5263499/diy-jarvis

RUN cd /go/src/github.com/kai5263499/diy-jarvis/cmd/text_processor && \
    go mod vendor && \
    go build && \
    ldd text_processor | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

FROM ubuntu:22.04

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

ENV LOG_LEVEL=info
ENV MQTT_BROKER=
ENV COMMAND_SPEC_YAML=commands.yaml

EXPOSE 6001

COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/text_processor/deps /
COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/text_processor/text_processor /text_processor

ENTRYPOINT [ "/text_processor" ]