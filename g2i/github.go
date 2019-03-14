package g2i

import (
	"context"
	"log"

	"encoding/json"
	"net/http"

	"github.com/google/go-github/github"
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

	client := &http.Client{}
	req, err := http.NewRequest("GET", eventsURL, nil)
	if err != nil {
		return events, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return events, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		return events, err
	}

	return events, nil
}

func (c *Config) IsEventProcessable(eventType string) bool {
	for i, _ := range c.Github.WatchedEventTypes {
		if eventType == c.Github.WatchedEventTypes[i] {
			return true
		}
	}
	return false
}
