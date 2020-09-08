package eventbus

// Coffee is the coffee schema
type Coffee struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Order is the schema for order, Status is of type Event but should be of type Status (separate type).
type Order struct {
	OrderID int
	Cof     Coffee
	Status  Event
}

// EventC is the complete event which is going to be pushed on kafka
type EventC struct {
	EventID int
	Event   Event
	Order   Order
}
