package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/asticode/go-astideepspeech"
	"github.com/cryptix/wav"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofrs/uuid"
	"github.com/mattetti/filebuffer"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/caarlos0/env"

	ds "github.com/asticode/go-astideepspeech"
	dj "github.com/kai5263499/diy-jarvis"
	pb "github.com/kai5263499/diy-jarvis/generated"
)

type config struct {
	MQTTBroker   string `env:"MQTT_BROKER"`
	MQTTClientID string `env:"MQTT_CLIENT_ID" envDefault:"deepspeech"`
	LogLevel     string `env:"LOG_LEVEL" envDefault:"info"`
	Model        string `env:"MODEL" envDefault:"/deepspeech_models/output_graph.pbmm"`
}

var (
	cfg       config
	mqttComms *dj.MqttComms
	model     *ds.Model
	channels  map[string]uint64
)

func mqttMessageHandler(client mqtt.Client, msg mqtt.Message) {
	logrus.Debugf("processing audio request with length %d from %s", len(msg.Payload()), msg.Topic())
	f := filebuffer.New(msg.Payload())

	r, err := wav.NewReader(f, int64(len(msg.Payload())))
	if err != nil {
		fmt.Printf("new reader error=%+#v\n", err)
		return
	}

	var d []int16
	for {
		s, err := r.ReadSample()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("sample read error=%+#v\n", err)
		}
		d = append(d, int16(s))
	}

	md, sttErr := model.SpeechToTextWithMetadata(d, 1)
	if sttErr != nil {
		logrus.WithError(sttErr).Error("error speech to text")
		return
	}
	defer md.Close()

	var output bytes.Buffer

	var confidence float64
	for ti, transcript := range md.Transcripts() {
		confidence = transcript.Confidence()
		for _, token := range transcript.Tokens() {
			output.WriteString(token.Text())
		}
		if ti < len(md.Transcripts())-1 {
			output.WriteString(" ")
		}
	}

	outputStr := output.String()
	logrus.Debugf("confidence=%f stt=%s", confidence, outputStr)

	b := pb.Base{
		Id:        uuid.Must(uuid.NewV4()).String(),
		Timestamp: uint64(time.Now().Unix()),
		Type:      pb.Type_TextRequestType,
		Text:      outputStr,
	}

	if err := mqttComms.SendRequest(b); err != nil {
		logrus.WithError(err).Error("error sending process audio response")
	}
}

func subscribeToRawAudioChannel(msg pb.Base) {
	if len(msg.SourceId) < 1 {
		logrus.Error("SourceId is empty, skipping registration")
		return
	}

	if _, found := channels[msg.SourceId]; found {
		logrus.Debugf("SourceId %s found in cache, skipping", msg.SourceId)
		return
	}

	logrus.Debugf("subscribing to SourceId %s", msg.SourceId)

	if token := mqttComms.MQTTClient().Subscribe(msg.SourceId, 0, mqttMessageHandler); token.Wait() && token.Error() != nil {
		logrus.WithError(token.Error()).Errorf("unable to subscribe to %s", msg.SourceId)
		return
	}

	logrus.Debugf("subscription to SourceId %s successful", msg.SourceId)

	channels[msg.SourceId] = msg.Timestamp
}

func main() {
	cfg = config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.WithError(err).Fatal("config parse")
	}

	if level, err := logrus.ParseLevel(cfg.LogLevel); err != nil {
		logrus.WithError(err).Fatal("parse log level")
	} else {
		logrus.SetLevel(level)
	}

	channels = make(map[string]uint64)

	var newModelErr error
	model, newModelErr = astideepspeech.New(cfg.Model)
	if newModelErr != nil {
		logrus.WithError(newModelErr).Fatal("unable to create new model")
	}
	defer model.Close()

	logrus.Info("Initialized DeepSpeech")

	var newMqttErr error
	mqttComms, newMqttErr = dj.NewMqttComms(cfg.MQTTClientID, cfg.MQTTBroker)
	if newMqttErr != nil {
		logrus.WithError(newMqttErr).Fatal("error creating mqtt comms")
	}

	logrus.Infof("Connected to MQTT broker %s with Client ID %s", cfg.MQTTBroker, cfg.MQTTClientID)

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		for {
			select {
			case msg := <-mqttComms.RequestChan():
				if msg.Type == pb.Type_RegisterAudioSourceRequestType {
					subscribeToRawAudioChannel(msg)
				}
			}
		}
	})

	logrus.Info("Listening for requests")

	g.Wait()
}
