package g2i

import (
	"encoding/json"
	"text/template"

	"bytes"

	"github.com/google/go-github/github"
	idec "github.com/idec-net/go-idec"
	log "github.com/sirupsen/logrus"
)

type IDECClient struct {
	FetchConfig idec.FetchConfig
	authstring  string
	config      *Config
	TopPostID   string
}

func NewIDECClient(config *Config) *IDECClient {
	return &IDECClient{
		config:     config,
		authstring: config.IDEC.Authstring,
		FetchConfig: idec.FetchConfig{
			Limit:  config.IDEC.Fetch.Limit,
			Offset: config.IDEC.Fetch.Offset,
			Node:   config.IDEC.NodeURL,
			Echoes: config.IDEC.Fetch.Echoes,
		},
		TopPostID: config.IDEC.TopPostID,
	}
}

type HelloMessage struct {
	Sources, Name, Owner, Repo string
	Issues                     []*github.Issue
}

func (i *IDECClient) PostHello() error {
	message := idec.PointMessage{}
	message.Echo = i.FetchConfig.Echoes[0]
	message.To = "All"
	message.Repto = "@repto:" + i.config.IDEC.TopPostID
	message.Subg = i.config.IDEC.MessageSubg

	body := &bytes.Buffer{}
	err := i.config.generateTemplate("hello_message.tpl", i.config.IDEC.HelloMessageTemplatePath, i.config.Github, body)
	if err != nil {
		log.Print("Error", err.Error())
		return err
	}

	message.Body = body.String()

	b64message := message.PrepareMessageForSend()
	err = i.FetchConfig.PostMessage(i.authstring, b64message)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (i *IDECClient) PostComment(event github.Event) error {
	var post idec.PointMessage
	comment, err := i.generateCommentTemplate(event)
	if err != nil {
		return err
	}
	post.Echo = i.FetchConfig.Echoes[0]
	post.To = "All"
	post.Repto = "@repto:" + i.TopPostID
	post.Body = comment
	post.Subg = i.config.IDEC.MessageSubg

	b64message := post.PrepareMessageForSend()
	err = i.FetchConfig.PostMessage(i.authstring, b64message)

	return err
}

func (i *IDECClient) generateCommentTemplate(event github.Event) (string, error) {
	payload, err := event.RawPayload.MarshalJSON()
	if err != nil {
		return "", err
	}

	var ghic GithubIssueComment
	err = json.Unmarshal(payload, &ghic)
	if err != nil {
		return "", err
	}

	body := &bytes.Buffer{}
	t, err := template.New("comment.tpl").ParseFiles(i.config.IDEC.CommentTemplatePath)
	if err != nil {
		return "", err
	}
	err = t.Execute(body, ghic)
	if err != nil {
		return "", err
	}
	log.Debug(body.String())

	return body.String(), nil
}
