package redis

import (
	"errors"
	"github.com/go-redis/redis"
)

type Config struct {
	Name    string
	Options *redis.Options
}

const (
	DefaultAddr = "127.0.0.1"
	DefaultName = "default"
)

var (
	Client     *redis.Client
	sourcePool = make(map[string]*redis.Client)
)

// Init redis connection
func Init(c *Config) error {
	err := checkConfig(c)
	if err != nil {
		return err
	}

	sourcePool[c.Name] = redis.NewClient(c.Options)
	if _, err = sourcePool[c.Name].Ping().Result(); err != nil {
		return err
	}

	if c.Name == DefaultName {
		Client = sourcePool[c.Name]
	}

	return nil
}

func checkConfig(c *Config) error {
	if c == nil {
		return errors.New("config is nil")
	}

	if c.Name == "" {
		c.Name = DefaultName
	}

	return nil
}

func Get(name string) *redis.Client {
	return sourcePool[name]
}
