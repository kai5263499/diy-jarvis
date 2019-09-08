FROM ubuntu:18.04
LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

ENV CGO_ENABLED=1 CGO_CPPFLAGS="-I/usr/include"
ENV GOPATH=/go
ENV PATH=/go/bin:/deepspeech:$PATH
ENV CGO_LDFLAGS="-L/deepspeech"
ENV CGO_CXXFLAGS="-I/deepspeech"
ENV LD_LIBRARY_PATH=/deepspeech:$LD_LIBRARY_PATH

RUN apt-get update && \
    apt-get install -y \
    git \
    golang \
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
    portaudio19-dev && \
	# Install deepspeech
	mkdir -p /deepspeech && \
	cd /deepspeech && \
	wget https://github.com/mozilla/DeepSpeech/raw/v0.5.1/native_client/deepspeech.h && \
	wget https://github.com/mozilla/DeepSpeech/releases/download/v0.5.1/native_client.amd64.cpu.linux.tar.xz && \
	tar -xvf native_client.amd64.cpu.linux.tar.xz && \
	rm native_client.amd64.cpu.linux.tar.xz && \
	go get -u github.com/asticode/go-astideepspeech/... && \
	# Misc golang libraries
	go get github.com/xlab/portaudio-go/portaudio && \
	go get github.com/xlab/closer && \
	go get github.com/zenwerk/go-wave && \
	go get github.com/nu7hatch/gouuid && \
	go get golang.org/x/net/context && \
	go get -u google.golang.org/grpc && \
	# Install protoc tools
	go get -u github.com/golang/protobuf/protoc-gen-go && \
	curl -sLO https://github.com/google/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip && \
    unzip protoc-3.7.1-linux-x86_64.zip -d protoc3 && \
    mv protoc3/bin/* /usr/local/bin/ && \
    mv protoc3/include/* /usr/local/include/ && \
    rm -rf protoc3 protoc-3.7.1-linux-x86_64.zip