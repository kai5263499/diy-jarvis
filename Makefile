# Build the diy-jarvis-builder image with go and other associated libraries
builder-image:
	docker build -t kai5263499/diy-jarvis-builder .

# Should tell you the pulseaudio daemon is running
check-pulseaudio:
	pulseaudio --check -v

# Helpful for Mac hosts
start-pulseaudio:
	pulseaudio --load="module-native-protocol-tcp auth-ip-acl=127.0.0.1;172.17.0.0/24 auth-anonymous=1" --exit-idle-time=-1 --daemon; true

# Build deepspeech model image
deepspeech-models-image:
	docker build -t kai5263499/diy-jarvis-deepspeech-models -f cmd/deepspeech/Dockerfile.model .

# Build deepspeech model image
deepspeech-image:
	docker build -t kai5263499/diy-jarvis-deepspeech -f cmd/deepspeech/Dockerfile .

# Build mic-capture image
mic-capture-image:
	docker build -t kai5263499/diy-jarvis-mic-capture -f cmd/mic_capture/Dockerfile .

# Build wav-slicer image
wav-slicer-image:
	docker build -t kai5263499/diy-jarvis-wav-slicer -f cmd/wav_slicer/Dockerfile .

# Run an interactive shell for development and testing
exec-interactive:
	docker run -it --rm \
	-e PULSE_SERVER=${PULSE_SERVER} \
	-v ${PULSE_CONFIG}:/home/pulseaudio/.config/pulse \
	-v ${DEEPSPEECH_MODELS}:/deepspeech_models \
	-v ~/code/deproot/src/github.com/kai5263499:/go/src/github.com/kai5263499 \
	-w /go/src/github.com/kai5263499/diy-jarvis \
	kai5263499/diy-jarvis-builder bash

# Run an image preconfigured with Mozilla Deep Speech and the latest English model
run-deepspeech:
	docker run -p 6000:6000 -d \
	-e LOG_LEVEL=debug \
	-e MQTT_BROKER=tcp://172.17.0.1:1883 \
	--name diy-jarvis-deepspeech \
	kai5263499/diy-jarvis-deepspeech

# Slice up a wav file (must be 16k sample rate and mono) and feed it to an audio processor (eg deepspeech-service)
run-wav-slicer:
	docker run -it --rm \
	-e FILE=${FILE} \
	--mount type=tmpfs,destination=/tmp \
	-v ${DATA_DIR}:/data \
	kai5263499/diy-jarvis-wav-slicer

# Take a slice of sampled audio and feed it to the audio processor
run-mic-capture:
	docker run -t -d \
	-e LOG_LEVEL=debug \
	-e MQTT_BROKER=tcp://172.17.0.1:1883 \
	-e DURATION=${DURATION} \
	-e PULSE_SERVER=${PULSE_SERVER} \
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
