package g2i

import (
	"log"
	"strings"

	"bytes"

	idec "github.com/Difrex/go-idec"
	"github.com/google/go-github/github"
)

type IDECClient struct {
	FetchConfig *idec.FetchConfig
	config      *Config
}

func NewIDECClient(config *Config) *IDECClient {
	return &IDECClient{
		config: config,
		FetchConfig: &idec.FetchConfig{
			Limit:  config.IDEC.Fetch.Limit,
			Offset: config.IDEC.Fetch.Offset,
			Node:   config.IDEC.NodeURL,
			Echoes: config.IDEC.Fetch.Echoes,
		},
	}
}

type HelloMessage struct {
	Sources, Name, Owner, Repo string
	Issues                     []*github.Issue
}

func (i *IDECClient) PostHello() {
	pointURL := strings.TrimRight(i.FetchConfig.Node, "/") + "/u/point"
	message := idec.PointMessage{}

	body := &bytes.Buffer{}
	err := i.config.generateTemplate("hello_message.tpl", i.config.IDEC.HelloMessageTemplatePath, i.FetchConfig, body)
	if err != nil {
		log.Print("Error", err.Error())
	}

	message.Body = body.String()
	log.Print(message, pointURL)
}
