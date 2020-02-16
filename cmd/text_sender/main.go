package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/caarlos0/env"
	"github.com/gofrs/uuid"
	"github.com/kai5263499/diy-jarvis/domain"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type config struct {
	TextProcessorAddress string `env:"TEXT_PROCESSOR_ADDRESS" envDefault:"localhost:6001"`
}

var (
	cfg          config
	stdInStopped int32
	numRequests  uint32
	numResponses uint32

	stream pb.EventResponder_SubscribeClient
)

func sendRequest(input string) {
	newUUID := uuid.Must(uuid.NewV4())

	req := &pb.TextEventRequest{
		RequestId: newUUID.String(),
		Text:      input,
	}

	stream.Send(req)
	atomic.AddUint32(&numRequests, 1)
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

	atomic.CompareAndSwapInt32(&stdInStopped, 0, 1)
}

func processResponses(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		resp, err := stream.Recv()
		domain.CheckError(err)

		atomic.AddUint32(&numResponses, 1)

		if resp.ResponseCode != pb.TextEventResponse_ACCEPTED {
			fmt.Printf("RequestId %s not accepted!\n", resp.RequestId)
		}

		stdInStopped := atomic.LoadInt32(&stdInStopped)
		requests := atomic.LoadUint32(&numRequests)
		responses := atomic.LoadUint32(&numResponses)

		if stdInStopped == 1 && requests == responses {
			return
		}
	}
}

func main() {
	var err error

	cfg = config{}
	err = env.Parse(&cfg)
	domain.CheckError(err)

	conn, err := grpc.Dial(cfg.TextProcessorAddress, grpc.WithInsecure())
	domain.CheckError(err)
	defer conn.Close()

	client := pb.NewEventResponderClient(conn)
	stream, err = client.Subscribe(context.Background())
	domain.CheckError(err)
	defer stream.CloseSend()

	var wg sync.WaitGroup
	wg.Add(2)
	go processInput(&wg)
	go processResponses(&wg)
	wg.Wait()
}
