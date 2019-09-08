# Build the diy-jarvis-builder image with go and other associated libraries
image:
	docker build -t kai5263499/diy-jarvis-builder .

pulseaudio:
	pulseaudio --load="module-native-protocol-tcp auth-ip-acl=127.0.0.1;172.17.0.0/24 auth-anonymous=1" --exit-idle-time=-1 --daemon
	pulseaudio --check -v

# Run an interactive shell for development and testing
exec-interactive:
	docker run -it --rm \
	-e PULSE_SERVER=docker.for.mac.localhost \
	-v ~/.config/pulse:/home/pulseaudio/.config/pulse \
	-v ~/Downloads/deepspeech-0.5.1-models:/deepspeech_models \
	-v ~/code/deproot/src/github.com/kai5263499:/go/src/github.com/kai5263499 \
	-w /go/src/github.com/kai5263499/diy-jarvis \
	kai5263499/diy-jarvis-builder bash

# Generate go stubs from proto definitions. This should be run inside of an interactive container
protos:
	protoc -I proto/ proto/*.proto --go_out=plugins=grpc:generated

.PHONY: image exec-interactive protos pulseaudio