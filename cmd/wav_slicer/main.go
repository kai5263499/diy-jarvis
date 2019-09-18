package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/caarlos0/env"
	"github.com/cryptix/wav"
	"github.com/kai5263499/diy-jarvis/domain"
	pb "github.com/kai5263499/diy-jarvis/generated"
	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type config struct {
	AudioFileToProcess    string `env:"FILE"`
	AudioProcessorAddress string `env:"AUDIO_PROCESSOR_ADDRESS"`
}

var (
	numRequests  uint32
	numResponses uint32
	wavFilesRead uint32
	cfg          config
)

type processAudioDataRequest struct {
	uuid        *uuid.UUID
	wavWriter   *wav.Writer
	tmpFileName string
	msg         *pb.ProcessAudioRequest
}

func newProcessAudioRequest(meta *wav.File) *processAudioDataRequest {
	newUUID, _ := uuid.NewV4()

	f, err := ioutil.TempFile("", newUUID.String())
	domain.CheckError(err)

	writer, _ := meta.NewWriter(f)

	msg := &pb.ProcessAudioRequest{
		RequestId:      newUUID.String(),
		AudioStartTime: uint64(time.Now().Unix()),
	}

	return &processAudioDataRequest{
		uuid:        newUUID,
		wavWriter:   writer,
		tmpFileName: f.Name(),
		msg:         msg,
	}
}

func processChunk(stream pb.AudioProcessor_SubscribeClient, req *processAudioDataRequest) {
	var err error
	err = req.wavWriter.Close()
	domain.CheckError(err)

	content, err := ioutil.ReadFile(req.tmpFileName)
	domain.CheckError(err)

	req.msg.AudioData = content

	err = stream.Send(req.msg)
	atomic.AddUint32(&numRequests, 1)
	domain.CheckError(err)

	os.Remove(req.tmpFileName)
}

func processWavFile(stream pb.AudioProcessor_SubscribeClient, wavFile string, wg *sync.WaitGroup) error {
	defer wg.Done()
	var err error

	fileStat, err := os.Stat(wavFile)
	domain.CheckError(err)

	f, err := os.Open(wavFile)
	domain.CheckError(err)

	r, err := wav.NewReader(f, fileStat.Size())
	domain.CheckError(err)

	samplesPerChunk := r.GetSampleRate() * uint32(10)

	meta := wav.File{
		Channels:        r.GetNumChannels(),
		SampleRate:      r.GetSampleRate(),
		SignificantBits: r.GetBitsPerSample(),
	}

	req := newProcessAudioRequest(&meta)

	sampleCnt := uint32(0)
	chunksCnt := 0

	keepReading := true
	for keepReading {
		data, err := r.ReadRawSample()
		if err == io.EOF {
			break
		} else if err != nil {
			domain.CheckError(err)
		}

		sampleCnt++

		req.wavWriter.WriteSample(data)

		if sampleCnt%samplesPerChunk == 0 {
			processChunk(stream, req)

			chunksCnt++
			req = newProcessAudioRequest(&meta)
		}
	}

	atomic.AddUint32(&wavFilesRead, 1)

	processChunk(stream, req)

	return nil
}

func processResponses(stream pb.AudioProcessor_SubscribeClient, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		resp, err := stream.Recv()
		domain.CheckError(err)

		fmt.Printf("%s\n", resp.Output)

		atomic.AddUint32(&numResponses, 1)

		if atomic.LoadUint32(&wavFilesRead) > 0 && (atomic.LoadUint32(&numRequests) == atomic.LoadUint32(&numResponses)) {
			return
		}
	}
}

func main() {
	cfg = config{}
	err := env.Parse(&cfg)
	domain.CheckError(err)

	conn, err := grpc.Dial(cfg.AudioProcessorAddress, grpc.WithInsecure())
	domain.CheckError(err)
	defer conn.Close()

	client := pb.NewAudioProcessorClient(conn)
	stream, err := client.Subscribe(context.Background())
	domain.CheckError(err)
	defer stream.CloseSend()

	var wg sync.WaitGroup
	wg.Add(2)
	go processWavFile(stream, cfg.AudioFileToProcess, &wg)
	go processResponses(stream, &wg)
	wg.Wait()
}
