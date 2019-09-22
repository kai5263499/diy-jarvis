package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/asticode/go-astideepspeech"

	"google.golang.org/grpc"

	"github.com/caarlos0/env"
	"github.com/kai5263499/diy-jarvis/domain"
	pb "github.com/kai5263499/diy-jarvis/generated"
	dsap "github.com/kai5263499/diy-jarvis/service/audio_processor/deepspeech"
)

type config struct {
	Alphabet             string  `env:"ALPHABET" envDefault:"/deepspeech_models/alphabet.txt"`
	LanguageModel        string  `env:"LM" envDefault:"/deepspeech_models/lm.binary"`
	Model                string  `env:"MODEL" envDefault:"/deepspeech_models/output_graph.pbmm"`
	Trie                 string  `env:"TRIE" envDefault:"/deepspeech_models/trie"`
	ListenPort           int     `env:"LISTEN_PORT" envDefault:"6000"`
	BeamWidth            int     `env:"BEAM_WIDTH" envDefault:"500"`
	NCep                 int     `env:"NCEP" envDefault:"26"`
	NContext             int     `env:"NCONTEXT" envDefault:"9"`
	LMWeight             float64 `env:"LM_WEIGHT" envDefault:"0.75"`
	ValidWordCountWeight float64 `env:"VALID_WORDCOUNT_WEIGHT" envDefault:"1.85"`
	UseTextProcessor     bool    `env:"USE_TEXT_PROCESSOR" envDefault:"true"`
	TextProcessorAddress string  `env:"TEXT_PROCESSOR_ADDRESS" envDefault:"localhost:6001"`
}

var (
	cfg                      config
	textEventProcessorClient pb.EventResponderClient
	textEventClientStream    pb.EventResponder_SubscribeClient
)

func main() {
	var err error
	cfg = config{}
	err = env.Parse(&cfg)
	domain.CheckError(err)

	fmt.Printf("Initialize DeepSpeech..\n")

	m := astideepspeech.New(cfg.Model, cfg.NCep, cfg.NContext, cfg.Alphabet, cfg.BeamWidth)
	defer m.Close()
	if cfg.LanguageModel != "" {
		m.EnableDecoderWithLM(cfg.Alphabet, cfg.LanguageModel, cfg.Trie, cfg.LMWeight, cfg.ValidWordCountWeight)
	}

	if cfg.UseTextProcessor {
		fmt.Printf("Connecting to text processor..")

		conn, err := grpc.Dial(cfg.TextProcessorAddress, grpc.WithInsecure())
		domain.CheckError(err)
		defer conn.Close()

		textEventProcessorClient = pb.NewEventResponderClient(conn)
		textEventClientStream, err = textEventProcessorClient.Subscribe(context.Background())
		domain.CheckError(err)
		defer textEventClientStream.CloseSend()
	}

	fmt.Printf("Starting server\n")

	ap := dsap.New(m, cfg.UseTextProcessor, &textEventClientStream)
	grpcServer := grpc.NewServer()
	pb.RegisterAudioProcessorServer(grpcServer, ap)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ListenPort))
	domain.CheckError(err)

	log.Printf("Listening on tcp://0.0.0.0:%d\n", cfg.ListenPort)
	grpcServer.Serve(l)
}
