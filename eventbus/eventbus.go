package eventbus

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
)

// Event matches event types
type Event string

const (
	orderAccepted   = "order_accepted"
	orderReceived   = "order_received"
	orderProcessing = "order_processing"
	orderRejected   = "order_rejected"
	orderFailed     = "order_failed"
	orderDelivered  = "order_delivered"
	orderProcessed  = "order_processed"
)

const (
	// OrderAccept is the event emitted by Beans Service after beans validation.
	OrderAccept Event = orderAccepted

	// OrderReceived is the event emitted by order service when the request is received .
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
	Publish(e Event, eventTimestamp int64, eventID, requetID string, jsonData []byte) error
}

// E is an interface to describe the entirety of eventbus.
type E interface {
	Publisher
}

// Config is the eventbus configuration
type Config struct {
	// ConsumerName is the name of the service
	ConsumerName string

	// Logger is the logger
	Logger *zerolog.Logger

	// KafkaClient is the *kafka.AdminClient
	KafkaClient *kafka.AdminClient
}

// EB is the eventbus struct, which satisfies Q
type EB struct {
	p *kafka.Producer
	c *kafka.Consumer

	l *zerolog.Logger
}

// compile time check: does *EB satisfy E interface
var _ E = (*EB)(nil)

// Publish publishes to kafka topic
func (eb *EB) Publish(e Event, eventTimestamp int64, eventID, requestID string, jsonData []byte) error {
	return nil
}

// New returns an instance
func New() {

}
