# diy-jarvis

This repo contains the services that make up a simple DIY voice assistant framework. The components of this framework are broken up into docker images for easier development. These images are

| Image | Description | Parameters |
| ----------- | ----------- | ----------- |
| diy-jarvis-builder | Used for development and as the build container in our CI/CD pipeline ||
| diy-jarvis-mic-capture | Captures a number of seconds of microphone data (determined by the `DURATION` environment variable)| `DURATION` controls how much audio to collect before sending it on to the processor|
| diy-jarvis-deepspeech | Uses the Mozilla Deep Speech library to transform raw wave files into ||
| diy-jarvis-wav-slicer | Slices up a wav `FILE` and feeds the chunks to the audio processor engine | `FILE` is the full path of the file (in the container) to process| 

## Getting started

The quickest way to get started with diy-jarvis is to use the Makefile at the root of our project to pull and execute the component service images. The simplest setup includes the following

~~~~bash
make deepspeech-service

export DURATION=3
make mic-capture
docker logs -f diy-jarvis-mic-capture
~~~~

## Using the diy-jarvis-builder image

This image includes all of the tooling required to build the various services in the diy-jarvis ecosystem.

~~~~bash
# Pull down the latest builder image
docker pull kai5263499/diy-jarvis-builder

# Run the builder image with pulse configured for localhost and
# mounted development directories
docker pull kai5263499/diy-jarvis-builder
	docker run -it --rm \
	-e PULSE_SERVER=docker.for.mac.localhost \
	-v ~/.config/pulse:/home/pulseaudio/.config/pulse \
	-v ~/Downloads/deepspeech-0.5.1-models:/deepspeech_models \
	-v ~/code/deproot/src/github.com/kai5263499:/go/src/github.com/kai5263499 \
	-w /go/src/github.com/kai5263499/diy-jarvis \
	kai5263499/diy-jarvis-builder bash
~~~~

## Containerized development with PulseAudio

We've found that working in a containerized development environment helps us make our finished product more portable. In order to do that, we need to run pulseaudio on the host and connect it to the container.

~~~~bash
# First, we need to check if pulse audio is running on the host
pulseaudio --check -v

# Most likely, its not running so we'll need to start pulse audio daemon on the host allowing anonymous connections from the docker ip range, assuming it's 172.17.0.0/24 which appears to be the default for Docker Desktop on my mac
pulseaudio --load="module-native-protocol-tcp auth-ip-acl=127.0.0.1;172.17.0.0/24 auth-anonymous=1" --exit-idle-time=-1 --daemon

# Now we can run a basic container that has pulseaudio installed to test our audio setup
docker run -it -e PULSE_SERVER=docker.for.mac.localhost -v ~/.config/pulse:/home/pulseaudio/.config/pulse --entrypoint bash --rm jess/pulseaudio

# We then need to set the default source and sink to and run a mic check with a 2 second delay from our selected default source (in) to default sink (out) to make sure everything's in order
pacat -r | pacat -p --latency-msec=2000
~~~~
