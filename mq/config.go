package mq

// KafkaGeneralConfig KafkaGeneralConfig
type KafkaGeneralConfig struct {
	TopicPrefix string   `toml:"topic_prefix"`
	Hosts       []string `toml:"hosts"`
	Version     string   `toml:"version"`
}

// KafkaProducerConfig KafkaProducerConfig
type KafkaProducerConfig struct {
	KafkaGeneralConfig
}

// KafkaConsumerConfig KafkaConsumerConfig
type KafkaConsumerConfig struct {
	KafkaGeneralConfig
	Topic []uint16 `toml:"topics"`
	Group string   `toml:"group"`
}
