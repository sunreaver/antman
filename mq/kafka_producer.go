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

// MakeKafkaProducer MakeKafkaProducer
func MakeKafkaProducer(c KafkaProducerConfig) (Sender, error) {
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

	return &KafkaProducer{
		kafka: kafka,
		cfg:   c,
		done:  make(chan bool),
		log:   logger.Empty,
	}, nil
}

// KafkaProducer KafkaProducer
type KafkaProducer struct {
	kafka sarama.AsyncProducer
	cfg   KafkaProducerConfig
	done  chan bool
	log   logger.Logger
}

// Stop Stop
func (m *KafkaProducer) Stop() {
	close(m.done)
}

// SyncSend SyncSend
func (m *KafkaProducer) SyncSend(t, st uint16, id string, body []byte) error {
	select {
	case m.kafka.Input() <- &sarama.ProducerMessage{
		Topic: fmt.Sprintf("%s%v", m.cfg.TopicPrefix, t),
		Key:   sarama.StringEncoder(fmt.Sprintf("%v", st)),
		Value: sarama.ByteEncoder(body),
		Metadata: map[string]string{
			"uid": id,
		},
	}:
	case err := <-m.kafka.Errors():
		return fmt.Errorf("kafka err: %v", err)
	case <-m.done:
		return errServerDone
	}
	return nil
}

func (m *KafkaProducer) SyncSendWithStringTopic(topic, key, id string, data []byte) error {
	select {
	case m.kafka.Input() <- &sarama.ProducerMessage{
		Topic: fmt.Sprintf("%s%v", m.cfg.TopicPrefix, topic),
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
		Metadata: map[string]string{
			"uid": id,
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
func (m *KafkaProducer) SetLogger(l logger.Logger) {
	m.log = l
}
