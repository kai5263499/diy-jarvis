package processor

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kai5263499/diy-jarvis/domain"
	pb "github.com/kai5263499/diy-jarvis/generated"
)

// New returns an instance of the Text struct
func New(keyword string, commands map[string]domain.TextEventCommand) *Text {
	return &Text{
		Keyword:  keyword,
		Commands: commands,
	}
}

// Text is an instance of a text event processor
type Text struct {
	Keyword  string
	Commands map[string]domain.TextEventCommand
}

func (p *Text) processText(command string) {
	if action, ok := p.Commands[command]; ok {
		fmt.Printf("got command %s, performing action %+#v\n", command, action)
	} else {
		// Ignore command
	}
}

// Subscribe opens a bidirectional text event request processing channel
func (p *Text) Subscribe(stream pb.EventResponder_SubscribeServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			var code codes.Code
			code = status.Convert(err).Code()
			if code == codes.Canceled {
				break
			}

			fmt.Printf("recv error=%+#v\n", err)
			break
		}

		// TODO test for keyword and copy only trailing text
		p.processText(req.Text)

		stream.Send(&pb.TextEventResponse{
			RequestId:    req.RequestId,
			ResponseCode: pb.TextEventResponse_ACCEPTED,
		})
	}

	return nil
}
