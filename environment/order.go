// Copyright © 2017 Alessandro Sanino <saninoale@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package environment

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

//OrderType is an enum {ASK, BID}
type OrderType int16

const (
	//Ask Represents an ASK Order.
	Ask OrderType = iota
	//Bid Represents a BID Order.
	Bid OrderType = iota
)

//OrderBook represents a standard orderbook implementation.
type OrderBook struct {
	Asks []Order `json:"asks,required"`
	Bids []Order `json:"bids,required"`
}

// String returns the string representation of the object.
func (book OrderBook) String() string {
	return fmt.Sprintln("ASKS") +
		fmt.Sprintln(book.Asks) +
		fmt.Sprintln("BIDS") +
		fmt.Sprintln(book.Bids)
}

//Order represents a single order in the Order Book for a market.
type Order struct {
	Value       decimal.Decimal //Value of the trade : e.g. in a BTC ETH is the value of a single ETH in BTC.
	Quantity    decimal.Decimal //Quantity of Coins of this order.
	OrderNumber string          //[optional] Order number as seen in echange archives.
	Timestamp   time.Time       //[optional] The timestamp of the order (as got from the exchange).
}

//Total returns order total in base currency.
func (order Order) Total() decimal.Decimal {
	return order.Quantity.Mul(order.Value)
}
