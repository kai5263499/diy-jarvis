package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/caarlos0/env"
	"github.com/gofrs/uuid"
	dj "github.com/kai5263499/diy-jarvis"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type config struct {
	MQTTBroker   string   `env:"MQTT_BROKER"`
	MQTTClientID string   `env:"MQTT_CLIENT_ID" envDefault:"slackbot"`
	LogLevel     string   `env:"LOG_LEVEL" envDefault:"info"`
	SlackToken   string   `env:"SLACK_TOKEN"`
	Channels     []string `env:"CHANNELS" envDefault:"jarvis"`
}

var (
	cfg           config
	slackClient   *slack.Client
	slackRtm      *slack.RTM
	slackChannel  *slack.Channel
	mqttComms     *dj.MqttComms
	sourceID      string
	slackChannels []*slack.Channel
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
	for {
		select {
		case msg := <-slackRtm.IncomingEvents:
			logrus.Debugf("event received %#v", msg)

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				logrus.Infof("Connection counter: %d", ev.ConnectionCount)
			case *slack.MessageEvent:
				logrus.Debugf("message event received!")
				go processMessage(ev)
			case *slack.RTMError:
				logrus.Errorf("RTMError: %s", ev.Error())
			case *slack.InvalidAuthEvent:
				logrus.Error("Invalid credentials!")
				return
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
		Type:     pb.Type_TextRequestType,
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

func joinChannels() {
	slackChannels = make([]*slack.Channel, len(cfg.Channels))
	for _, chanStr := range cfg.Channels {
		channel, err := slackClient.JoinChannel(chanStr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"channel": chanStr,
				"err":     err,
			}).Error("failed to join channel")
		} else {
			logrus.WithFields(logrus.Fields{
				"channel": chanStr,
			}).Debug("joined channel")
			slackChannels = append(slackChannels, channel)
		}
	}
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

	logrus.Debug("connecting to slack")

	slackClient = slack.New(cfg.SlackToken)
	slackRtm = slackClient.NewRTM()
	go slackRtm.ManageConnection()
	joinChannels()
	go slackReadLoop()

	logrus.Infof("connected to Slack with sourceID %s", sourceID)

	var newMqttErr error
	mqttComms, newMqttErr = dj.NewMqttComms(cfg.MQTTClientID, cfg.MQTTBroker)
	if newMqttErr != nil {
		logrus.WithError(newMqttErr).Fatal("error creating mqtt comms")
	}

	waitForCtrlC()
}
