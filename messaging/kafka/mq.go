package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"

	tntContext "github.com/Dert12318/Utilities/context"
	"github.com/Dert12318/Utilities/encoding"
	jsoniter "github.com/Dert12318/Utilities/encoding/jsontier"
	"github.com/Dert12318/Utilities/logs/logrus"
	"github.com/Dert12318/Utilities/messaging"
)

type kafka struct {
	Option        option
	Client        sarama.Client
	ConsumerGroup sarama.ConsumerGroup
	consumer      *consumer
	jsonEncoding  encoding.Encoding
}

func New(options ...Option) (messaging.Queue, error) {
	var err error

	option := option{
		Host:                 make([]string, 0),
		KafkaVersion:         "",
		ConsumerGroup:        "",
		ListTopics:           make([]string, 0),
		ConsumerWorker:       DefaultConsumerWorker,
		ConsumerRetryMax:     DefaultConsumerRetryMax,
		ConsumerRetryBackoff: DefaultConsumerRetryBackoff,
		Strategy:             DefaultStrategy,
		Heartbeat:            DefaultHeartbeat,
		ProducerMaxBytes:     DefaultProducerMaxBytes,
		ProducerRetryMax:     DefaultProducerRetryMax,
		ProducerRetryBackOff: DefaultProducerRetryBackoff,
		Log:                  logrus.DefaultLog(),
		//Apm:                  defaultAPM,
	}

	for _, opt := range options {
		opt.Apply(&option)
	}

	if err := validate(option); err != nil {
		return nil, err
	}

	l := kafka{
		Option:       option,
		jsonEncoding: jsoniter.NewEncoding(),
	}

	version, err := sarama.ParseKafkaVersion(l.Option.KafkaVersion)
	if err != nil {
		return nil, err
	}

	sarama.Logger = l.Option.Log

	cfg := sarama.NewConfig()
	cfg.Version = version

	if option.EnableSASL {
		cfg.Net.SASL.Mechanism = getMechanism(l.Option)
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = option.Username
		cfg.Net.SASL.Password = option.Password
		cfg.Net.SASL.SCRAMClientGeneratorFunc = getClientGeneratorFunc(cfg.Net.SASL.Mechanism)
	}

	if len(l.Option.ClientID) > 0 {
		cfg.ClientID = l.Option.ClientID
	}

	if !option.WithoutConsumer {
		// - consumer
		cfg.Consumer.Group.Rebalance.GroupStrategies = getStrategy(l.Option)
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
		cfg.Consumer.Retry.Backoff = l.Option.ConsumerRetryBackoff
		cfg.Consumer.Return.Errors = true
		l.ConsumerGroup, err = sarama.NewConsumerGroup(l.Option.Host, l.Option.ConsumerGroup, cfg)
		if err != nil {
			return nil, err
		}
	}

	if !option.WithoutProducer {
		// - producer
		cfg.Producer.Return.Errors = true
		cfg.Producer.Return.Successes = true
		cfg.Producer.MaxMessageBytes = l.Option.ProducerMaxBytes
		cfg.Producer.Retry.Max = l.Option.ProducerRetryMax
		cfg.Producer.Retry.Backoff = l.Option.ProducerRetryBackOff
	}

	l.Client, err = sarama.NewClient(l.Option.Host, cfg)
	if err != nil {
		return nil, err
	}

	if !option.WithoutConsumer {
		l.consumer = &consumer{
			mu:     &sync.Mutex{},
			topics: make(map[string]messaging.Dispatcher),
			ready:  make(chan bool),
			option: l.Option,
			apm:    option.Apm,
		}
	}

	return &l, nil
}

func (k *kafka) Subscribe(topic string, dispatcher messaging.Dispatcher) error {
	if k.Option.WithoutConsumer {
		return errors.New("kafka is initialize without consumer")
	}

	k.consumer.mu.Lock()
	defer k.consumer.mu.Unlock()

	k.consumer.topics[topic] = dispatcher
	return nil
}

func (k *kafka) Listen() {
	if k.Option.WithoutConsumer {
		k.Option.Log.Error(errors.New("kafka is initialize without consumer"))
		return
	}

	if k.consumer.listening {
		k.Option.Log.Info("already listening to kafka")
		return
	}

	var (
		topics = make([]string, 0)
	)

	for key := range k.consumer.topics {
		topics = append(topics, key)
	}

	go func() {
		for {
			if err := k.ConsumerGroup.Consume(context.Background(), topics, k.consumer); err != nil {
				k.Option.Log.Errorf("error from consumer: %s", err.Error())
			}

			k.consumer.ready <- true
		}
	}()

	// Await till the consumer has been set up
	<-k.consumer.ready
}

func (k *kafka) Publish(topic string, msg messaging.Message) error {
	if k.Option.WithoutProducer {
		return errors.New("kafka is initialize without producer")
	}
	ctx := tntContext.New()
	ctx.Transaction = k.Option.Apm.StartTransaction(EventPublish)

	asyncProducer, err := sarama.NewAsyncProducerFromClient(k.Client)
	if err != nil {
		return errors.New(fmt.Sprintf("error create async client message: %s", err.Error()))
	}

	producer := &producer{asyncProducer: asyncProducer, apm: k.Option.Apm}
	defer func() {
		if err := producer.Close(); err != nil {
			k.Option.Log.Error(errors.Wrapf(err, "Failed to Close producer"))
		}
	}()

	return producer.Publish(ctx, topic, msg)
}

func (k *kafka) PublishWithContext(ctx *tntContext.Context, topic string, msg messaging.Message) error {
	if k.Option.WithoutProducer {
		return errors.New("kafka is initialize without producer")
	}
	if ctx.Transaction == nil {
		ctx.Transaction = k.Option.Apm.StartTransaction(EventPublish)
	}

	asyncProducer, err := sarama.NewAsyncProducerFromClient(k.Client)
	if err != nil {
		return errors.New(fmt.Sprintf("error create async client message: %s", err.Error()))
	}

	producer := &producer{asyncProducer: asyncProducer, apm: k.Option.Apm}
	defer func() {
		if err := producer.Close(); err != nil {
			k.Option.Log.Error(errors.Wrapf(err, "Failed to Close producer"))
		}
	}()

	return producer.Publish(ctx, topic, msg)
}

func (k *kafka) Close() error {
	if k.ConsumerGroup != nil {
		if err := k.ConsumerGroup.Close(); err != nil {
			return errors.Wrapf(err, "Failed to Close Consumer")
		}
	}

	if k.Client != nil {
		if err := k.Client.Close(); err != nil {
			return errors.Wrapf(err, "Failed to Close Producer")
		}
	}

	return nil
}

func (k *kafka) Ping(ctx *tntContext.Context) error {
	if k.consumer.listening {
		return nil
	}
	return errors.New("kafka is not rebalanced")
}
