package mq

import (
	"github.com/sunreaver/logger"
)

type RecvUIntTopicFunc func(topic, key uint16, data []byte) error
type RecvStringTopicFunc func(topic, key string, data []byte) error

// Stoper Stoper
type Stoper interface {
	Stop()
}

// Logger Logger
type Logger interface {
	SetLogger(logger.Logger)
}

// Sender Sender
type Sender interface {
	Stoper
	Logger
	SyncSend(topic uint16, key uint16, id string, data []byte) error
	SyncSendWithStringTopic(topic, key, id string, data []byte) error
}

// Recver Recver
type Recver interface {
	Stoper
	Logger
	SyncRecv(RecvUIntTopicFunc) error
	SyncRecvStringTopic(RecvStringTopicFunc) error
}
