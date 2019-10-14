# Build the diy-jarvis-builder image with go and other associated libraries
builder-image:
	docker build -t kai5263499/diy-jarvis-builder .

check-pulseaudio:
	pulseaudio --check -v

start-pulseaudio:
	pulseaudio --load="module-native-protocol-tcp auth-ip-acl=127.0.0.1;172.17.0.0/24 auth-anonymous=1" --exit-idle-time=-1 --daemon; true

# Run an interactive shell for development and testing
exec-interactive:
	docker pull kai5263499/diy-jarvis-builder
	docker run -it --rm \
	-e PULSE_SERVER=docker.for.mac.localhost \
	-v ~/.config/pulse:/home/pulseaudio/.config/pulse \
	-v ~/Downloads/deepspeech-0.5.1-models:/deepspeech_models \
	-v ~/code/deproot/src/github.com/kai5263499:/go/src/github.com/kai5263499 \
	-w /go/src/github.com/kai5263499/diy-jarvis \
	kai5263499/diy-jarvis-builder bash

# Run an image preconfigured with Mozilla Deep Speech and the latest English model
deepspeech-service:
	docker pull kai5263499/diy-jarvis-deepspeech
	docker run -p 6000:6000 -d \
	-e TEXT_PROCESSOR_ADDRESS="docker.for.mac.localhost:6001" \
	--name diy-jarvis-deepspeech \
	kai5263499/diy-jarvis-deepspeech

# Slice up a wav file (must be 16k sample rate and mono) and feed it to an audio processor (eg deepspeech-service)
wav-slicer:
	docker pull kai5263499/diy-jarvis-wav-slicer
	docker run -it --rm \
	-e AUDIO_PROCESSOR_ADDRESS="docker.for.mac.localhost:6000" \
	-e FILE=${FILE} \
	--mount type=tmpfs,destination=/tmp \
	-v ${DATA_DIR}:/data \
	kai5263499/diy-jarvis-wav-slicer

# Take a slice of sampled audio and feed it to the audio processor
mic-capture:
	docker pull kai5263499/diy-jarvis-mic-capture
	docker rm -f diy-jarvis-mic-capture; true
	docker run -t -d \
	-e DURATION=${DURATION} \
	-e AUDIO_PROCESSOR_ADDRESS="docker.for.mac.localhost:6000" \
	-e PULSE_SERVER=docker.for.mac.localhost \
	-v ~/.config/pulse:/home/pulseaudio/.config/pulse \
	--name diy-jarvis-mic-capture \
	kai5263499/diy-jarvis-mic-capture

text-processor:
	docker pull kai5263499/diy-jarvis-text-processor
	docker run -d \
	-e OUTPUT_PROCESSOR_ADDRESS="docker.for.mac.localhost:6002" \
	--name diy-jarvis-text-processor \
	kai5263499/diy-jarvis-text-processor

# Generate go stubs from proto definitions. This should be run inside of an interactive container
go-protos:
	protoc -I proto/ proto/*.proto --go_out=plugins=grpc:generated

.PHONY: image exec-interactive protos pulseaudio deepspeech-service mic_capture