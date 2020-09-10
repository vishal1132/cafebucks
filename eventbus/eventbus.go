package eventbus

import (
	"context"
	"fmt"
	"log"

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

	// GroupID is the consumer group
	GroupID string

	// Brokers is kafka brokers
	Brokers []string
}

// EB is the eventbus struct, which satisfies E
type EB struct {
	w *kafka.Writer
	r *kafka.Reader

	l   *zerolog.Logger
	g   *kafka.ConsumerGroup
	gen *kafka.Generation
	c   *kafka.Client
}

// GetOffset returns the offset for the reader
func (eb *EB) GetOffset(ctx context.Context, topic, group string) (int64, error) {
	m, err := eb.c.ConsumerOffsets(ctx, kafka.TopicAndGroup{
		Topic:   topic,
		GroupId: group,
	})
	if err != nil {
		return -5, err
	}
	return m[0], nil
}

// SetOffset sets the offset for the reader
func (eb *EB) SetOffset(offset int64) error {
	return eb.r.SetOffset(offset)
}

// GetOffsetGen returns the offset for the group
func (eb *EB) GetOffsetGen(ctx context.Context, topic string) (int64, error) {

	assignments := eb.gen.Assignments[topic]
	var offset int64
	for _, assignment := range assignments {
		_, offset = assignment.ID, assignment.Offset
	}

	return offset, nil

	// return eb.gen.CommitOffsets(map[string]map[int]int64{topic: {0: offset}})
}

// SetOffsetGen sets the offset for the group and returns error if any
func (eb *EB) SetOffsetGen(ctx context.Context, topic string, offset int64) error {
	var err error
	eb.gen, err = eb.g.Next(ctx)
	if err != nil {
		return err
	}
	return eb.gen.CommitOffsets(map[string]map[int]int64{topic: {0: offset}})
}

// compile time check: does *EB satisfy E interface?
var _ E = (*EB)(nil)

// New returns an instance
func New(cfg Config) (*EB, error) {

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
	})

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})

	// group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
	// 	ID:      cfg.GroupID,
	// 	Brokers: cfg.Brokers,
	// 	Topics:  []string{"my-topic"},
	// })

	// if err != nil {
	// 	return nil, err
	// }

	c := kafka.NewClient(cfg.Brokers...)

	eb := &EB{
		w: w,
		r: r,
		l: cfg.Logger,
		c: c,
	}

	return eb, nil
}

// NewReader returns a reader
func (eb *EB) NewReader(brokers []string, topic string) *EB {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return &EB{r: r}
}

// NewGen creates a new generation from kafka group
func (eb *EB) NewGen(ctx context.Context) error {
	var err error
	eb.gen, err = eb.g.Next(ctx)
	return err
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

// StartGen starts a new gen
func (eb *EB) StartGen() {

}

// ReadEventsGen read events for a consumer group
func (eb *EB) ReadEventsGen(ctx context.Context, topic string) {
	var offset int64
	eb.gen.Start(func(ctx context.Context) {
		// create reader for this partition.
		reader := eb.r
		defer reader.Close()
		reader.SetOffset(offset)
		for {
			msg, err := reader.ReadMessage(ctx)
			switch err {
			case kafka.ErrGenerationEnded:
				// generation has ended.  commit offsets.  in a real app,
				// offsets would be committed periodically.
				log.Println("committing offset now")
				eb.gen.CommitOffsets(map[string]map[int]int64{"test": {0: offset + 1}})
				return
			case nil:
				offset = msg.Offset
				return
			default:
				fmt.Printf("error reading message: %+v\n", err)
				return
			}
		}
	})

}
