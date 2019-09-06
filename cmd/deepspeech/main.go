package main

import (
	"flag"
	"fmt"
	"net"
	"log"
	"os"

	"github.com/asticode/go-astideepspeech"

	"google.golang.org/grpc"

	dsap "github.com/kai5263499/diy-jarvis/service/audio_processor/deepspeech"
	pb "github.com/kai5263499/diy-jarvis/generated"
)

// Constants
const (
	beamWidth            = 500
	nCep                 = 26
	nContext             = 9
	lmWeight             = 0.75
	validWordCountWeight = 1.85
)

var model = flag.String("model", "", "Path to the model (protocol buffer binary file)")
var alphabet = flag.String("alphabet", "", "Path to the configuration file specifying the alphabet used by the network")
var lm = flag.String("lm", "", "Path to the language model binary file")
var trie = flag.String("trie", "", "Path to the language model trie file created with native_client/generate_trie")
var version = flag.Bool("version", false, "Print version and exits")
var listenPort = flag.Int("port", 6000, "Port to listen for requests on")

func main() {
	flag.Parse()

	if *version {
		astideepspeech.PrintVersions()
		return
	}

	if *model == "" || *alphabet == "" {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	fmt.Printf("Initialize DeepSpeech..")

	m := astideepspeech.New(*model, nCep, nContext, *alphabet, beamWidth)
	defer m.Close()
	if *lm != "" {
		m.EnableDecoderWithLM(*alphabet, *lm, *trie, lmWeight, validWordCountWeight)
	}

	fmt.Printf("Done!\n")

	ap := dsap.New(m)
	grpcServer := grpc.NewServer()
	pb.RegisterAudioProcessorServer(grpcServer, ap)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *listenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening on tcp://localhost:%d", *listenPort)
	grpcServer.Serve(l)
}
