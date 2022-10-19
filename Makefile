DOCKER_REPO := kai5263499

MODULES=$(subst cmd/,,$(wildcard cmd/*))
TEST_FILTER?=.

GIT_COMMIT := $(shell git rev-parse HEAD | cut -c 1-8)
GIT_BRANCH := $(shell git branch --show-current)

which = $(shell which $1 2> /dev/null || echo $1)

GO_PATH := $(call which,go)
$(GO_PATH):
	$(error Missing go)

# Build the diy-jarvis-builder image with go and other associated libraries
image/builder:
	docker build -t kai5263499/diy-jarvis-builder .

# Should tell you the pulseaudio daemon is running
pulseaudio/check:
	pulseaudio --check -v

# Helpful for Mac hosts
pulseaudio/start:
	pulseaudio --load="module-native-protocol-tcp auth-ip-acl=127.0.0.1;172.17.0.0/24 auth-anonymous=1" --exit-idle-time=-1 --daemon; true

image/deepspeech-models:
	docker build -t kai5263499/diy-jarvis-deepspeech-models -f cmd/deepspeech/Dockerfile.model .

image/deepspeech:
	docker build -t kai5263499/diy-jarvis-deepspeech -f cmd/deepspeech/Dockerfile .

image/mic_capture:
	docker build -t kai5263499/diy-jarvis-mic-capture -f cmd/mic_capture/Dockerfile .

image/wav-slicer:
	docker build -t kai5263499/diy-jarvis-wav-slicer -f cmd/wav_slicer/Dockerfile .

image/text_processor:
	docker build -t kai5263499/diy-jarvis-text-processor -f cmd/text_processor/Dockerfile .

image/slack_bot:
	docker build -t kai5263499/diy-jarvis-slack-bot -f cmd/slack_bot/Dockerfile .

images: image/builder image/deepspeech-models image/deepspeech image/mic_capture image/wav-slicer image/text_processor image/slack_bot

# Run an interactive shell for development and testing
debug:
	docker run -it --rm \
	--net=host --add-host host.docker.internal:host-gateway \
	-e PULSE_SERVER=host.docker.internal \
	-v ${PULSE_CONFIG}:/home/pulseaudio/.config/pulse \
	-v ${DEEPSPEECH_MODELS}:/deepspeech_models \
	-v ~/code/deproot/src/github.com/kai5263499:/go/src/github.com/kai5263499 \
	-w /go/src/github.com/kai5263499/diy-jarvis \
	kai5263499/diy-jarvis-builder bash

# Run an image preconfigured with Mozilla Deep Speech and the latest English model
run/deepspeech:
	docker run -p 6000:6000 -d \
	-e LOG_LEVEL=debug \
	--net=host --add-host host.docker.internal:host-gateway \
	-e MQTT_BROKER=tcp://host.docker.internal:1883 \
	--name diy-jarvis-deepspeech \
	kai5263499/diy-jarvis-deepspeech

# Slice up a wav file (must be 16k sample rate and mono) and feed it to an audio processor (eg deepspeech-service)
run/wav-slicer:
	docker run -it --rm \
	-e FILE=${FILE} \
	--mount type=tmpfs,destination=/tmp \
	-v ${DATA_DIR}:/data \
	kai5263499/diy-jarvis-wav-slicer

# Take a slice of sampled audio and feed it to the audio processor
run/mic-capture:
	docker run -t -d \
	--net=host --add-host host.docker.internal:host-gateway \
	-e PULSE_SERVER=host.docker.internal \
	-e LOG_LEVEL=debug \
	-e MQTT_BROKER=tcp://host.docker.internal:1883 \
	-e DURATION=3s \
	-v ${PULSE_CONFIG}:/home/pulseaudio/.config/pulse \
	--name diy-jarvis-mic-capture \
	kai5263499/diy-jarvis-mic-capture

run-text-processor:
	docker run -d \
	--name diy-jarvis-text-processor \
	kai5263499/diy-jarvis-text-processor

run-mqtt-mosquito:
	docker run -d \
	-p 1883:1883 -p 9001:9001 \
	--name mosquitto \
	eclipse-mosquitto

# Generate go stubs from proto definitions. This should be run inside of an interactive container
gen-protos:
	protoc -I proto/ proto/*.proto --go_out=generated

LINTER_PATH := $(call which,golangci-lint)
$(LINTER_PATH):
	$(error Missing golangci: https://golangci-lint.run/usage/install)
lint:
	@rm -rf ./vendor
	@$(GO_PATH) mod vendor
	export GOMODCACHE=./vendor
	golangci-lint run ./...