package main

import (
	"flag"
	"io"
	"os"
	"fmt"
	"time"
	"io/ioutil"
	"sync"
	"sync/atomic"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/cryptix/wav"
	pb "github.com/kai5263499/diy-jarvis/generated"
	uuid "github.com/nu7hatch/gouuid"
)

var (
	audioFileToProcess = flag.String("in", "", "Path to the audio file to run (WAV format)")
	audioProcessorAddress = flag.String("server", "localhost:6000", "Address of the audio processor to send samples to")

	numRequests uint32
	numResponses uint32
	wavFilesRead uint32
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

func newProcessAudioRequest(meta *wav.File) *processAudioDataRequest {
	newUUID, _ := uuid.NewV4()

	f, err := ioutil.TempFile("", newUUID.String())
	checkError(err)

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

func processChunk(stream pb.AudioProcessor_SubscribeClient, req *processAudioDataRequest) {
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

func processWavFile(stream pb.AudioProcessor_SubscribeClient, wavFile string, wg *sync.WaitGroup) error {
	defer wg.Done()
	var err error

	fileStat, err := os.Stat(wavFile)
	checkError(err)

	f, err := os.Open(wavFile)
	checkError(err)

	r, err := wav.NewReader(f, fileStat.Size())
	checkError(err)
	
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
			checkError(err)
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
		checkError(err)

		fmt.Printf("%s\n", resp.Output)
		
		atomic.AddUint32(&numResponses, 1)

		if atomic.LoadUint32(&wavFilesRead) > 0 && (atomic.LoadUint32(&numRequests) == atomic.LoadUint32(&numResponses)) {
			return
		}
	}
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*audioProcessorAddress, grpc.WithInsecure())
	checkError(err)
	defer conn.Close()

	client := pb.NewAudioProcessorClient(conn)
	stream, err := client.Subscribe(context.Background())
	checkError(err)
	defer stream.CloseSend()

	var wg sync.WaitGroup
	wg.Add(2)
	go processWavFile(stream, *audioFileToProcess, &wg)
	go processResponses(stream, &wg)
	wg.Wait()
}
