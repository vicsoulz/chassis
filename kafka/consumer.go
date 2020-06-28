package kafka

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

type ConsumerCallback func(message *sarama.ConsumerMessage) error

//NewConsumer return Consumer instance
func NewConsumer(custom *ConsumerCustomConfig) (*cluster.Consumer, error) {
	conf, err := NewConsumerConfig(custom)
	if err != nil {
		return nil, err
	}

	return cluster.NewConsumer(
		custom.Address,
		custom.ConsumerGroup,
		custom.Topics,
		conf)
}

func RunConsumer(consumer *cluster.Consumer, callBack ConsumerCallback) {
	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg, more := <-consumer.Messages():
			if more {
				if err := callBack(msg); err == nil {
					consumer.MarkOffset(msg, "") // mark message as processed
				} else {
					fmt.Printf("consumer handle msg error %+v", err)
				}
			}
		case err, more := <-consumer.Errors():
			if more {
				fmt.Printf("Kafka consumer error: %v", err.Error())
			}
		case ntf, more := <-consumer.Notifications():
			if more {
				fmt.Printf("Kafka consumer rebalance: %v", ntf)
			}
		case <-signals:
			return
		}
	}

}
