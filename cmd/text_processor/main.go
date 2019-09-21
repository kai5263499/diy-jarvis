package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"

	"gopkg.in/yaml.v2"

	"github.com/caarlos0/env"
	"github.com/kai5263499/diy-jarvis/domain"
	pb "github.com/kai5263499/diy-jarvis/generated"
	ep "github.com/kai5263499/diy-jarvis/service/event_processor/text"
)

type config struct {
	ListenPort      int    `env:"LISTEN_PORT" envDefault:"6001"`
	CommandSpecYaml string `env:"COMMAND_SPEC_YAML" envDefault:"commands.yaml"`
	Keyword         string `env:"KEYWORD", envDefault:"Jarvis"`
}

var (
	cfg config
)

func main() {
	var err error
	cfg = config{}
	err = env.Parse(&cfg)
	domain.CheckError(err)

	fmt.Printf("Initialize Text Event Processor..\n")

	var yamlCommands domain.TextEventCommands
	yamlData, err := ioutil.ReadFile(cfg.CommandSpecYaml)
	domain.CheckError(err)

	err = yaml.Unmarshal([]byte(yamlData), &yamlCommands)
	domain.CheckError(err)

	commands := make(map[string]domain.TextEventCommand)
	for _, c := range yamlCommands.Commands {
		commands[c.Command] = c
	}

	tp := ep.New(cfg.Keyword, commands)
	grpcServer := grpc.NewServer()
	pb.RegisterEventResponderServer(grpcServer, tp)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ListenPort))
	domain.CheckError(err)

	log.Printf("Listening on tcp://0.0.0.0:%d\n", cfg.ListenPort)
	grpcServer.Serve(l)
}
