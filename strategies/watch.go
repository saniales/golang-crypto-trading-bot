package strategies

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AlessandroSanino1994/golang-crypto-trading-bot/environment"
	"github.com/AlessandroSanino1994/golang-crypto-trading-bot/exchangeWrappers"
)

func init() {
	AddStrategy(&WatchStrategy{
		Label:           "WatchEveryHour",
		RefreshInterval: time.Hour,
	})
	AddStrategy(&WatchStrategy{
		Label:           "WatchEvery30Minutes",
		RefreshInterval: time.Minute * 30,
	})
	AddStrategy(&WatchStrategy{
		Label:           "WatchEvery5Minutes",
		RefreshInterval: time.Minute * 5,
	})
	AddStrategy(&WatchStrategy{
		Label:           "TestWatch",
		RefreshInterval: time.Minute * 5,
	})
}

// WatchStrategy represents a strategy which does nothing than print
// what it gets from markets.
type WatchStrategy struct {
	Label               string        // The Label used to name the strategy.
	RefreshInterval     time.Duration // Interval represents how often summaries will be requested.
	skipFirstCycleDelay bool          // Tells if the strategy has just been initialized and is started. If true the first cycle will be delayed of time expressed by RefreshInterval variable.
}

// Name returns the name of the strategy.
func (w WatchStrategy) Name() string {
	if w.Label == "" {
		return "Watch"
	}
	return w.Label
}

// OnCandleUpdate Prints info about markets on every candle tick in JSON format.
func (w *WatchStrategy) OnCandleUpdate(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) (Action, float64, float64, error) {
	if w.skipFirstCycleDelay {
		time.Sleep(w.RefreshInterval)
	}
	w.skipFirstCycleDelay = true

	err := wrapper.GetTicker(market)
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
func (w WatchStrategy) SetUpStrategy(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) error {
	return nil
}

// TearDownStrategy does nothing for this strategy.
func (w WatchStrategy) TearDownStrategy(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) error {
	return nil
}
