package mq

import (
	"github.com/sunreaver/logger"
)

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
	SyncSend(uint16, uint16, string, []byte) error
}

// Recver Recver
type Recver interface {
	Stoper
	Logger
	SyncRecv(recvFunc) error
}
