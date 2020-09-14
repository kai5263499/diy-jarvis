package main

import (
	"os"

	"github.com/caarlos0/env"
	"github.com/gofrs/uuid"
	dj "github.com/kai5263499/diy-jarvis"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"github.com/sirupsen/logrus"
)

type config struct {
	MQTTBroker   string `env:"MQTT_BROKER"`
	MQTTClientID string `env:"MQTT_CLIENT_ID" envDefault:"textsender"`
	LogLevel     string `env:"LOG_LEVEL" envDefault:"info"`
}

var (
	cfg          config
	stdInStopped int32
	numRequests  uint32
	numResponses uint32
	mqttComms    *dj.MqttComms
	sourceID     string
)

func main() {
	cfg = config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.WithError(err).Fatal("parse configs")
	}

	if level, err := logrus.ParseLevel(cfg.LogLevel); err != nil {
		logrus.WithError(err).Fatal("parse log level")
	} else {
		logrus.SetLevel(level)
	}

	sourceID = uuid.Must(uuid.NewV4()).String()

	var newMqttErr error
	mqttComms, newMqttErr = dj.NewMqttComms(cfg.MQTTClientID, cfg.MQTTBroker)
	if newMqttErr != nil {
		logrus.WithError(newMqttErr).Fatal("error creating mqtt comms")
	}

	if len(os.Args) < 2 {
		logrus.Fatal("must include text to send")
	}

	input := os.Args[1]

	req := pb.Base{
		Id:       uuid.Must(uuid.NewV4()).String(),
		Type:     pb.Type_TextRequestType,
		SourceId: sourceID,
		Text:     input,
	}

	if err := mqttComms.SendRequest(req); err != nil {
		logrus.WithError(err).Fatal("error sending request")
	}

	mqttComms.Close()
}
