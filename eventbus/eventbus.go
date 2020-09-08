package eventbus

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

// Event matches event types
type Event string

const (
	orderAccepted      = "order_accepted"
	orderReceived      = "order_received"
	orderProcessing    = "order_processing"
	orderRejected      = "order_rejected"
	orderFailed        = "order_failed"
	orderDelivered     = "order_delivered"
	orderProcessed     = "order_processed"
	coffeeNotAvailable = "coffee_not_available"
)

const (
	// OrderAccept is the event emitted by Beans Service after beans validation.
	OrderAccept Event = orderAccepted

	// OrderReceived is the event emitted by order service when the request is received.
	OrderReceived Event = orderReceived

	// OrderProcessing is the event emitted by brewerie every 5 seconds.
	OrderProcessing Event = orderProcessing

	// OrderRejected is the event emitted by beans service if fails to validate beans.
	OrderRejected Event = orderRejected

	// OrderFailed is the event emitted by brewerie service if it fails to complete the order.
	OrderFailed Event = orderFailed

	// OrderDelivered is the event emitted by order service after delivering the order.
	OrderDelivered Event = orderDelivered

	// OrderProcessed is the event emitted by brewerie service after completing the order.
	OrderProcessed Event = orderProcessed
)

// Publisher is the interface for the eventbus publish behaviour
type Publisher interface {
	Publish(ctx context.Context, e Event, jsonData []byte) error
}

// Reader is the interface for the eventbus read behaviour
type Reader interface {
	ReadEvents(ctx context.Context) (kafka.Message, error)
}

// E is an interface to describe the entirety of eventbus.
type E interface {
	Publisher
	Reader
}

// Config is the eventbus configuration
type Config struct {
	// ConsumerName is the name of the service
	ConsumerName string

	// Logger is the logger
	Logger *zerolog.Logger

	// Kafka topic to produce event on
	Topic string

	// KafkaClient is the *kafka.AdminClient
	KafkaClient *kafka.Client
}

// EB is the eventbus struct, which satisfies E
type EB struct {
	w *kafka.Writer
	r *kafka.Reader

	l *zerolog.Logger
}

// compile time check: does *EB satisfy E interface?
var _ E = (*EB)(nil)

// New returns an instance
func New(cfg Config) (*EB, error) {

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   cfg.Topic,
	})

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   cfg.Topic,
	})

	eb := &EB{
		w: w,
		r: r,
		l: cfg.Logger,
	}
	return eb, nil
}

// Publish publishes to kafka topic
func (eb *EB) Publish(ctx context.Context, e Event, jsonData []byte) error {
	// e is the event type
	eb.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(e),
		Value: jsonData,
	})
	return nil
}

// ReadEvents to read event
func (eb *EB) ReadEvents(ctx context.Context) (kafka.Message, error) {
	msg, err := eb.r.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, err
	}
	return msg, nil
}
