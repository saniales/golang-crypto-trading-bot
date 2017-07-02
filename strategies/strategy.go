package strategies

import (
	"github.com/AlessandroSanino1994/golang-crypto-trading-bot/environment"
	"github.com/AlessandroSanino1994/golang-crypto-trading-bot/exchangeWrappers"
)

//Action provides which action should the bot take with the current configuration.
type Action int8

const (
	//Buy represents a BUY action.
	Buy Action = iota
	//BuyLimit represents a BUY-LIMIT action.
	BuyLimit Action = iota
	//Sell represents a SELL action.
	Sell Action = iota
	//SellLimit represents a SELL-LIMIT action.
	SellLimit Action = iota
	//DoNothing represents a DO-NOTHING action.
	DoNothing Action = iota
	//CancelOrder represents a CANCEL-ORDER action.
	CancelOrder Action = iota
	//Invalid represents an invalid action.
	Invalid Action = iota
)

// All represents all strategies built into the system.
var All map[string]Strategy

func init() {
	All = make(map[string]Strategy)
}

//Strategy represents a strategy to attach a bot on a market.
type Strategy interface {
	Name() string                                                                                                         // Returns the name of the strategy.
	OnCandleUpdate(wrapper exchangeWrappers.ExchangeWrapper, market environment.Market) (Action, float64, float64, error) // OnCandleUpdate represents what to do when new data has been synced.
	SetUpStrategy(wrapper exchangeWrappers.ExchangeWrapper, market environment.Market) error                              // SetUpStrategy represents what to do when strategy is attached.
	TearDownStrategy(wrapper exchangeWrappers.ExchangeWrapper, market environment.Market) error                           // TearDownStrategy represents what to do when strategy is detached.
}
