package messaging

import (
	tntContext "github.com/Dert12318/Utilities/context"
)

const PublishTime = "publish_time"

type (
	Message struct {
		MsgID         string            `json:"msg_id"`
		MsgData       []byte            `json:"msg_data"`
		MsgAttributes map[string]string `json:"msg_attributes"`
	}

	Queue interface {
		Ping(ctx *tntContext.Context) error
		Listen()
		Subscribe(topic string, dispatcher Dispatcher) error
		Publish(topic string, msg Message) error
		PublishWithContext(ctx *tntContext.Context, topic string, msg Message) error
		Close() error
	}
)
