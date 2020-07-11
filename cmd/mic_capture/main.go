package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"io/ioutil"
	"os"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/caarlos0/env"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/cryptix/wav"
	"github.com/gofrs/uuid"
	dj "github.com/kai5263499/diy-jarvis"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"github.com/xlab/closer"
	"github.com/xlab/portaudio-go/portaudio"
)

const (
	sampleRate   = 16000
	channels     = 1
	sampleFormat = portaudio.PaInt16
)

type config struct {
	MQTTBroker           string        `env:"MQTT_BROKER"`
	MQTTClientID         string        `env:"MQTT_CLIENT_ID" envDefault:"miccapture"`
	LogLevel             string        `env:"LOG_LEVEL" envDefault:"info"`
	PulseInterval        time.Duration `env:"PULSE_DURATION" envDefault:"30s"`
	AudioCaptureDuration uint          `env:"DURATION" envDefault:"3"`
}

var (
	samplesPerChannel = uint(16000)

	numRequests  uint32
	numResponses uint32
	cfg          config
	mqttComms    *dj.MqttComms
	sourceID     string
)

type processAudioDataRequest struct {
	uuid        uuid.UUID
	wavWriter   *wav.Writer
	tmpFileName string
}

func newProcessAudioRequest() *processAudioDataRequest {
	newUUID := uuid.Must(uuid.NewV4())

	f, tmpFileErr := ioutil.TempFile("", newUUID.String())
	if tmpFileErr != nil {
		logrus.WithError(tmpFileErr).Errorf("creating temp file")
		return nil
	}

	meta := wav.File{
		Channels:        channels,
		SampleRate:      sampleRate,
		SignificantBits: 16,
	}
	writer, newWriterErr := meta.NewWriter(f)
	if newWriterErr != nil {
		logrus.WithError(newWriterErr).Errorf("temp wav writer")
		return nil
	}

	return &processAudioDataRequest{
		uuid:        newUUID,
		wavWriter:   writer,
		tmpFileName: f.Name(),
	}
}

func processChunk(req *processAudioDataRequest) {
	if err := req.wavWriter.Close(); err != nil {
		logrus.WithError(err).Errorf("error closing wavWriter")
		return
	}

	content, readFileErr := ioutil.ReadFile(req.tmpFileName)
	if readFileErr != nil {
		logrus.WithError(readFileErr).Errorf("error reading tmpFileName %s", req.tmpFileName)
	}

	if token := mqttComms.MQTTClient().Publish(sourceID, 0, false, content); token.Error() != nil {
		logrus.WithError(token.Error()).Errorf("error sending audio data to the topic %s", sourceID)
		return
	}

	logrus.Debugf("sent %d of audio data to channel %s", len(content), sourceID)

	atomic.AddUint32(&numRequests, 1)

	if err := os.Remove(req.tmpFileName); err != nil {
		logrus.WithError(err).Errorf("error removing tmpFileName %s", req.tmpFileName)
	}

	logrus.Debugf("removed temp file")
}

func processMicInput() error {
	defer closer.Close()

	if paErr := portaudio.Initialize(); portaudio.ErrorCode(paErr) != portaudio.PaNoError {
		logrus.WithFields(logrus.Fields{
			"pa-error-code": portaudio.ErrorCode(paErr),
			"pa-error-text": portaudio.GetErrorText(paErr),
		}).Fatal("initializing port audio")
		return errors.New("unable to initialize mic input")
	}

	closer.Bind(func() {
		if paErr := portaudio.Terminate(); portaudio.ErrorCode(paErr) != portaudio.PaNoError {
			logrus.WithFields(logrus.Fields{
				"pa-error-code": portaudio.ErrorCode(paErr),
				"pa-error-text": portaudio.GetErrorText(paErr),
			}).Fatal("port audio terminate")
			return
		}
	})

	var paStream *portaudio.Stream
	if paErr := portaudio.OpenDefaultStream(&paStream, channels, 0, sampleFormat, sampleRate,
		samplesPerChannel, paCallback, nil); portaudio.ErrorCode(paErr) != portaudio.PaNoError {
		logrus.WithFields(logrus.Fields{
			"pa-error-code": portaudio.ErrorCode(paErr),
			"pa-error-text": portaudio.GetErrorText(paErr),
		}).Fatal("port audio open default stream")
		return errors.New("error opening default stream")
	}

	closer.Bind(func() {
		if paErr := portaudio.CloseStream(paStream); portaudio.ErrorCode(paErr) != portaudio.PaNoError {
			logrus.WithFields(logrus.Fields{
				"pa-error-code": portaudio.ErrorCode(paErr),
				"pa-error-text": portaudio.GetErrorText(paErr),
			}).Fatal("portaudio close stream")
			return
		}
	})

	if paErr := portaudio.StartStream(paStream); portaudio.ErrorCode(paErr) != portaudio.PaNoError {
		logrus.WithFields(logrus.Fields{
			"pa-error-code": portaudio.ErrorCode(paErr),
		}).Fatal("portaudio start stream")
		return errors.New("error starting portaudio stream")
	}

	closer.Bind(func() {
		if paErr := portaudio.StopStream(paStream); portaudio.ErrorCode(paErr) != portaudio.PaNoError {
			logrus.WithFields(logrus.Fields{
				"pa-error-code": portaudio.ErrorCode(paErr),
				"pa-error-text": portaudio.GetErrorText(paErr),
			}).Fatal("portaudio stop stream")
			return
		}
	})

	closer.Hold()

	return nil
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

	return statusContinue
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

	var newMqttErr error
	mqttComms, newMqttErr = dj.NewMqttComms(cfg.MQTTClientID, cfg.MQTTBroker)
	if newMqttErr != nil {
		logrus.WithError(newMqttErr).Fatal("error creating mqtt comms")
	}

	sourceID = uuid.Must(uuid.NewV4()).String()
	sendAudioRegistrationRequest()

	samplesPerChannel = sampleRate * cfg.AudioCaptureDuration

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		for {
			ticker := time.NewTicker(cfg.PulseInterval)
			for {
				select {
				case <-ticker.C:
					if err := sendAudioRegistrationRequest(); err != nil {
						return err
					}
				}
			}
		}
	})

	g.Go(func() error {
		for {
			if err := processMicInput(); err != nil {
				return err
			}
		}
	})

	g.Go(func() error {
		for {
			select {
			case _ = <-mqttComms.RequestChan():
				// react to some messages here
			}
		}
	})

	g.Wait()
}
