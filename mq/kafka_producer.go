package mq

import (
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/sunreaver/logger"
)

var (
	errServerDone = errors.New("server done")
)

// MakeKafkaAsyncProducer MakeKafkaAsyncProducer
func MakeKafkaAsyncProducer(c KafkaProducerConfig) (AsyncSender, error) {
	cfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(c.Version)
	if err != nil {
		return nil, fmt.Errorf("parse version err: %v", err)
	}
	cfg.Version = version
	cfg.Producer.Partitioner = NewUIDPartitioner
	kafka, e := sarama.NewAsyncProducer(c.Hosts, cfg)
	if e != nil {
		return nil, fmt.Errorf("new async producer err: %v[hosts: %v]", e, c.Hosts)
	}

	return &KafkaAsyncProducer{
		kafka: kafka,
		cfg:   c,
		done:  make(chan bool),
		log:   logger.Empty,
	}, nil
}

// KafkaAsyncProducer KafkaAsyncProducer
type KafkaAsyncProducer struct {
	kafka sarama.AsyncProducer
	cfg   KafkaProducerConfig
	done  chan bool
	log   logger.Logger
}

// Stop Stop
func (m *KafkaAsyncProducer) Stop() {
	close(m.done)
}

// SyncSend SyncSend
func (m *KafkaAsyncProducer) AsyncSend(topic, key uint16, uid string, body []byte) error {
	return m.AsyncSendWithStringTopic(fmt.Sprintf("%v", topic), fmt.Sprintf("%v", key), uid, body)
}

func (m *KafkaAsyncProducer) AsyncSendWithStringTopic(topic, key, uid string, data []byte) error {
	select {
	case m.kafka.Input() <- &sarama.ProducerMessage{
		Topic: fmt.Sprintf("%s%v", m.cfg.TopicPrefix, topic),
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
		Metadata: map[string]string{
			"uid": uid,
		},
	}:
	case err := <-m.kafka.Errors():
		return fmt.Errorf("kafka err: %v", err)
	case <-m.done:
		return errServerDone
	}
	return nil
}

// SetLogger SetLogger
func (m *KafkaAsyncProducer) SetLogger(l logger.Logger) {
	m.log = l
}
