package g2i

import (
	"github.com/boltdb/bolt"
)

type Config struct {
	Data   Data   `json:"data"`
	IDEC   IDEC   `json:"idec"`
	Github Github `json:"github"`
}

type Data struct {
	Path string `json:"path"`
	db   *bolt.DB
}

type Fetch struct {
	Echoes []string `json:"echoes"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

type IDEC struct {
	Fetch                    Fetch  `json:"fetch"`
	Authstring               string `json:"authstring"`
	NodeURL                  string `json:"node_url"`
	HelloMessage             bool   `json:"hello_message"`
	HelloMessageTemplatePath string `json:"hello_message_template_path"`
	TopPostID                string `json:"top_post_id"`
}

type Github struct {
	RepoOwner         string   `json:"repo_owner"`
	Repo              string   `json:"repo"`
	Token             string   `json:"token"`
	WatchedEventTypes []string `json:"watched_event_types"`
}
