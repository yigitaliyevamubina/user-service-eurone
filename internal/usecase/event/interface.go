package event

import (
	"context"
)

type ConsumerConfig interface {
	GetBrokers() []string
	GetTopic() string
	GetGroupID() string
	GetHandler() func(ctx context.Context, key, value []byte) error
}

type BrokerConsumer interface {
	Run() error
	RegisterConsumer(config ConsumerConfig)
	Close()
}