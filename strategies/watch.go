package strategies

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AlessandroSanino1994/gobot/environment"
	"github.com/AlessandroSanino1994/gobot/exchangeWrappers"
)

func init() {
	All["watch"] = WatchStrategy{RefreshEvery: "hour"}
}

// WatchStrategy represents a strategy which does nothing than print
// what it gets from markets.
type WatchStrategy struct {
	RefreshEvery string // Interval represents how often summaries will be requested.
}

// Name returns the name of the strategy.
func (w WatchStrategy) Name() string {
	return "Watch"
}

// RefreshInterval represents how often summaries will be requested.
func (w WatchStrategy) RefreshInterval() (time.Duration, error) {
	return time.ParseDuration(w.RefreshEvery)
}

// OnCandleUpdate Prints info about markets on every candle tick in JSON format.
func (w WatchStrategy) OnCandleUpdate(wrapper exchangeWrappers.ExchangeWrapper, market environment.Market) (Action, float64, float64, error) {
	err := wrapper.GetMarketSummary(&market)
	if err != nil {
		return Invalid, -1, -1, err
	}

	JSONContent, err := json.Marshal(market.Summary)
	if err != nil {
		return Invalid, -1, -1, err
	}

	fmt.Println(string(JSONContent))
	return DoNothing, 0, 0, nil
}

// SetUpStrategy does nothing for this strategy.
func (w WatchStrategy) SetUpStrategy(wrapper exchangeWrappers.ExchangeWrapper, market environment.Market) error {
	return nil
}

// TearDownStrategy does nothing for this strategy.
func (w WatchStrategy) TearDownStrategy(wrapper exchangeWrappers.ExchangeWrapper, market environment.Market) error {
	return nil
}
