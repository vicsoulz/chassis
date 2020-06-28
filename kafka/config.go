package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/spf13/viper"
	"time"
)

// 本地关于kafka的配置结构
type CustomConfig struct {
	Topics   []string `mapstructure:"topics"`
	Address  []string `mapstructure:"address"`
	UserName string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	Cert     string   `mapstructure:"cert_file_content"`
}

type ConsumerCustomConfig struct {
	CustomConfig
	ConsumerGroup string `mapstructure:"consumer_group"`
	Offset        int64 // 自定义Offset, 默认为OffsetOldest
}

var (
	minVersion = sarama.V0_9_0_0
)

func NewConfigFromViper() (*sarama.Config, error) {
	custom, err := GetCustomConfigFromViper()
	if err != nil {
		return nil, err
	}

	return NewConfig(custom)
}

func NewConfig(custom *CustomConfig) (*sarama.Config, error) {
	conf := sarama.NewConfig()

	conf.Net.SASL.Enable = false
	conf.Net.SASL.Handshake = true
	conf.Net.TLS.Enable = false //true 阿里商用版这个参数要为false
	conf.Producer.Return.Successes = true
	conf.Net.DialTimeout = 10 * time.Second
	conf.Net.ReadTimeout = 10 * time.Second
	conf.Net.WriteTimeout = 10 * time.Second

	if custom.UserName != "" {
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = custom.UserName
		conf.Net.SASL.Password = custom.Password
	}

	if custom.Cert != "" {
		clientCertPool := x509.NewCertPool()
		ok := clientCertPool.AppendCertsFromPEM([]byte(custom.Cert))
		if !ok {
			return nil, errors.New("kafka syncProducer failed to parse root certificate")
		}
		conf.Net.TLS.Config = &tls.Config{
			//Certificates:       []tls.Certificate{},
			RootCAs:            clientCertPool,
			InsecureSkipVerify: true,
		}
	}

	if err := conf.Validate(); err != nil {
		return nil, fmt.Errorf("Kafka syncProducer config invalidate. config: %v. err: %v", conf, err)
	}

	return conf, nil
}

func NewConsumerConfig(custom *ConsumerCustomConfig) (*cluster.Config, error) {
	config, err := NewConfig(&custom.CustomConfig)
	if err != nil {
		return nil, err
	}

	c := &cluster.Config{
		Config: *config,
	}
	c.Group.PartitionStrategy = cluster.StrategyRange
	c.Group.Offsets.Retry.Max = 3
	c.Group.Offsets.Synchronization.DwellTime = c.Consumer.MaxProcessingTime
	c.Group.Session.Timeout = 30 * time.Second
	c.Group.Heartbeat.Interval = 3 * time.Second
	c.Config.Version = minVersion

	if custom.Offset == 0 {
		c.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		c.Consumer.Offsets.Initial = custom.Offset
	}

	return c, nil
}

func GetCustomConfigFromViper() (custom *CustomConfig, err error) {
	err = viper.UnmarshalKey("kafka", &custom)
	return
}

func GetConsumerCustomConfigFromViper() (custom *ConsumerCustomConfig, err error) {
	err = viper.UnmarshalKey("kafka", &custom)
	if err != nil {
		return
	}

	c, err := GetCustomConfigFromViper()
	if err != nil {
		return
	}

	custom.CustomConfig = *c
	return
}
