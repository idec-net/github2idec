package g2i

import (
	"encoding/json"
	"fmt"
	"os"

	"strconv"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
)

const (
	ISSUES_KEY    = "issues"
	ISSUES_BUCKET = "issues_bucket"
	EVENTS_KEY    = "events"
	EVENTS_BUCKET = "events_bucket"
)

func (c *Config) storeEvents(events []github.Event) error {
	if err := c.checkDB(); err != nil {
		return err
	}

	data, err := json.Marshal(events)
	if err != nil {
		return err
	}
	err = c.Data.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(EVENTS_BUCKET))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(EVENTS_KEY), data)
		if err != nil {
			return err
		}

		// Also store an events list size
		err = bucket.Put([]byte(EVENTS_KEY+"_count"), []byte(string(len(events))))
		return err
	})

	return err
}

func (c *Config) getEvents() ([]github.Event, error) {
	var events []github.Event

	if err := c.checkDB(); err != nil {
		return nil, err
	}

	err := c.Data.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(EVENTS_BUCKET))
		if bucket == nil {
			return fmt.Errorf("Events bucket not created yet!")
		}

		data := bucket.Get([]byte(EVENTS_KEY))
		err := json.Unmarshal(data, &events)

		return err
	})
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (c *Config) getEventsCount() int {
	if err := c.checkDB(); err != nil {
		log.Error(err)
		return -1
	}

	var count int
	err := c.Data.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(EVENTS_BUCKET))
		if bucket == nil {
			log.Error("Events bucket not created yet!")
			return nil
		}

		data := bucket.Get([]byte(EVENTS_KEY + "_count"))
		if data != nil {
			c, err := strconv.Atoi(string(data))
			if err != nil {
				return err
			}
			count = c
		}

		return nil
	})
	if err != nil {
		return -1
	}

	return count
}

func (c *Config) storeIssues(issues []*github.Issue) error {
	if err := c.checkDB(); err != nil {
		return err
	}

	data, err := json.Marshal(issues)
	if err != nil {
		return err
	}
	err = c.Data.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(ISSUES_BUCKET))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(ISSUES_KEY), data)
		return err
	})

	return err
}

func (c *Config) getIssues() ([]*github.Issue, error) {
	var issues []*github.Issue

	if err := c.checkDB(); err != nil {
		return nil, err
	}

	err := c.Data.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ISSUES_BUCKET))
		if bucket == nil {
			return fmt.Errorf("Issues bucket not created yet!")
		}

		data := bucket.Get([]byte(ISSUES_KEY))
		err := json.Unmarshal(data, &issues)

		return err
	})
	if err != nil {
		return nil, err
	}

	return issues, nil
}

func (c *Config) checkDB() error {
	if c.Data.db == nil {
		err := c.openDB()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) openDB() error {
	err := c.createDB()
	if err != nil {
		return err
	}

	db, err := bolt.Open(c.Data.Path+"/db.bolt", 0600, nil)
	if err != nil {
		return err
	}
	c.Data.db = db

	log.Info("Database is open")

	return nil
}

func (c *Config) createDB() error {
	c.createDataDir()

	db, err := bolt.Open(c.Data.Path+"/db.bolt", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}

func (c *Config) createDataDir() {
	_, err := os.Open(c.Data.Path)
	if !os.IsExist(err) {
		os.Mkdir(c.Data.Path, 0700)
	}
}
