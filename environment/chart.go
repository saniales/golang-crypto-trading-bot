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

//CandleStick represents a single candle in the graph.
import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

//CandleStick represents a single candlestick in a chart.
type CandleStick struct {
	High   decimal.Decimal //Represents the highest value obtained during candle period.
	Open   decimal.Decimal //Represents the first value of the candle period.
	Close  decimal.Decimal //Represents the last value of the candle period.
	Low    decimal.Decimal //Represents the lowest value obtained during candle period.
	Volume decimal.Decimal //Represents the volume of trades during the candle period.
}

// String returns the string representation of the object.
func (cs CandleStick) String() string {
	var color string
	if cs.Open.GreaterThan(cs.Close) {
		color = "Green/Bullish"
	} else if cs.Open.LessThan(cs.Close) {
		color = "Red/Bearish"
	} else {
		color = "Neutral"
	}
	ret := fmt.Sprintln(color, "Candle")
	ret += fmt.Sprintln("High:", cs.High)
	ret += fmt.Sprintln("Open:", cs.Open)
	ret += fmt.Sprintln("Close:", cs.Close)
	ret += fmt.Sprintln("Low:", cs.Low)
	ret += fmt.Sprintln("Volume:", cs.Volume)
	return strings.TrimSpace(ret)
}

//CandleStickChart represents a chart of a market expresed using Candle Sticks.
type CandleStickChart struct {
	CandlePeriod time.Duration //Represents the candle period (expressed in time.Duration).
	CandleSticks []CandleStick //Represents the last Candle Sticks used for evaluation of current state.
	OrderBook    []Order       //Represents the Book of current trades.
}
