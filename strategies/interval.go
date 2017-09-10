package strategies

import (
	"time"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchangeWrappers"
)

// IntervalStrategy is an interval based strategy.
type IntervalStrategy struct {
	Model    StrategyModel
	Interval time.Duration
}

// Name returns the name of the strategy.
func (is IntervalStrategy) Name() string {
	return is.Model.Name
}

// String returns a string representation of the object.
func (is IntervalStrategy) String() string {
	return is.Name()
}

// Apply executes Cyclically the On Update, basing on provided interval.
func (is IntervalStrategy) Apply(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) {
	var err error
	if is.Model.Setup != nil {
		err = is.Model.Setup(wrapper, market)
		if err != nil && is.Model.OnError != nil {
			is.Model.OnError(err)
		}
	}
	for err == nil {
		err = is.Model.OnUpdate(wrapper, market)
		if err != nil && is.Model.OnError != nil {
			is.Model.OnError(err)
		}
		time.Sleep(is.Interval)
	}
	if is.Model.TearDown != nil {
		err = is.Model.TearDown(wrapper, market)
		if err != nil && is.Model.OnError != nil {
			is.Model.OnError(err)
		}
	}
}
