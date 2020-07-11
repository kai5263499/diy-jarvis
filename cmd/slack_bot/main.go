package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/caarlos0/env"
	"github.com/gofrs/uuid"
	dj "github.com/kai5263499/diy-jarvis"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type config struct {
	MQTTBroker   string `env:"MQTT_BROKER"`
	MQTTClientID string `env:"MQTT_CLIENT_ID" envDefault:"slackbot"`
	LogLevel     string `env:"LOG_LEVEL" envDefault:"info"`
	SlackToken   string `env:"SLACK_TOKEN"`
}

var (
	cfg          config
	slackClient  *slack.Client
	slackRtm     *slack.RTM
	slackChannel *slack.Channel
	mqttComms    *dj.MqttComms
	sourceID     string
)

func processFileUpload(ev *slack.FileSharedEvent) {
	file, _, _, getFileInfoErr := slackClient.GetFileInfo(ev.File.ID, 0, 0)
	if getFileInfoErr != nil {
		logrus.WithError(getFileInfoErr).Error("error getting file info")
		return
	}

	client := &http.Client{}
	req, newRequestErr := http.NewRequest("GET", file.URLPrivate, nil)
	if newRequestErr != nil {
		logrus.WithError(newRequestErr).Error("error creating file request")
		return
	}

	req.Header.Set("Authorization", "Bearer "+cfg.SlackToken)
	resp, clientDoErr := client.Do(req)
	if clientDoErr != nil {
		logrus.WithError(clientDoErr).Error("error performing client request")
		return
	}

	body, readAllErr := ioutil.ReadAll(resp.Body)
	if readAllErr != nil {
		logrus.WithError(readAllErr).Error("error reading response")
		return
	}

	bodyStr := string(body)
	sendRequest(bodyStr)
}

func processMessage(ev *slack.MessageEvent) {
	sendRequest(ev.Text)
}

func slackReadLoop() {
Loop:
	for {
		select {
		case msg := <-slackRtm.IncomingEvents:
			logrus.Debugf("event received %#v", msg)

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				logrus.Debugf("Connection counter: %d", ev.ConnectionCount)
			case *slack.MessageEvent:
				go processMessage(ev)
			case *slack.FileSharedEvent:
				go processFileUpload(ev)
			case *slack.RTMError:
				logrus.Warnf("RTMError: %s\n", ev.Error())
			case *slack.InvalidAuthEvent:
				logrus.Errorf("Invalid credentials!\n")
				break Loop
			default:
				//Take no action
			}
		}
	}
}

func sendRequest(input string) {
	newUUID := uuid.Must(uuid.NewV4())

	req := pb.Base{
		Id:       newUUID.String(),
		Text:     input,
		SourceId: sourceID,
	}

	mqttComms.SendRequest(req)
}

func waitForCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalCh chan os.Signal
	signalCh = make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		<-signalCh
		endWaiter.Done()
	}()
	endWaiter.Wait()
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

	sourceID = uuid.Must(uuid.NewV4()).String()

	logrus.Info("Connecting to slack")

	slackClient = slack.New(cfg.SlackToken)
	slackClient.SetUserAsActive()
	slackRtm = slackClient.NewRTM()

	var newMqttErr error
	mqttComms, newMqttErr = dj.NewMqttComms(cfg.MQTTClientID, cfg.MQTTBroker)
	if newMqttErr != nil {
		logrus.WithError(newMqttErr).Fatal("error creating mqtt comms")
	}

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		slackReadLoop()
		return nil
	})

	g.Go(func() error {
		slackRtm.ManageConnection()
		return nil
	})

	g.Wait()

	waitForCtrlC()
}
