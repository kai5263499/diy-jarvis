FROM kai5263499/diy-jarvis-builder as builder

COPY / /go/src/github.com/kai5263499/diy-jarvis

RUN cd /go/src/github.com/kai5263499/diy-jarvis/cmd/slack_bot && \
    go mod vendor && \
    go build && \
    ldd slack_bot | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

FROM scratch

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

ENV LOG_LEVEL=info
ENV MQTT_BROKER=
ENV SLACK_TOKEN=""

COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/slack_bot/deps /
COPY --from=builder /go/src/github.com/kai5263499/diy-jarvis/cmd/slack_bot/slack_bot /slack_bot

ENTRYPOINT [ "/slack_bot" ]