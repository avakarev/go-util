package natsutil

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/avakarev/go-util/envutil"
)

// Conn implements nats connection
type Conn struct {
	env           envutil.AppEnv
	conn          *nats.Conn
	subscriptions map[*nats.Subscription]struct{}
}

func (c *Conn) envSubj(subj string) string {
	if strings.HasPrefix(subj, envutil.EnvDev) ||
		strings.HasPrefix(subj, envutil.EnvBeta) ||
		strings.HasPrefix(subj, envutil.EnvProd) {
		return subj
	}
	return fmt.Sprintf("%s.%s", c.env.String(), subj)
}

// Publish sends byte slice to the given subject
func (c *Conn) Publish(subj string, data []byte) error {
	return c.conn.Publish(c.envSubj(subj), data)
}

// PublishJSON marshalls given pointer destination into JSON and sends to the given subject
func (c *Conn) PublishJSON(subj string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.conn.Publish(c.envSubj(subj), bytes)
}

// Subscribe subscribes given handler to the given subject
func (c *Conn) Subscribe(subj string, fn nats.MsgHandler) error {
	sub, err := c.conn.Subscribe(c.envSubj(subj), fn)
	if err != nil {
		return err
	}
	log.Debug().Str("subject", sub.Subject).Msg("nats: subscribed")
	c.subscriptions[sub] = struct{}{}
	return nil
}

// QueueSubscribe subscribes given handler to the given subject
func (c *Conn) QueueSubscribe(subj string, queue string, fn nats.MsgHandler) error {
	sub, err := c.conn.QueueSubscribe(c.envSubj(subj), queue, fn)
	if err != nil {
		return err
	}
	log.Debug().Str("subject", sub.Subject).Str("queue", queue).Msg("nats: subscribed")
	c.subscriptions[sub] = struct{}{}
	return nil
}

// Request sends request and returns reply's message
func (c *Conn) Request(subj string, v any, timeout time.Duration) (*nats.Msg, error) {
	if v != nil {
		dataBytes, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return c.conn.Request(c.envSubj(subj), dataBytes, timeout)
	}
	return c.conn.Request(c.envSubj(subj), nil, timeout)
}

// RequestBytes sends request and returns reply's bytes
func (c *Conn) RequestBytes(subj string, v any, timeout time.Duration) ([]byte, error) {
	resp, err := c.Request(subj, v, timeout)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// RequestJSON sends requests and unmarshals reply's json bytes into given destination
func (c *Conn) RequestJSON(subj string, v any, timeout time.Duration, destPtr any) error {
	bytes, err := c.RequestBytes(subj, v, timeout)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, destPtr)
}

// Close unsubscribes consumers and closes connections
func (c *Conn) Close() error {
	for sub := range c.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			return err
		}
		delete(c.subscriptions, sub)
	}
	c.conn.Close()
	return nil
}

// ConnConfig defines connection configuration
type ConnConfig struct {
	Env           envutil.AppEnv
	URL           string
	User          string
	Password      string
	Timeout       time.Duration
}

// NewConn returns new connection value
func NewConn(config *ConnConfig) (*Conn, error) {
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}
	conn, err := nats.Connect(
		config.URL,
		nats.UserInfo(config.User, config.Password),
		nats.Timeout(config.Timeout),
	)
	if err != nil {
		return nil, err
	}
	return &Conn{
		env:           config.Env,
		conn:          conn,
		subscriptions: make(map[*nats.Subscription]struct{}),
	}, nil
}

// DefaultConn returns new connection initialized from env variables
func DefaultConn() (*Conn, error) {
	env, err := envutil.NewAppEnv()
	if err != nil {
		return nil, err
	}
	url, err := envutil.MustStr("NATS_URL")
	if err != nil {
		return nil, err
	}
	user, err := envutil.MustStr("NATS_USER")
	if err != nil {
		return nil, err
	}
	password, err := envutil.MustStr("NATS_PASSWORD")
	if err != nil {
		return nil, err
	}
	return NewConn(&ConnConfig{
		Env:      env,
		URL:      url,
		User:     user,
		Password: password,
	})
}
