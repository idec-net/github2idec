package g2i

import (
	"time"
)

type GithubIssueComment struct {
	Action  string  `json:"action"`
	Issue   Issue   `json:"issue"`
	Comment Comment `json:"comment"`
}

type User struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Labels struct {
	ID      int    `json:"id"`
	NodeID  string `json:"node_id"`
	URL     string `json:"url"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Default bool   `json:"default"`
}

type PullRequest struct {
	URL      string `json:"url"`
	HTMLURL  string `json:"html_url"`
	DiffURL  string `json:"diff_url"`
	PatchURL string `json:"patch_url"`
}

type Issue struct {
	URL               string        `json:"url"`
	RepositoryURL     string        `json:"repository_url"`
	LabelsURL         string        `json:"labels_url"`
	CommentsURL       string        `json:"comments_url"`
	EventsURL         string        `json:"events_url"`
	HTMLURL           string        `json:"html_url"`
	ID                int           `json:"id"`
	NodeID            string        `json:"node_id"`
	Number            int           `json:"number"`
	Title             string        `json:"title"`
	User              User          `json:"user"`
	Labels            []Labels      `json:"labels"`
	State             string        `json:"state"`
	Locked            bool          `json:"locked"`
	Assignee          interface{}   `json:"assignee"`
	Assignees         []interface{} `json:"assignees"`
	Milestone         interface{}   `json:"milestone"`
	Comments          int           `json:"comments"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	ClosedAt          interface{}   `json:"closed_at"`
	AuthorAssociation string        `json:"author_association"`
	PullRequest       PullRequest   `json:"pull_request"`
	Body              string        `json:"body"`
}

type Comment struct {
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	IssueURL          string    `json:"issue_url"`
	ID                int       `json:"id"`
	NodeID            string    `json:"node_id"`
	User              User      `json:"user"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	AuthorAssociation string    `json:"author_association"`
	Body              string    `json:"body"`
}
