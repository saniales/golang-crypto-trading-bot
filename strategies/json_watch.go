package strategies

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AlessandroSanino1994/golang-crypto-trading-bot/environment"
	"github.com/AlessandroSanino1994/golang-crypto-trading-bot/exchangeWrappers"
)

func init() {
	AddStrategy(&JSONWatchStrategy{
		Label:           "JSONWatchEveryHour",
		RefreshInterval: time.Hour,
	})
	AddStrategy(&JSONWatchStrategy{
		Label:           "JSONWatchEvery30Minutes",
		RefreshInterval: time.Minute * 30,
	})
	AddStrategy(&JSONWatchStrategy{
		Label:           "JSONWatchEvery5Minutes",
		RefreshInterval: time.Second * 5,
	})
	AddStrategy(&JSONWatchStrategy{
		Label:           "JSONTestWatch",
		RefreshInterval: time.Second,
	})
}

// JSONWatchStrategy represents a strategy which does nothing than print
// what it gets from markets.
type JSONWatchStrategy struct {
	Label               string        // The Label used to name the strategy.
	RefreshInterval     time.Duration // Interval represents how often summaries will be requested.
	skipFirstCycleDelay bool          // Tells if the strategy has just been initialized and is started. If true the first cycle will be delayed of time expressed by RefreshInterval variable.
}

// Name returns the name of the strategy.
func (w JSONWatchStrategy) Name() string {
	if w.Label == "" {
		return "Watch"
	}
	return w.Label
}

// OnCandleUpdate Prints info about markets on every candle tick in JSON format.
func (w *JSONWatchStrategy) OnCandleUpdate(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) (Action, float64, float64, error) {
	if w.skipFirstCycleDelay {
		time.Sleep(w.RefreshInterval)
	}
	w.skipFirstCycleDelay = true

	err := wrapper.GetTicker(market)
	if err != nil {
		return Invalid, -1, -1, err
	}

	JSONContent, err := json.Marshal(market)
	if err != nil {
		return Invalid, -1, -1, err
	}

	fmt.Println(string(JSONContent))
	return DoNothing, 0, 0, nil
}

// SetUpStrategy does nothing for this strategy.
func (w JSONWatchStrategy) SetUpStrategy(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) error {
	return nil
}

// TearDownStrategy does nothing for this strategy.
func (w JSONWatchStrategy) TearDownStrategy(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) error {
	return nil
}
