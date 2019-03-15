package g2i

import (
	"context"
	"time"

	"encoding/json"
	"net/http"

	"io/ioutil"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

type GithubClient struct {
	c     *github.Client
	Owner string
	Repo  string
	token string
	ctx   context.Context
}

func NewGithubClient(config *Config, ctx context.Context) *GithubClient {
	if config.Github.Token != "" {

	}
	client := github.NewClient(nil)
	return &GithubClient{client, config.Github.RepoOwner, config.Github.Repo, config.Github.Token, ctx}
}

// GetIssues returns all repo issues
func (g *GithubClient) GetIssues() ([]*github.Issue, error) {
	log.Print("Get repository issues")
	issues, _, err := g.c.Issues.ListByRepo(g.ctx, g.Owner, g.Repo, nil)
	return issues, err
}

// GetRepoEvents
func (g *GithubClient) EventsURL() (string, error) {
	var eventsURL string
	repos, _, err := g.c.Repositories.ListByOrg(g.ctx, g.Owner, nil)
	if err != nil {
		return eventsURL, err
	}
	for i, _ := range repos {
		if *repos[i].Name == g.Repo {
			eventsURL = *repos[i].EventsURL
			break
		}
	}
	return eventsURL, nil
}

// GetEvents returns repo events
func (g *GithubClient) GetEvents(eventsURL string) ([]github.Event, error) {
	var events []github.Event

	log.Info("Get repository events")

	client := &http.Client{Timeout: time.Second * 15}
	req, err := http.NewRequest("GET", eventsURL, nil)
	if err != nil {
		log.Error(err)
		return events, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return events, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return events, err
	}

	log.Info("Data: %s", string(data))

	err = json.Unmarshal(data, &events)
	if err != nil {
		return events, err
	}

	log.Info("Events received")
	return events, nil
}

func (c *Config) IsEventProcessable(eventType string) bool {
	for i, _ := range c.Github.WatchedEventTypes {
		if eventType == c.Github.WatchedEventTypes[i] {
			log.Info("Processing event")
			return true
		}
	}
	return false
}
