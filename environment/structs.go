package environment

import (
	"fmt"
	"time"

	"github.com/AlessandroSanino1994/gobot/strategies"
)

//OrderType is an enum {ASK, BID}
type OrderType int16

const (
	Ask OrderType = iota //Represents an ASK Order.
	Bid OrderType = iota //Represents a BID Order.
)

//CandleStick represents a single candle in the graph.
type CandleStick struct {
	High  float64 //Represents the highest value obtained during candle period.
	Open  float64 //Represents the first value of the candle period.
	Close float64 //Represents the last value of the candle period.
	Low   float64 //Represents the lowest value obtained during candle period.
}

//Order represents a single order in the Order Book for a market.
type Order struct {
	Type         OrderType
	Value        float64
	Sum          float64
	Quantity     float64
	Countervalue float64
}

//CandleStickChart represents a chart of a market expresed using Candle Sticks.
type CandleStickChart struct {
	CandlePeriod time.Duration //Represents the candle period (expressed in time.Duration)
	MarketName   string        //Represents the name of the market (e.g. BTC-LTC)
	CandleSticks []CandleStick //Represents the last Candle Sticks used for evaluation of current state.
	OrderBook    []Order       //Represents the Book of current trades.
}

//Market represents the environment the bot is trading in.
type Market struct {
	Name         string              //Represents the name of the market (e.g. ETH-BTC).
	Volume24h    float64             //Represents the 24 hour volume of the market.
	WatchedChart CandleStickChart    //Represents a map which contains all the charts watched currently by the bot.
	strategy     strategies.Strategy //Represents the current strategy for the chart.
}

//ApplyStrategy returns an action as a consequence for applying a strategy of the market.
func (market *Market) ApplyStrategy() strategies.Action {
	return market.strategy.OnCandleUpdate()
}

func (market *Market) AttachStrategy(strategy strategies.Strategy) {
	if market.strategy != nil {
		market.DetachStrategy()
	}
	market.strategy = strategy
	market.strategy.Setup()
}

func (market *Market) DetachStrategy() error {
	if market.strategy == nil {
		return fmt.Errorf("Market %s error : No strategy to detach", market.Name)
	}
	market.strategy.TearDown()
	market.strategy = nil
	return nil
}
