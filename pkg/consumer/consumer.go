package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	kafkaPlayloads "github.com/crodriguezde/rtdashs/pkg/kafkaPayloads"
)

type ConsumerGroupHandler interface {
	sarama.ConsumerGroupHandler
	WaitReady()
	Reset()
}

type ConsumerGroup struct {
	cg sarama.ConsumerGroup
}

func NewConsumerGroup(broker string, topics []string, group string, handler ConsumerGroupHandler, ver string) (*ConsumerGroup, error) {
	ctx := context.Background()
	cfg := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(ver)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	cfg.Version = version
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	client, err := sarama.NewConsumerGroup([]string{broker}, group, cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err := client.Consume(ctx, topics, handler)
			if err != nil {
				if err == sarama.ErrClosedConsumerGroup {
					break
				} else {
					panic(err)
				}
			}
			if ctx.Err() != nil {
				return
			}
			handler.Reset()
		}
	}()

	handler.WaitReady() // Await till the consumer has been set up

	return &ConsumerGroup{
		cg: client,
	}, nil
}

func (c *ConsumerGroup) Close() error {
	return c.cg.Close()
}

type ConsumerSessionMessage struct {
	Session sarama.ConsumerGroupSession
	Message *sarama.ConsumerMessage
}

func decodeMessage(data []byte) error {
	var msg kafkaPlayloads.Generic
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}
	return nil
}

func StartSync(broker, topic string, version string) (*ConsumerGroup, error) {
	var count int64
	var start = time.Now()
	handler := NewSyncConsumerGroupHandler(func(data []byte) error {
		if err := decodeMessage(data); err != nil {
			return err
		}
		count++
		if count%5000 == 0 {
			fmt.Printf("sync consumer consumed %d messages at speed %.2f/s\n", count, float64(count)/time.Since(start).Seconds())
		}
		return nil
	})
	consumer, err := NewConsumerGroup(broker, []string{topic}, "sync-consumer-"+fmt.Sprintf("%d", time.Now().Unix()), handler, version)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
