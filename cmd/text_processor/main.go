package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"

	"github.com/caarlos0/env"
	dj "github.com/kai5263499/diy-jarvis"
	"github.com/kai5263499/diy-jarvis/domain"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"github.com/sirupsen/logrus"
)

type config struct {
	MQTTBroker      string `env:"MQTT_BROKER"`
	MQTTClientID    string `env:"MQTT_CLIENT_ID" envDefault:"textprocessor"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
	CommandSpecYaml string `env:"COMMAND_SPEC_YAML" envDefault:"commands.yaml"`
	Keyword         string `env:"KEYWORD", envDefault:"Jarvis"`
}

var (
	cfg       config
	mqttComms *dj.MqttComms
	commands  map[string]domain.TextEventCommand
)

func processText(evt *pb.Base) {
	if action, ok := commands[evt.Text]; ok {
		fmt.Printf("got command %s, performing action %+#v\n", evt.Text, action)
	} else {
		// Ignore command
	}
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

	fmt.Printf("Initialize Text Event Processor..\n")

	var yamlCommands domain.TextEventCommands
	yamlData, readFileErr := ioutil.ReadFile(cfg.CommandSpecYaml)
	if readFileErr != nil {
		logrus.WithError(readFileErr).Fatal("error reading file")
	}

	if err := yaml.Unmarshal([]byte(yamlData), &yamlCommands); err != nil {
		logrus.WithError(readFileErr).Fatal("error unmarshaling yaml")
	}

	commands = make(map[string]domain.TextEventCommand)
	for _, c := range yamlCommands.Commands {
		commands[c.Command] = c
	}

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		for {
			select {
			case msg := <-mqttComms.RequestChan():
				if msg.Type == pb.Type_TextRequestType {
					processText(&msg)
				}
			}
		}
	})

	g.Wait()

}
