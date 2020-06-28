package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote" // adjusts to remote configuration center

)

type Config struct {
	Local  *LocalConfig
	Remote *RemoteConfig
}

type LocalConfig struct {
	Name string
	Path []string
}

type RemoteConfig struct {
	Provider string
	EndPoint string
	Path     string
	Name     string
	Type     string
}

const (
	DefaultName = "config"
	DefaultType = "json"
)

var (
	conf *Config
)

// New set config file search path
func Init(c *Config) error {
	if c == nil {
		return errors.New("config is nil")
	}

	conf = c

	viper.AutomaticEnv()

	if c.Local != nil {
		return c.Local.init()
	}

	if c.Remote != nil {
		return c.Remote.init()
	}

	return LocalConfig{}.init()
}

func (c RemoteConfig) init() error {
	fmt.Println("init remote config")
	if err := viper.AddRemoteProvider(c.Provider, c.EndPoint, c.Path); err != nil {
		return err
	}

	name := DefaultName
	if c.Name != "" {
		name = c.Name
	}

	t := DefaultType
	if c.Type != "" {
		t = c.Type
	}

	viper.SetConfigName(name)
	viper.SetConfigType(t)

	if err := viper.ReadRemoteConfig(); err != nil {
		return err
	}
	return nil
}

// 加载本地配置, 包含环境变量及本地文件.
func (c LocalConfig) init() (err error) {
	fmt.Println("init local config")
	if c.Path != nil {
		for _, p := range c.Path {
			viper.AddConfigPath(p)
		}
	}

	// 设置家目录为最后的查询路径
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	name := DefaultName
	if c.Name != "" {
		name = c.Name
	}
	viper.SetConfigName(name)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return
}
