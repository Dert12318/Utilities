package kafka

import (
	"fmt"
	"sync"

	"github.com/Shopify/sarama"

	"github.com/Dert12318/Utilities/apm"
	"github.com/Dert12318/Utilities/common/constant/header"
	"github.com/Dert12318/Utilities/messaging"
)

type (
	Strategy string
	consumer struct {
		mu        *sync.Mutex
		topics    map[string]messaging.Dispatcher
		ready     chan bool
		listening bool
		option    option
		apm       apm.APM
	}
	messageShown struct {
		MsgID         string            `json:"msg_id"`
		MsgData       interface{}       `json:"msg_data"`
		MsgAttributes map[string]string `json:"msg_attributes"`
	}
)

func (c *consumer) Setup(session sarama.ConsumerGroupSession) error {
	if c.mu == nil {
		c.mu = &sync.Mutex{}
	}

	if c.topics == nil {
		c.topics = make(map[string]messaging.Dispatcher)
	}

	c.option.Log.Info("Start Listening")
	c.listening = true
	return nil
}

func (c *consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		dispatcher := c.topics[message.Topic]

		if message == nil {
			continue
		}

		if dispatcher == nil {
			continue
		}

		c.processMessage(message.Topic, session, message, dispatcher)
	}

	return nil
}

func (c *consumer) getMessage(msg sarama.ConsumerMessage) messaging.Message {
	attributes := make(map[string]string)
	for _, attr := range msg.Headers {
		attributes[string(attr.Key)] = string(attr.Value)
	}

	return messaging.Message{
		MsgID:         string(msg.Key),
		MsgData:       msg.Value,
		MsgAttributes: attributes,
	}
}

func (c *consumer) processMessage(topic string, session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage, dispatcher messaging.Dispatcher) {
	messageData := c.getMessage(*msg)
	requestID := messageData.MsgAttributes[header.MessagingRequestID]
	msgType := messageData.MsgAttributes["message"]

	var err error
	var message = messaging.DispatchDTO{
		Type:      messaging.Handle,
		Source:    fmt.Sprintf("Kafka - %s", topic),
		RequestID: requestID,
		MsgType:   msgType,
		Msg:       messageData,
		Log:       c.option.Log,
	}
	for i := 0; i <= c.option.ConsumerRetryMax; i++ {
		message.Err = nil
		if err = dispatcher.Dispatch(message); err == nil {
			session.MarkMessage(msg, "")
			return
		} else {
			c.option.Log.Error("error on dispatch message from kafka: ", err.Error())
			message.Err = err
		}
	}

	errMessage := messaging.DispatchDTO{
		Type:      messaging.Error,
		Source:    fmt.Sprintf("Kafka - %s", topic),
		RequestID: requestID,
		MsgType:   msgType,
		Msg: messaging.Message{
			MsgID:         messageData.MsgID,
			MsgData:       messageData.MsgData,
			MsgAttributes: messageData.MsgAttributes,
		},
		Log: c.option.Log,
		Err: err,
	}
	_ = dispatcher.Dispatch(errMessage)
	session.MarkMessage(msg, "")
}
