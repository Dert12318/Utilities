package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"

	"github.com/Dert12318/Utilities/apm"
	"github.com/Dert12318/Utilities/logs"
)

const (
	DefaultConsumerWorker       = 10
	DefaultConsumerRetryMax     = 3
	DefaultConsumerRetryBackoff = 2 * time.Second
	DefaultFailedDeadline       = 60 * time.Second
	DefaultStrategy             = "BalanceStrategyRoundRobin"
	DefaultHeartbeat            = 3
	DefaultProducerMaxBytes     = 1000000
	DefaultProducerRetryMax     = 3
	DefaultProducerRetryBackoff = 100
	BalanceStrategySticky       = "BalanceStrategySticky"
	BalanceStrategyRoundRobin   = "BalanceStrategyRoundRobin"
	BalanceStrategyRange        = "BalanceStrategyRange"
)

type (
	Option interface {
		Apply(o *option)
	}
	option struct {
		Host                 []string
		ClientID             string
		ConsumerWorker       int
		ConsumerGroup        string
		ConsumerRetryMax     int
		ConsumerRetryBackoff time.Duration
		Strategy             Strategy
		Heartbeat            int
		ProducerMaxBytes     int
		ProducerRetryMax     int
		ProducerRetryBackOff time.Duration
		KafkaVersion         string
		ListTopics           []string
		Log                  logs.Logger
		Apm                  apm.APM
		WithoutProducer      bool
		WithoutConsumer      bool
		EnableSASL           bool
		SASLMechanism        string
		Username             string
		Password             string
	}
)

func validate(option option) error {
	if len(option.Host) < 1 {
		return errors.New("invalid kafka host")
	}
	if option.ConsumerGroup == "" {
		return errors.New("invalid kafka consumer group")
	}
	if option.KafkaVersion == "" {
		return errors.New("invalid kafka version")
	}
	return nil
}

func getClientGeneratorFunc(mechanism sarama.SASLMechanism) func() sarama.SCRAMClient {
	if mechanism == sarama.SASLTypeSCRAMSHA512 {
		return func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
	} else if mechanism == sarama.SASLTypeSCRAMSHA256 {
		return func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
	}
	return nil
}

func getMechanism(option option) sarama.SASLMechanism {
	if option.SASLMechanism == sarama.SASLTypeSCRAMSHA512 {
		return sarama.SASLTypeSCRAMSHA512
	}

	if option.SASLMechanism == sarama.SASLTypeSCRAMSHA256 {
		return sarama.SASLTypeSCRAMSHA256
	}

	return sarama.SASLTypePlaintext
}

func getStrategy(option option) []sarama.BalanceStrategy {
	if option.Strategy == BalanceStrategyRange {
		return []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	}

	if option.Strategy == BalanceStrategyRoundRobin {
		return []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	}

	return []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
}

type withHost []string

func WithHost(host []string) Option {
	return withHost(host)
}

func (w withHost) Apply(o *option) {
	o.Host = w
}

type withClientID string

func WithClientID(clientID string) Option {
	return withClientID(clientID)
}

func (w withClientID) Apply(o *option) {
	o.ClientID = string(w)
}

type withConsumerWorker int

func WithConsumerWorker(worker int) Option {
	return withConsumerWorker(worker)
}

func (w withConsumerWorker) Apply(o *option) {
	o.ConsumerWorker = int(w)
}

type withConsumerGroup string

func WithConsumerGroup(group string) Option {
	return withConsumerGroup(group)
}

func (w withConsumerGroup) Apply(o *option) {
	o.ConsumerGroup = string(w)
}

type withConsumerRetryMax int

func WithConsumerRetryMax(maxRetry int) Option {
	return withProducerRetryMax(maxRetry)
}

func (w withConsumerRetryMax) Apply(o *option) {
	o.ProducerRetryMax = int(w)
}

type withStrategy Strategy

func WithStrategy(strategy Strategy) Option {
	return withStrategy(strategy)
}

func (w withStrategy) Apply(o *option) {
	o.Strategy = Strategy(w)
}

type withHeartbeat int

func WithHeartbeat(heartbeat int) Option {
	return withHeartbeat(heartbeat)
}

func (w withHeartbeat) Apply(o *option) {
	o.Heartbeat = int(w)
}

type withProducerMaxBytes int

func WithProducerMaxBytes(maxBytes int) Option {
	return withProducerMaxBytes(maxBytes)
}

func (w withProducerMaxBytes) Apply(o *option) {
	o.ProducerMaxBytes = int(w)
}

type withProducerRetryMax int

func WithProducerRetryMax(maxRetry int) Option {
	return withProducerRetryMax(maxRetry)
}

func (w withProducerRetryMax) Apply(o *option) {
	o.ProducerRetryMax = int(w)
}

type withProducerRetryBackOff time.Duration

func WithProducerRetryBackOff(retryBackoff time.Duration) Option {
	return withProducerRetryBackOff(retryBackoff)
}

func (w withProducerRetryBackOff) Apply(o *option) {
	o.ProducerRetryBackOff = time.Duration(w)
}

type withKafkaVersion string

func WithKafkaVersion(version string) Option {
	return withKafkaVersion(version)
}

func (w withKafkaVersion) Apply(o *option) {
	o.KafkaVersion = string(w)
}

type withListTopics []string

func WithListTopics(topics []string) Option {
	return withListTopics(topics)
}

func (w withListTopics) Apply(o *option) {
	o.ListTopics = w
}

type withLog struct{ logs.Logger }

func WithLog(logger logs.Logger) Option {
	return withLog{logger}
}

func (w withLog) Apply(o *option) {
	o.Log = w.Logger
}

type withApm struct{ apm.APM }

func WithApm(apm apm.APM) Option {
	return withApm{apm}
}

func (w withApm) Apply(o *option) {
	o.Apm = w.APM
}

type withoutProducer bool

func WithoutProducer() Option {
	return withoutProducer(true)
}

func (w withoutProducer) Apply(o *option) {
	o.WithoutProducer = bool(w)
}

type withoutConsumer bool

func WithoutConsumer() Option {
	return withoutConsumer(true)
}

func (w withoutConsumer) Apply(o *option) {
	o.WithoutConsumer = bool(w)
}

type enableSASL bool

func EnableSASL(enable bool) Option {
	return enableSASL(enable)
}

func (w enableSASL) Apply(o *option) {
	o.EnableSASL = bool(w)
}

type saslMechanism string

func SASLMechanism(mechanism string) Option {
	return saslMechanism(mechanism)
}

func (u saslMechanism) Apply(o *option) {
	o.SASLMechanism = string(u)
}

type username string

func Username(user string) Option {
	return username(user)
}

func (u username) Apply(o *option) {
	o.Username = string(u)
}

type password string

func Password(pass string) Option {
	return password(pass)
}

func (p password) Apply(o *option) {
	o.Password = string(p)
}
