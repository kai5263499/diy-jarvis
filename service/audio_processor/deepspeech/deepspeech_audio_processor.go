package deepspeech_audio_processor

import (
	"fmt"
	"io"

	ds "github.com/asticode/go-astideepspeech"
	"github.com/cryptix/wav"
	pb "github.com/kai5263499/diy-jarvis/generated"
	"github.com/mattetti/filebuffer"
)

func New(m *ds.Model, useTextProcessor bool, s *pb.EventResponder_SubscribeClient) *DeepSpeechAudioProcessor {
	return &DeepSpeechAudioProcessor{
		model:                 m,
		useTextProcessor:      useTextProcessor,
		textEventClientStream: s,
	}
}

type DeepSpeechAudioProcessor struct {
	model                 *ds.Model
	useTextProcessor      bool
	textEventClientStream *pb.EventResponder_SubscribeClient
}

func metadataToString(m *ds.Metadata) string {
	retval := ""
	for _, item := range m.Items() {
		retval += item.Character()
	}
	return retval
}

func (p *DeepSpeechAudioProcessor) sendTextEventRequest(input string, req *pb.ProcessAudioRequest) {
	textEventReq := &pb.TextEventRequest{
		RequestId: req.RequestId,
		SourceId:  req.SourceId,
		Text:      input,
	}

	fmt.Printf("sending text event request %+#v\n", textEventReq)

	(*p.textEventClientStream).Send(textEventReq)
}

func (p *DeepSpeechAudioProcessor) Subscribe(stream pb.AudioProcessor_SubscribeServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			fmt.Printf("recv error=%+#v\n", err)
			break
		}

		if len(req.AudioData) < 1 {
			continue
		}

		f := filebuffer.New(req.AudioData)

		r, err := wav.NewReader(f, int64(len(req.AudioData)))
		if err != nil {
			continue
		}

		var d []int16
		for {
			s, err := r.ReadSample()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("sample read error=%+#v\n", err)
			}
			d = append(d, int16(s))
		}

		output := p.model.SpeechToText(d, uint(len(d)), 16000)

		if p.useTextProcessor && len(output) > 0 {
			p.sendTextEventRequest(output, req)
		}

		resp := &pb.ProcessAudioResponse{
			RequestId:      req.RequestId,
			AudioStartTime: req.AudioStartTime,
			ResponseCode:   pb.ProcessAudioResponse_ACCEPTED,
			Output:         output,
		}

		stream.Send(resp)
	}

	return nil
}
