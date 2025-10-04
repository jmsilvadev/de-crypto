package kafka

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jmsilvadev/de-crypto/pkg/config"
	k "github.com/segmentio/kafka-go"
)

type Publisher struct {
	ctx context.Context
	w   *k.Writer
}

func NewPublisher(ctx context.Context, brokers []string, topic string) (*Publisher, error) {
	log.Printf("registering new publisher Topic: %s", topic)
	if len(brokers) == 0 || topic == "" {
		return nil, errors.New("kafka: missing brokers or topic")
	}

	if err := ensureTopicExists(ctx, brokers, topic); err != nil {
		return nil, err
	}

	w := &k.Writer{
		Addr:         k.TCP(brokers...),
		Topic:        topic,
		Balancer:     &k.Hash{},
		RequiredAcks: k.RequireAll,
		BatchSize:    config.DefaultBatchSize,
		BatchTimeout: config.DefaultFlushInterval,
		WriteTimeout: config.DefaultRequestTimeout,
	}
	return &Publisher{ctx: ctx, w: w}, nil
}

func (p *Publisher) Publish(b []byte) error {
	return p.w.WriteMessages(p.ctx, k.Message{Value: b})
}

func (p *Publisher) Close() error {
	return p.w.Close()
}

func ensureTopicExists(ctx context.Context, brokers []string, topic string) error {
	conn, err := k.DialLeader(ctx, "tcp", brokers[0], topic, 0)
	if err != nil {
		maxRetries := 3
		for attempt := 1; attempt <= maxRetries; attempt++ {
			admin := k.Client{
				Addr: k.TCP(brokers...),
			}

			_, createErr := admin.CreateTopics(ctx, &k.CreateTopicsRequest{
				Topics: []k.TopicConfig{
					{
						Topic:             topic,
						NumPartitions:     1,
						ReplicationFactor: 1,
					},
				},
			})

			if createErr != nil {
				log.Printf("Attempt %d: Failed to create topic %s: %v", attempt, topic, createErr)
				if attempt == maxRetries {
					return createErr
				}
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}

			log.Printf("Successfully created topic: %s", topic)
			break
		}

		time.Sleep(2 * time.Second)

		conn, err = k.DialLeader(ctx, "tcp", brokers[0], topic, 0)
		if err != nil {
			log.Printf("Failed to connect to topic %s after creation: %v", topic, err)
			return err
		}
	}

	conn.Close()
	return nil
}
