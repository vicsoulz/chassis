package mongo

import (
	"errors"
	"strings"

	"github.com/globalsign/mgo"
)

type Config struct {
	Uri           string
	DataBase      string
	PoolLimit     int
	MinPoolSize   int
	MaxIdleTimeMS int
}

const (
	DefaultPoolLimit     = 20
	DefaultDataBase      = "test"
	DefaultUri           = "mongodb://localhost:27017"
	DefaultMinPoolSize   = 20
	DefaultMaxIdleTimeMS = 2 * 60 * 1000
)

var (
	Conf    *Config
	Session *mgo.Session
	DB      *mgo.Database
)

func Init(c *Config) error {
	if c == nil {
		return errors.New("config is nil")
	}

	if c.Uri == "" {
		return errors.New("uri is nil")
	}
	initConfig(c)
	Conf = c

	info, err := mgo.ParseURL(c.Uri)
	if err != nil {
		return err
	}

	info.MinPoolSize = c.MinPoolSize
	info.MaxIdleTimeMS = c.MaxIdleTimeMS

	Session, err = mgo.DialWithInfo(info)
	if err != nil {
		return err
	}

	err = Session.Ping()
	if err != nil {
		return err
	}

	Session.SetPoolLimit(DefaultPoolLimit)
	Session.SetMode(mgo.Eventual, true)

	DB = Session.DB(c.DataBase)

	return nil
}

func InitDefault() error {
	return Init(DefaultConfig())
}

func DefaultConfig() *Config {
	return &Config{
		Uri:           DefaultUri,
		DataBase:      DefaultDataBase,
		PoolLimit:     DefaultPoolLimit,
		MinPoolSize:   DefaultMinPoolSize,
		MaxIdleTimeMS: DefaultMaxIdleTimeMS,
	}
}

func initConfig(c *Config) {
	if c.MinPoolSize == 0 {
		c.MinPoolSize = DefaultMinPoolSize
	}

	if c.MaxIdleTimeMS == 0 {
		c.MaxIdleTimeMS = DefaultMaxIdleTimeMS
	}

	if c.MinPoolSize == 0 {
		c.PoolLimit = DefaultPoolLimit
	}

	if c.DataBase == "" {
		c.DataBase = DefaultDataBase
	}
}

func C(name string) *mgo.Collection {
	return DB.C(name)
}

func NotFound(err error) bool {
	if err == nil {
		return false
	}

	if err == mgo.ErrNotFound || strings.Contains(err.Error(), "not found") {
		return true
	}
	return false
}
