FROM ubuntu:18.04

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

# GPU_PLATFORM is the location where model processing takes place. cpu or cuda
ARG GPU_PLATFORM=cpu

# DEEPSPEECH_VERSION is the version of deepspeech to use
ARG DEEPSPEECH_VERSION=0.6.0

ENV CGO_ENABLED=1 CGO_CPPFLAGS="-I/usr/include"
ENV GOPATH=/go
ENV PATH=/go/bin:/usr/local/go/bin:/deepspeech:$PATH
ENV CGO_LDFLAGS="-L/deepspeech"
ENV CGO_CXXFLAGS="-I/deepspeech"
ENV LD_LIBRARY_PATH=/deepspeech:$LD_LIBRARY_PATH
ENV DEBIAN_FRONTEND=noninteractive

COPY . /go/src/github.com/kai5263499/diy-jarvis

WORKDIR /go/src/github.com/kai5263499/diy-jarvis

RUN apt-get update && \
    apt-get install -y \
    git \
    curl \
	unzip \
	wget \
	sox \
	ffmpeg \
	python3-pip \
	alsa-utils \
	pulseaudio \
	pulseaudio-utils \
	libsoxr-dev \
	tzdata \
    portaudio19-dev

RUN echo "Install golang" && \
	curl -sLO https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz && \
	tar -xf go1.13.3.linux-amd64.tar.gz && \
	mv go /usr/local && \
	rm -rf go1.13.3.linux-amd64.tar.gz

RUN	echo "Install deepspeech" && \
	mkdir -p /deepspeech && \
	cd /deepspeech && \
	wget https://github.com/mozilla/DeepSpeech/raw/v${DEEPSPEECH_VERSION}/native_client/deepspeech.h && \
	wget https://github.com/mozilla/DeepSpeech/releases/download/v${DEEPSPEECH_VERSION}/native_client.amd64.${GPU_PLATFORM}.linux.tar.xz && \
	tar -xvf native_client.amd64.${GPU_PLATFORM}.linux.tar.xz && \
	rm native_client.amd64.${GPU_PLATFORM}.linux.tar.xz

RUN echo "Caching golang modules" && \
	go mod vendor

RUN	echo "Install protoc tools" && \
	go get -u github.com/golang/protobuf/protoc-gen-go && \
	curl -sLO https://github.com/google/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip && \
    unzip protoc-3.7.1-linux-x86_64.zip -d protoc3 && \
    mv protoc3/bin/* /usr/local/bin/ && \
    mv protoc3/include/* /usr/local/include/ && \
    rm -rf protoc3 protoc-3.7.1-linux-x86_64.zip