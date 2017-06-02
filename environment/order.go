package environment

//OrderType is an enum {ASK, BID}
type OrderType int16

const (
	//Ask Represents an ASK Order.
	Ask OrderType = iota
	//Bid Represents a BID Order.
	Bid OrderType = iota
)

//Order represents a single order in the Order Book for a market.
type Order struct {
	Type        OrderType //Type of the order. Can be Ask or Bid.
	Value       float64   //Value of the trade : e.g. in a BTC ETH is the value of a single ETH in BTC.
	Quantity    float64   //Quantity of Coins of this order.
	OrderNumber string    //[optional]Order number as seen in echange archives.
}

//Total returns order total in base currency.
func (order Order) Total() float64 {
	return order.Quantity * order.Value
}
