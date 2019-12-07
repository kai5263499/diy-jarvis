FROM ubuntu:18.04 as models

# DEEPSPEECH_VERSION is the version of deepspeech to use
ARG DEEPSPEECH_VERSION=0.6.0

RUN apt-get update && \
    apt-get install -y wget && \
    mkdir -p /deepspeech_models

RUN wget https://github.com/mozilla/DeepSpeech/releases/download/v${DEEPSPEECH_VERSION}/deepspeech-${DEEPSPEECH_VERSION}-models.tar.gz -O deepepeech-models.tar.gz -nv && \
    tar -zxf deepepeech-models.tar.gz --strip-components=1 -C /deepspeech_models

FROM scratch

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

COPY --from=models /deepspeech_models /deepspeech_models