package mq

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/sunreaver/logger"
)

var errServerDone = errors.New("server done")

// MakeKafkaAsyncProducer MakeKafkaAsyncProducer
func MakeKafkaAsyncProducer(c KafkaProducerConfig) (AsyncSender, error) {
	cfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(c.Version)
	if err != nil {
		return nil, errors.Wrap(err, "parse version")
	}
	cfg.Version = version
	cfg.Producer.Partitioner = NewUIDPartitioner
	kafka, e := sarama.NewAsyncProducer(c.Hosts, cfg)
	if e != nil {
		return nil, errors.Wrapf(e, "new async producer [hosts: %v]", c.Hosts)
	}
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case e := <-kafka.Errors():
				fmt.Println("[sarama] kafka producer err:", e.Error())
			}
		}
	}()

	return &KafkaAsyncProducer{
		kafka: kafka,
		cfg:   c,
		done:  done,
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
	case <-m.done:
		return errors.Wrap(errServerDone, "async done")
	}
	return nil
}

// SetLogger SetLogger
func (m *KafkaAsyncProducer) SetLogger(l logger.Logger) {
	m.log = l
}
