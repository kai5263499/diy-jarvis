# gortana

Gortana is a shameless clone of the example code included in the golang pulseaudio repo. It's based on the pocketsphinx library which is targeted to work on embedded systems. Because these embedded systems have limited resources, we also use a keyword list to reduce the search space that pocketsphinx needs to search through.

~~~~bash
# Start by running an image that contains pocketsphinx and it's golang wrapper along with the pulseaudio subsystem
docker run -it --rm -e PULSE_SERVER=docker.for.mac.localhost -v ~/.config/pulse:/home/pulseaudio/.config/pulse -v ~/code/deproot/src/github.com/kai5263499/diy-jarvis:/go/src/github.com/kai5263499/diy_jarvis pocketsphinx-go bash

# Then, run pocketsphinx without keywords. This should illustrate how inefficient pocketsphinx is when searching through the entire english dictionary.
pocketsphinx_continuous -hmm /usr/share/pocketsphinx/model/en-us/en-us -lm /usr/share/pocketsphinx/model/en-us/en-us.lm.bin -dict /usr/share/pocketsphinx/model/en-us/cmudict-en-us.dict -inmic yes

# Next, run pocketsphinx with a limited set of keywords
pocketsphinx_continuous -hmm /usr/share/pocketsphinx/model/en-us/en-us -lm /go/src/github.com/kai5263499/diy_jarvis/commands/6087.lm -dict /go/src/github.com/kai5263499/diy_jarvis/commands/6087.dic -keyphrase "JARVIS" -kws_threshold 1e-20 -inmic yes

# Finally, run gortana using our limited set of keywords 
go run ~/dep-root/src/github.com/kai5263499/diy-jarvis/gortana/main.go --hmm=/usr/share/pocketsphinx/model/en-us/en-us --dict=~/dep-root/src/github.com/kai5263499/diy-jarvis/commands/6087.dic --lm=~/dep-root/src/github.com/kai5263499/diy-jarvis/commands/6087.lm --stdout

# TODO: Get this to cross-compile properly for ARM (Raspberry Pi)
cd ~/dep-root/src/github.com/kai5263499/diy-jarvis/gortana
GOOS=linux GOARCH=arm GOARM=5 go build

~~~~
