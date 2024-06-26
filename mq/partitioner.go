package mq

import (
	"hash/crc32"

	"github.com/IBM/sarama"
)

var random = sarama.NewRandomPartitioner("")

// UIDPartitioner UIDPartitioner
type UIDPartitioner struct {
	topic string
}

// NewUIDPartitioner NewUIDPartitioner
func NewUIDPartitioner(topic string) sarama.Partitioner {
	return &UIDPartitioner{
		topic: topic,
	}
}

func getMetaDataWithTag(message *sarama.ProducerMessage, tag string) (string, bool) {
	meta, ok := message.Metadata.(map[string]string)
	if !ok {
		return "", false
	}
	t, ok := meta[tag]
	return t, ok
}

// Partition 根据消息uid的前10位和，取余numPartitions
func (p *UIDPartitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	uid, ok := getMetaDataWithTag(message, "uid")
	if !ok {
		return random.Partition(message, numPartitions)
	}

	sum := UID2IndexSDBMHash(uid)
	return sum % numPartitions, nil
}

// RequiresConsistency 分区是否一致
func (p *UIDPartitioner) RequiresConsistency() bool {
	return true
}

// MessageRequiresConsistency 消息分区是否一致
func (p *UIDPartitioner) MessageRequiresConsistency(message *sarama.ProducerMessage) bool {
	return message.Key != nil
}

// UID2Index UID2Index
func UID2Index(uid string) int32 {
	var sum int32
	for _, v := range uid {
		// if index > 20 {
		// 	// 最多只取10位
		// 	break
		// }
		sum += v
	}
	return (sum & 0x7FFFFFFF)
}

// UID2IndexSDBMHash UID2IndexSDBMHash
func UID2IndexSDBMHash(uid string) int32 {
	var sum int32
	for _, v := range uid {
		// if index > 20 {
		// 	// 最多只取10位
		// 	break
		// }
		sum = v + (sum << 6) + (sum << 16) - sum
	}
	return (sum & 0x7FFFFFFF)
}

// UID2IndexCRC32 UID2IndexCRC32
func UID2IndexCRC32(uid string) int32 {
	return int32(crc32.ChecksumIEEE([]byte(uid)))
}
