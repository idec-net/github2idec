package g2i

import (
	"context"

	"time"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

const (
	EVENTS_UPDATER_SLEEP = 60
	ISSUES_UPDATER_SLEEP = 60
)

type Client struct {
	GHClient   *GithubClient
	IDECClient *IDECClient
	config     *Config
}

func (c *Config) NewClient(ctx context.Context) *Client {
	client := &Client{}
	ghc := NewGithubClient(c, ctx)
	ic := NewIDECClient(c)
	client.GHClient = ghc
	client.IDECClient = ic
	client.config = c
	return client
}

func (c *Client) Run() {
	// open db first
	err := c.config.openDB()
	if err != nil {
		panic(err)
	}
	defer c.config.Data.db.Close()

	evetsCH := make(chan github.Event)
	// Run issues updater
	go c.issuesUpdater()
	// Run events updater
	go c.eventsUpdater(evetsCH)
	// Run IDEC messager
	go c.idecEventMessager(evetsCH)

	// Main loop
	for {
	}
}

func (c *Client) idecEventMessager(ch chan github.Event) {
	for {
		event := <-ch
		log.Info(event)
	}
}

func (c *Client) eventsUpdater(ch chan github.Event) {
	log.Info("Events updater is running")
	eventsURL, err := c.GHClient.EventsURL()
	if err != nil {
		log.Error(err)
	}
	for {
		events, err := c.GHClient.GetEvents(eventsURL)
		if err != nil {
			log.Error(err)
			time.Sleep(time.Second * EVENTS_UPDATER_SLEEP)
		}

		log.Info(events)

		var newEvents []github.Event
		for i, _ := range events {
			if c.config.IsEventProcessable(*events[i].Type) {
				newEvents = append(newEvents, events[i])
			}
		}

		log.Info(newEvents)
		prevEvents, err := c.config.getEvents()
		for i, _ := range newEvents {
			if err == nil {
				if idNotInEvents(newEvents[i].GetID(), prevEvents) {
					ch <- newEvents[i]
				}
			} else {
				ch <- newEvents[i]
			}
		}
		if err := c.config.storeEvents(newEvents); err != nil {
			log.Error(err)
		}
		time.Sleep(time.Second * 60 * 60)
	}
}

func idNotInEvents(id string, events []github.Event) bool {
	for i, _ := range events {
		if events[i].GetID() == id {
			return false
		}
	}
	return true
}

func (c *Client) issuesUpdater() {
	for {
		issues, err := c.GHClient.GetIssues()
		if err != nil {
			log.Error(err)
		} else {
			if err := c.config.storeIssues(issues); err != nil {
				log.Error(err)
			}
		}
		time.Sleep(time.Second * ISSUES_UPDATER_SLEEP)
	}
}
