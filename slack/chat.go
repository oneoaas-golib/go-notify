package slack

import (
	"github.com/appscode/go-notify"
	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

const Uid = "slack"

type Options struct {
	AuthToken string   `envconfig:"AUTH_TOKEN" required:"true"`
	Channel   []string `envconfig:"CHANNEL" required:"true"`
}

type client struct {
	opt     Options
	channel []string
	body    string
}

var _ notify.ByChat = &client{}

func New(opt Options) *client {
	return &client{
		opt:     opt,
		channel: opt.Channel,
	}
}

func Default() (*client, error) {
	var opt Options
	err := envconfig.Process(Uid, &opt)
	if err != nil {
		return nil, err
	}
	return New(opt), nil
}

func (c *client) WithBody(body string) notify.ByChat {
	c.body = body
	return c
}

func (c *client) To(to string, cc ...string) notify.ByChat {
	c.channel = append([]string{to}, cc...)
	return c
}

func (c *client) Send() error {
	s := slack.New(c.opt.AuthToken)
	for _, channel := range c.channel {
		if _, _, err := s.PostMessage(channel, c.body, slack.PostMessageParameters{}); err != nil {
			return err
		}
	}
	return nil
}
