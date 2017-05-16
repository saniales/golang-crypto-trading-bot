package environment

import (
	"fmt"

	"github.com/AlessandroSanino1994/gobot/strategies"
)

//Market represents the environment the bot is trading in.
type Market struct {
	Name           string              //Represents the name of the market (e.g. ETH-BTC).
	BaseCurrency   string              //Represents the base currency of the market.
	MarketCurrency string              //Represents the currency to exchange by using base currency.
	WatchedChart   *CandleStickChart   //Represents a map which contains all the charts watched currently by the bot.
	strategy       strategies.Strategy //Represents the current strategy for the chart.
}

//ApplyStrategy returns an action as a consequence for applying a strategy of the market.
func (market *Market) ApplyStrategy() strategies.Action {
	return market.strategy.OnCandleUpdate(market)
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
