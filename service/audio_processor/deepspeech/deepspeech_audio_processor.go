package deepspeech_audio_processor

import (
	"fmt"
	"io"
	"github.com/pkg/errors"
	"github.com/mattetti/filebuffer"
	"github.com/cryptix/wav"
	ds "github.com/asticode/go-astideepspeech"
	"github.com/asticode/go-astilog"
	pb "github.com/kai5263499/diy-jarvis/generated"
)

func New(m *ds.Model) *DeepSpeechAudioProcessor {
	return &DeepSpeechAudioProcessor{
		model: m,
	}
}

type DeepSpeechAudioProcessor struct {
	model *ds.Model
}

func metadataToString(m *ds.Metadata) string {
	retval := ""
	for _, item := range m.Items() {
		retval += item.Character()
	}
	return retval
}

func (p *DeepSpeechAudioProcessor) Subscribe(stream pb.AudioProcessor_SubscribeServer) error {
	fmt.Printf("new subscription!\n")
	for {
		req, err := stream.Recv()
		if err != nil {
			fmt.Printf("error recv request %+#v\n", err)
			break
		}

		if len(req.AudioData) < 1 {
			continue
		}
	
		f := filebuffer.New(req.AudioData)

		r, err := wav.NewReader(f, int64(len(req.AudioData)))
		if err != nil {
			astilog.Fatal(errors.Wrap(err, "creating new reader failed"))
			continue
		}

		// Read
		var d []int16
		for {
			// Read sample
			s, err := r.ReadSample()
			if err == io.EOF {
				break
			} else if err != nil {
				astilog.Fatal(errors.Wrap(err, "reading sample failed"))
			}

			// Append
			d = append(d, int16(s))
		}

		output := p.model.SpeechToText(d, uint(len(d)), 16000)

		resp := &pb.ProcessAudioResponse{
			RequestId: req.RequestId,
			AudioStartTime: req.AudioStartTime,
			ResponseCode: pb.ProcessAudioResponse_ACCEPTED,
			Output: output,
			
		}

		stream.Send(resp)
	}

	return nil
}