package g2i

import (
	"context"
)

type Client struct {
	GHClient   *GithubClient
	IDECClient *IDECClient
}

func (c *Config) NewClient(ctx context.Context) (*Client, error) {
	client := &Client{}
	ghc := NewGithubClient(c, ctx)
	return client, nil
}
