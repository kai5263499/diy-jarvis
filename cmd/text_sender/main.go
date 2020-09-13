package main

import (
	"bufio"
	"io"
	"os"
	"sync"

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

func sendRequest(input string) {
	req := pb.Base{
		Id:       uuid.Must(uuid.NewV4()).String(),
		Type:     pb.Type_TextRequestType,
		SourceId: sourceID,
		Text:     input,
	}

	mqttComms.SendRequest(req)
}

func processInput(wg *sync.WaitGroup) {
	defer wg.Done()

	var err error
	var input string

	reader := bufio.NewReader(os.Stdin)
	for {
		input, err = reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		sendRequest(input)
	}

	sendRequest(input)
}

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

	var wg sync.WaitGroup
	wg.Add(1)
	go processInput(&wg)
	wg.Wait()
}
