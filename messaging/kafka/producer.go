package kafka

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"

	"github.com/Dert12318/Utilities/apm"
	tntContext "github.com/Dert12318/Utilities/context"
	"github.com/Dert12318/Utilities/messaging"
)

const (
	EventPublish = "kafka_publish"
)

type (
	producer struct {
		asyncProducer sarama.AsyncProducer
		apm           apm.APM
	}
)

func (p *producer) Publish(ctx *tntContext.Context, topic string, msg messaging.Message) error {
	//s := ctx.Transaction.StartMessageProducerSegment(apm.MessageProducerSegmentDTO{
	//	Library:              "Kafka",
	//	DestinationType:      apm.MessageTopic,
	//	DestinationName:      topic,
	//	DestinationTemporary: false,
	//})
	//defer s.End()

	headers := make([]sarama.RecordHeader, 0)

	if msg.MsgAttributes == nil {
		msg.MsgAttributes = make(map[string]string)
	}

	//msg.MsgAttributes[header.MessagingApiKey] = ctx.MandatoryRequest().APIKey()
	//msg.MsgAttributes[header.MessagingRequestID] = ctx.MandatoryRequest().RequestID()
	//msg.MsgAttributes[header.MessagingServiceID] = ctx.MandatoryRequest().ServiceID()
	//msg.MsgAttributes[header.MessagingServiceSecret] = ctx.MandatoryRequest().ServiceSecret()
	//msg.MsgAttributes[header.MessagingAuthorization] = ctx.MandatoryRequest().Token()
	//msg.MsgAttributes[header.MessagingDeviceID] = ctx.MandatoryRequest().DeviceID()
	//msg.MsgAttributes[header.MessagingUserAgent] = ctx.MandatoryRequest().UserAgent()

	for key, attr := range msg.MsgAttributes {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(attr),
		})
	}
	headers = append(headers, sarama.RecordHeader{
		Key:   []byte(messaging.PublishTime),
		Value: []byte(time.Now().Format(time.RFC3339)),
	})

	if "" == msg.MsgID {
		msg.MsgID = uuid.New().String()
	}

	p.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(msg.MsgID),
		Value:     sarama.StringEncoder(string(msg.MsgData)),
		Headers:   headers,
		Timestamp: time.Now(),
	}

	if p.apm != nil {
		key := fmt.Sprintf("%s:%s", EventPublish, strings.ReplaceAll(topic, ".", "_"))
		p.apm.RecordCustomEvent(key, map[string]interface{}{
			"topic":        topic,
			"message_id":   msg.MsgID,
			"message_attr": msg.MsgAttributes,
			"message_data": string(msg.MsgData),
		})
	}
	return nil
}

func (p *producer) Close() error {
	return p.asyncProducer.Close()
}
