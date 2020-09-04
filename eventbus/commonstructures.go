package eventbus

// Coffee is the coffee schema
type Coffee struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Order is the schema for order
type Order struct {
	OrderID int
	Cof     Coffee
	Status  Event
}
