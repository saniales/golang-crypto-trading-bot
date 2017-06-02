package environment

//CandleStick represents a single candle in the graph.
import "time"

//CandleStick represents a single candlestick in a chart.
type CandleStick struct {
	High  float64 //Represents the highest value obtained during candle period.
	Open  float64 //Represents the first value of the candle period.
	Close float64 //Represents the last value of the candle period.
	Low   float64 //Represents the lowest value obtained during candle period.
}

//CandleStickChart represents a chart of a market expresed using Candle Sticks.
type CandleStickChart struct {
	CandlePeriod time.Duration //Represents the candle period (expressed in time.Duration).
	CandleSticks []CandleStick //Represents the last Candle Sticks used for evaluation of current state.
	OrderBook    []Order       //Represents the Book of current trades.
}
