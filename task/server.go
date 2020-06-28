package task

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

const (
	DefaultQueue = "default_task"
)

var (
	Server     *machinery.Server
	serverPool = make(map[string]*machinery.Server)
)

func InitServer(conf *config.Config) error {
	var err error
	Server, err = machinery.NewServer(conf)
	if err != nil {
		return err
	}

	serverPool[DefaultQueue] = Server

	return nil
}

func GetServer(queue string) (*machinery.Server, error) {
	if exists, ok := serverPool[queue]; ok {
		return exists, nil
	}

	conf := *Server.GetConfig()
	conf.DefaultQueue = queue

	s, err := machinery.NewServer(&conf)
	if err != nil {
		return nil, err
	}
	serverPool[queue] = s
	return s, err
}
