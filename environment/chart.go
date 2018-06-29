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
	High  decimal.Decimal //Represents the highest value obtained during candle period.
	Open  decimal.Decimal //Represents the first value of the candle period.
	Close decimal.Decimal //Represents the last value of the candle period.
	Low   decimal.Decimal //Represents the lowest value obtained during candle period.
}

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
	return strings.TrimSpace(ret)
}

//CandleStickChart represents a chart of a market expresed using Candle Sticks.
type CandleStickChart struct {
	CandlePeriod time.Duration //Represents the candle period (expressed in time.Duration).
	CandleSticks []CandleStick //Represents the last Candle Sticks used for evaluation of current state.
	OrderBook    []Order       //Represents the Book of current trades.
}
