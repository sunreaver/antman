package mq

import "github.com/sunreaver/logger/v3"

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

type AsyncSender interface {
	Stoper
	Logger
	AsyncSend(topic, key uint16, id string, data []byte) error
	AsyncSendWithStringTopic(topic, key, id string, data []byte) error
}

type SyncSender interface {
	Stoper
	Logger
	SyncSend(topic, key uint16, id string, data []byte) error
	SyncSendWithStringTopic(topic, key, id string, data []byte) error
}

// Recver Recver
type Recver interface {
	Stoper
	Logger
	SyncRecvUintTopic(RecvUIntTopicFunc) error
	SyncRecvStringTopic(RecvStringTopicFunc) error
}
