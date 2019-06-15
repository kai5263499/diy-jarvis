# diy-jarvis

~~~~bash
# Check if pulse audio is running on the host
pulseaudio --check -v
~~~~

~~~~bash
# Start pulse audio daemon on the host allowing anonymous connections from the docker ip range
pulseaudio --load="module-native-protocol-tcp auth-ip-acl=127.0.0.1;172.17.0.0/24 auth-anonymous=1" --exit-idle-time=-1 --daemon
~~~~

~~~~bash
# mic check - 2 second delay from default in to default out
pacat -r | pacat -p --latency-msec=2000
~~~~

~~~~bash
# run container with diy_jarvis code and pulseaudio TCP connection to host
docker run -it --rm -e PULSE_SERVER=docker.for.mac.localhost -v ~/.config/pulse:/home/pulseaudio/.config/pulse -v ~/code/deproot/src/github.com/kai5263499/diy-jarvis:/go/src/github.com/kai5263499/diy_jarvis pocketsphinx-go bash
~~~~

~~~~bash
# run pocketsphinx without keywords
pocketsphinx_continuous -hmm /usr/share/pocketsphinx/model/en-us/en-us -lm /usr/share/pocketsphinx/model/en-us/en-us.lm.bin -dict /usr/share/pocketsphinx/model/en-us/cmudict-en-us.dict -inmic yes
~~~~

~~~~bash
# run pocketsphinx with keywords
pocketsphinx_continuous -hmm /usr/share/pocketsphinx/model/en-us/en-us -lm /go/src/github.com/kai5263499/diy_jarvis/commands/6087.lm -dict /go/src/github.com/kai5263499/diy_jarvis/commands/6087.dic -keyphrase "JARVIS" -kws_threshold 1e-20 -inmic yes
~~~~

~~~~bash
# run gortana with keywords
go run src/github.com/kai5263499/diy_jarvis/gortana/main.go --hmm=/usr/share/pocketsphinx/model/en-us/en-us --dict=/go/src/github.com/kai5263499/diy_jarvis/commands/6087.dic --lm=/go/src/github.com/kai5263499/diy_jarvis/commands/6087.lm --stdout
~~~~