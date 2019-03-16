package g2i

import (
	"context"

	"time"

	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

const (
	EVENTS_UPDATER_SLEEP = 60 * 10
	ISSUES_UPDATER_SLEEP = 60 * 10
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
	// Run IDEC messager
	go c.idecEventMessenger(evetsCH)
	// Run events updater
	c.eventsUpdater(evetsCH)
}

func (c *Client) idecEventMessenger(ch chan github.Event) {
	if !c.config.isHelloSent() {
		c.IDECClient.PostHello()
		c.config.storeHello()
	}
	for {
		event := <-ch
		err := c.IDECClient.PostComment(event)
		if err != nil {
			log.Error(err)
		}
	}
}

func (c *Client) eventsUpdater(ch chan github.Event) {
	log.Info("Events updater is running")
	eventsURL, err := c.GHClient.EventsURL()
	if err != nil {
		log.Error(err)
	}
	for {
		log.Info("Get repository events")
		events, err := c.GHClient.GetEvents(eventsURL)
		if err != nil {
			log.Error(err)
			time.Sleep(time.Second * EVENTS_UPDATER_SLEEP)
		}

		var newEvents []github.Event
		for i, _ := range events {
			if c.config.IsEventProcessable(*events[i].Type) {
				newEvents = append(newEvents, events[i])
			}
		}

		prevEvents, err := c.config.getEvents()
		for i := len(newEvents) - 1; i >= 0; i-- {
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
	log.Info("Issues updater is running")
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
