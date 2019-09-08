# Build the diy-jarvis-builder image with go and other associated libraries
image:
	docker build -t kai5263499/diy-jarvis-builder .

# Run an interactive shell for development and testing
exec-interactive:
	docker run -it --rm \
	-e PULSE_SERVER=docker.for.mac.localhost \
	-v ~/.config/pulse:/home/pulseaudio/.config/pulse \
	-v ~/code/deproot/src/github.com/kai5263499:/go/src/github.com/kai5263499 \
	kai5263499/diy-jarvis-builder bash

# Generate go stubs from proto definitions. This should be run inside of an interactive container
protos:
	protoc -I proto/ proto/*.proto --go_out=plugins=grpc:generated

.PHONY: image exec-interactive protos