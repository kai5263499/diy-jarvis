package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/caarlos0/env"
	"github.com/kai5263499/diy-jarvis/domain"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"github.com/nlopes/slack"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
)

type config struct {
	SlackToken           string `env:"SLACK_TOKEN"`
	TextProcessorAddress string `env:"TEXT_PROCESSOR_ADDRESS" envDefault:"localhost:6001"`
}

var (
	cfg          config
	stream       pb.EventResponder_SubscribeClient
	slackClient  *slack.Client
	slackRtm     *slack.RTM
	slackChannel *slack.Channel
)

func processFileUpload(ev *slack.FileSharedEvent) {
	var err error

	file, _, _, err := slackClient.GetFileInfo(ev.File.ID, 0, 0)
	domain.CheckError(err)

	client := &http.Client{}
	req, err := http.NewRequest("GET", file.URLPrivate, nil)
	domain.CheckError(err)
	req.Header.Set("Authorization", "Bearer "+cfg.SlackToken)
	resp, err := client.Do(req)
	domain.CheckError(err)

	body, err := ioutil.ReadAll(resp.Body)
	domain.CheckError(err)

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
			fmt.Printf("event received %#v\n", msg)

			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Printf("Connection counter: %d", ev.ConnectionCount)
			case *slack.MessageEvent:
				go processMessage(ev)
			case *slack.FileSharedEvent:
				go processFileUpload(ev)
			case *slack.RTMError:
				fmt.Printf("RTMError: %s\n", ev.Error())
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials!\n")
				break Loop
			default:
				//Take no action
			}
		}
	}
}

func sendRequest(input string) {
	newUUID := uuid.Must(uuid.NewV4())

	req := &pb.TextEventRequest{
		RequestId: newUUID.String(),
		Text:      input,
	}

	stream.Send(req)
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
	var err error

	cfg = config{}
	err = env.Parse(&cfg)
	domain.CheckError(err)

	fmt.Printf("Connecting to Slack..\n")

	slackClient = slack.New(cfg.SlackToken)
	slackClient.SetUserAsActive()
	slackRtm = slackClient.NewRTM()

	fmt.Printf("Connecting to text processor\n")

	conn, err := grpc.Dial(cfg.TextProcessorAddress, grpc.WithInsecure())
	domain.CheckError(err)
	defer conn.Close()

	client := pb.NewEventResponderClient(conn)
	stream, err = client.Subscribe(context.Background())
	domain.CheckError(err)
	defer stream.CloseSend()

	fmt.Printf("Beginning read loop\n")

	go slackRtm.ManageConnection()
	go slackReadLoop()

	waitForCtrlC()
}
