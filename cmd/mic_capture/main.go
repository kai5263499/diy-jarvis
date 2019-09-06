package main

import (
	"os"
	"flag"
	"unsafe"
	"fmt"
	"time"
	"bytes"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"encoding/binary"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"github.com/cryptix/wav"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/xlab/portaudio-go/portaudio"
	"github.com/xlab/closer"
)

const (
	sampleRate  = 16000
	channels    = 1
	sampleFormat= portaudio.PaInt16
)

var (
	samplesPerChannel = uint(16000)
	duration = flag.Int("duration",10, "Number of seconds of audio to collect before sending to processor")
	audioProcessorAddress = flag.String("server", "localhost:6000", "Address of the audio processor to send samples to")

	numRequests uint32
	numResponses uint32
	stream pb.AudioProcessor_SubscribeClient
)

type processAudioDataRequest struct {
	uuid *uuid.UUID
	wavWriter *wav.Writer
	tmpFileName string
	msg *pb.ProcessAudioRequest
}

func checkError(err error) {
	if err != nil {
		panic(fmt.Sprintf("err=%#+v", err))
	}
}

func checkPAError(err portaudio.Error) {
	if portaudio.ErrorCode(err) != portaudio.PaNoError {
		panic(fmt.Sprintf("err=%s", paErrorText(err)))
	}
}

func paErrorText(err portaudio.Error) string {
	return portaudio.GetErrorText(err)
}

func newProcessAudioRequest() *processAudioDataRequest {
	newUUID, _ := uuid.NewV4()

	f, err := ioutil.TempFile("", newUUID.String())
	checkError(err)

	meta := wav.File{
		Channels:        channels,
		SampleRate:      sampleRate,
		SignificantBits: 16,
	}
	writer, _ := meta.NewWriter(f)

	msg := &pb.ProcessAudioRequest{
		RequestId: newUUID.String(),
		AudioStartTime: uint64(time.Now().Unix()),
	}

	return &processAudioDataRequest{
		uuid: newUUID,
		wavWriter: writer,
		tmpFileName: f.Name(),
		msg: msg,
	}
}

func processChunk(req *processAudioDataRequest) {
	var err error

	err = req.wavWriter.Close()
	checkError(err)

	content, err := ioutil.ReadFile(req.tmpFileName)
	checkError(err)

	req.msg.AudioData = content

	err = stream.Send(req.msg)
	atomic.AddUint32(&numRequests, 1)
	checkError(err)

	os.Remove(req.tmpFileName)
}

func processMicInput(wg *sync.WaitGroup) {
	defer wg.Done()
	defer closer.Close()

	var paErr portaudio.Error

	paErr = portaudio.Initialize()
	checkPAError(paErr)

	closer.Bind(func() {
		paErr := portaudio.Terminate()
		checkPAError(paErr)
	})
	
	var paStream *portaudio.Stream
	paErr = portaudio.OpenDefaultStream(&paStream, channels, 0, sampleFormat, sampleRate,
		samplesPerChannel, paCallback, nil)
	checkPAError(paErr)

	closer.Bind(func() {
		paErr := portaudio.CloseStream(paStream)
		checkPAError(paErr)
	})

	paErr = portaudio.StartStream(paStream)
	checkPAError(paErr)

	closer.Bind(func() {
		paErr := portaudio.StopStream(paStream)
		checkPAError(paErr)
	})

	closer.Hold()
}

func paCallback(input unsafe.Pointer, _ unsafe.Pointer, sampleCount uint,
	_ *portaudio.StreamCallbackTimeInfo, _ portaudio.StreamCallbackFlags, _ unsafe.Pointer) int32 {

	const (
		statusContinue = int32(portaudio.PaContinue)
		statusAbort    = int32(portaudio.PaAbort)
	)

	in := (*(*[1 << 24]int16)(input))[:int(sampleCount)*channels]

	req := newProcessAudioRequest()
	buf := new(bytes.Buffer)

	for frame := range in {
		binary.Write(buf, binary.LittleEndian, in[frame])
	}
	req.wavWriter.Write(buf.Bytes())

	processChunk(req)
	// if !ok {
	// 	return statusAbort
	// }
	
	return statusContinue
}

func processResponses(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		resp, err := stream.Recv()
		checkError(err)

		fmt.Printf("%s\n", resp.Output)
	}
}

func main() {
	flag.Parse()
	var err error
	
	conn, err := grpc.Dial(*audioProcessorAddress, grpc.WithInsecure())
	checkError(err)
	defer conn.Close()

	client := pb.NewAudioProcessorClient(conn)
	stream, err = client.Subscribe(context.Background())
	checkError(err)
	defer stream.CloseSend()

	samplesPerChannel = sampleRate * uint(*duration) 

	var wg sync.WaitGroup
	wg.Add(2)
	go processMicInput(&wg)
	go processResponses(&wg)
	wg.Wait()
}
