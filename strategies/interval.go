package strategies

import (
	"time"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchangeWrappers"
)

// IntervalStrategy is an interval based strategy.
type IntervalStrategy struct {
	model    StrategyModel
	Interval time.Duration
}

// Name returns the name of the strategy.
func (is IntervalStrategy) Name() string {
	return is.model.Name
}

// String returns a string representation of the object.
func (is IntervalStrategy) String() string {
	return is.Name()
}

// Apply executes Cyclically the On Update, basing on provided interval.
func (is IntervalStrategy) Apply(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) {
	var err error
	if is.model.Setup != nil {
		err = is.model.Setup(wrapper, market)
		if err != nil {
			is.model.OnError(err)
		}
	}
	for err == nil {
		err = is.model.OnUpdate(wrapper, market)
		if err != nil {
			is.model.OnError(err)
		}
		time.Sleep(is.Interval)
	}
	if is.model.TearDown != nil {
		err = is.model.TearDown(wrapper, market)
		if err != nil {
			is.model.OnError(err)
		}
	}
}
