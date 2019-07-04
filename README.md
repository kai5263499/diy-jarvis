# diy-jarvis

This repo contains various scripts and notebooks I've come across for processing audio in general and speech in particular. The overall goal of all this is to create a simple yet flexible system for responding to audio events.



## Containerized development

I've found that working in a containerized development environment helps me make my finished product more portable. In order to do that, we need to run pulseaudio on the host and connect it to the container.

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

At this point your environment should be configured to run audio processing scripts and applications using pulseaudio inside of a docker container. This will make development easier and more reproducable when we go to transfer our final application to an embedded system such as a Raspberry Pi.
