package main

import (
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/cryptix/wav"
	"github.com/gofrs/uuid"
	dj "github.com/kai5263499/diy-jarvis"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type config struct {
	MQTTBroker         string        `env:"MQTT_BROKER"`
	MQTTClientID       string        `env:"MQTT_CLIENT_ID" envDefault:"wavslicer"`
	LogLevel           string        `env:"LOG_LEVEL" envDefault:"info"`
	AudioFileToProcess string        `env:"FILE"`
	AudioSampleSize    time.Duration `env:"AUDIO_SAMPLE_SIZE" envDefault:"3s"`
}

var (
	wavFilesRead uint32
	cfg          config
	mqttComms    *dj.MqttComms
	sourceID     string
)

type processAudioDataRequest struct {
	uuid        uuid.UUID
	wavWriter   *wav.Writer
	tmpFileName string
}

func newProcessAudioRequest(meta *wav.File) *processAudioDataRequest {
	newUUID := uuid.Must(uuid.NewV4())

	f, tmpFileErr := ioutil.TempFile("", newUUID.String())
	if tmpFileErr != nil {
		logrus.WithError(tmpFileErr).Fatal("error creating temporary file")
	}

	writer, _ := meta.NewWriter(f)

	return &processAudioDataRequest{
		uuid:        newUUID,
		wavWriter:   writer,
		tmpFileName: f.Name(),
	}
}

func processChunk(req *processAudioDataRequest) {
	if err := req.wavWriter.Close(); err != nil {
		logrus.WithError(err).Fatal("error closing wav writer")
	}

	content, readFileErr := ioutil.ReadFile(req.tmpFileName)
	if readFileErr != nil {
		logrus.WithError(readFileErr).Fatal("error reading tmp file")
	}

	if token := mqttComms.MQTTClient().Publish(sourceID, 0, false, content); token.Error() != nil {
		logrus.WithError(token.Error()).Errorf("error sending audio data to the topic %s", sourceID)
		return
	}

	logrus.Debugf("sent %d of audio data to channel %s", len(content), sourceID)

	if err := os.Remove(req.tmpFileName); err != nil {
		logrus.WithError(err).Errorf("error removing tmpFileName %s", req.tmpFileName)
	}

	logrus.Debug("removed temp file")
}

func processWavFile(wavFile string) error {
	fileStat, osStatErr := os.Stat(wavFile)
	if osStatErr != nil {
		logrus.WithError(osStatErr).Fatal("error stating wav file")
		return osStatErr
	}

	f, wavOpenErr := os.Open(wavFile)
	if wavOpenErr != nil {
		logrus.WithError(wavOpenErr).Fatal("error opening wav file")
		return wavOpenErr
	}

	r, newReaderErr := wav.NewReader(f, fileStat.Size())
	if newReaderErr != nil {
		logrus.WithError(newReaderErr).Fatal("error creating new wav reader")
		return newReaderErr
	}

	samplesPerChunk := r.GetSampleRate() * uint32(cfg.AudioSampleSize.Seconds())

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
		data, readRawSampleErr := r.ReadRawSample()
		if readRawSampleErr == io.EOF {
			break
		} else if readRawSampleErr != nil {
			logrus.WithError(readRawSampleErr).Fatal("error reading raw sample")
		}

		sampleCnt++

		req.wavWriter.WriteSample(data)

		if sampleCnt%samplesPerChunk == 0 {
			processChunk(req)

			chunksCnt++
			req = newProcessAudioRequest(&meta)
		}
	}

	processChunk(req)

	return nil
}

func sendAudioRegistrationRequest() error {
	msg := pb.Base{
		Id:        uuid.Must(uuid.NewV4()).String(),
		Timestamp: uint64(time.Now().Unix()),
		Type:      pb.Type_RegisterAudioSourceRequestType,
		SourceId:  sourceID,
	}

	if err := mqttComms.SendRequest(msg); err != nil {
		logrus.WithError(err).Error("error sending audio channel id pulse")
		return errors.Wrap(err, "error sending audio channel id pulse")
	}

	return nil
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

	var newMqttErr error
	mqttComms, newMqttErr = dj.NewMqttComms(cfg.MQTTClientID, cfg.MQTTBroker)
	if newMqttErr != nil {
		logrus.WithError(newMqttErr).Fatal("error creating mqtt comms")
	}

	sendAudioRegistrationRequest()
	processWavFile(cfg.AudioFileToProcess)
}
