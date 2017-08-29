package strategies

import (
	"fmt"
	"time"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchangeWrappers"
)

// Watch5Min prints out the info of the market every 5 minutes.
var Watch5Min Strategy = IntervalStrategy{
	model: StrategyModel{
		Name: "Watch5Min",
		Setup: func(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) error {
			fmt.Println("Watch5Min starting")
			return nil
		},
		OnUpdate: func(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) error {
			err := wrapper.GetTicker(market)
			if err != nil {
				return err
			}
			fmt.Println(market)
			return nil
		},
		OnError: func(err error) {
			fmt.Println(err)
		},
		TearDown: func(wrapper exchangeWrappers.ExchangeWrapper, market *environment.Market) error {
			fmt.Println("Watch5Min exited")
			return nil
		},
	},
	Interval: time.Minute * 5,
}
