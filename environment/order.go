package environment

//OrderType is an enum {ASK, BID}
type OrderType int16

const (
	Ask OrderType = iota //Represents an ASK Order.
	Bid OrderType = iota //Represents a BID Order.
)

//Order represents a single order in the Order Book for a market.
type Order struct {
	Type        OrderType
	Value       float64
	Quantity    float64
	OrderNumber string
}

//Total returns order total in base currency.
func (order Order) Total() float64 {
	return order.Quantity * order.Value
}
