package mq

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/sunreaver/logger"
)

// MakeKafkaConsumer MakeKafkaConsumer
func MakeKafkaConsumer(c KafkaConsumerConfig) (Recver, error) {
	cfg := sarama.NewConfig()
	v, err := sarama.ParseKafkaVersion(c.Version)
	if err != nil {
		return nil, fmt.Errorf("parse version err: %v", err)
	}
	cfg.Version = v

	client, err := sarama.NewConsumerGroup(c.Hosts, c.Group, cfg)
	if err != nil {
		return nil, fmt.Errorf("new consumer err: %v", err)
	}
	topics := make([]string, len(c.Topic))
	for index, v := range c.Topic {
		topics[index] = fmt.Sprintf("%v%v", c.TopicPrefix, v)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaConsumer{
		cfg:            c,
		topics:         topics,
		topicPrefixLen: len(c.TopicPrefix),
		client:         client,
		ctx:            ctx,
		cancel:         cancel,
		log:            logger.Empty,
	}, nil
}

// KafkaConsumer represents a Sarama consumer group consumer
type KafkaConsumer struct {
	cfg            KafkaConsumerConfig
	topics         []string
	topicPrefixLen int
	ctx            context.Context
	cancel         context.CancelFunc
	client         sarama.ConsumerGroup
	intFN          RecvUIntTopicFunc
	stringFN       RecvStringTopicFunc
	log            logger.Logger
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29

PROCESS_MSG:
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				break PROCESS_MSG
			}
			var err error
			if consumer.intFN != nil {
				t, st, e := consumer.getTypeAndSubtype(msg)
				if e != nil {
					err = e
				} else {
					err = consumer.intFN(t, st, msg.Value)
				}
			} else if consumer.stringFN != nil {
				err = consumer.stringFN(msg.Topic, string(msg.Key), msg.Value)
			}
			if err != nil {
				consumer.log.Errorw("process",
					"topic", msg.Topic,
					"key", string(msg.Key),
					"mq_time", msg.Timestamp,
					"data", string(msg.Value),
					"err", err,
				)
			} else {
				consumer.log.Debugw("process",
					"topic", msg.Topic,
					"key", string(msg.Key),
					"mq_time", msg.Timestamp,
					"data", string(msg.Value),
				)
			}
			session.MarkMessage(msg, "")
		case <-consumer.ctx.Done():
			break PROCESS_MSG
		}
	}

	return nil
}

// SyncRecv SyncRecv
func (consumer *KafkaConsumer) SyncRecv(fn RecvUIntTopicFunc) error {
	consumer.intFN = fn
	return consumer.startConsume()
}

// SyncRecv SyncRecv
func (consumer *KafkaConsumer) SyncRecvStringTopic(fn RecvStringTopicFunc) error {
	consumer.stringFN = fn
	return consumer.startConsume()
}

// SyncRecv SyncRecv
func (consumer *KafkaConsumer) startConsume() error {
	errChan := make(chan error)
	go func() {
		for {
			err := consumer.client.Consume(consumer.ctx, consumer.topics, consumer)
			if err != nil {
				errChan <- err
				return
			}
		}
	}()
	consumer.log.Infow("start_sync_recv",
		"cfg", consumer.cfg,
	)
	defer consumer.client.Close()
	err := <-errChan
	return err
}

// Stop Stop
func (consumer *KafkaConsumer) Stop() {
	consumer.log.Infow("stop_sync_recv", "cfg", consumer.cfg)
	consumer.cancel()
}

func (consumer *KafkaConsumer) getTypeAndSubtype(msg *sarama.ConsumerMessage) (uint16, uint16, error) {
	if !strings.HasPrefix(msg.Topic, consumer.cfg.TopicPrefix) {
		return 0, 0, fmt.Errorf("unsupport kafka topic: %v", msg.Topic)
	}
	t, e := strconv.Atoi(msg.Topic[consumer.topicPrefixLen:])
	if e != nil {
		return 0, 0, fmt.Errorf("topic err: %v", e)
	}
	st, e := strconv.Atoi(string(msg.Key))
	if e != nil {
		return 0, 0, fmt.Errorf("st err: %v", e)
	}
	return uint16(t), uint16(st), nil
}

// SetLogger SetLogger
func (consumer *KafkaConsumer) SetLogger(l logger.Logger) {
	consumer.log = l
}
