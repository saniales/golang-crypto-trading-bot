package strategies

import (
	"fmt"
	"sync"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchangeWrappers"
)

var available map[string]Strategy                   //mapped name -> strategy
var appliedTactics map[*environment.Market]Strategy //mapped strategy -> marketName

// Strategy represents a generic strategy.
type Strategy interface {
	Name() string                                                // Name returns the name of the strategy.
	Apply(exchangeWrappers.ExchangeWrapper, *environment.Market) // Apply applies the strategy when called, using the specified wrapper.
}

//StrategyModel represents a strategy model used by strategies.
type StrategyModel struct {
	Name     string
	Setup    func(exchangeWrappers.ExchangeWrapper, *environment.Market) error
	TearDown func(exchangeWrappers.ExchangeWrapper, *environment.Market) error
	OnUpdate func(exchangeWrappers.ExchangeWrapper, *environment.Market) error
	OnError  func(error)
}

func init() {
	available = make(map[string]Strategy)
	appliedTactics = make(map[*environment.Market]Strategy)

	AddCustomStrategy(Watch5Min)
}

// AddCustomStrategy adds a strategy to the available set.
func AddCustomStrategy(s Strategy) {
	available[s.Name()] = s
}

// MatchWithMarket matches a strategy with a market.
func MatchWithMarket(strategyName string, market *environment.Market) error {
	s, exists := available[strategyName]
	if !exists {
		return fmt.Errorf("Strategy %s does not exist, cannot bind to market %s", strategyName, market.Name)
	}
	appliedTactics[market] = s
	return nil
}

// ApplyAllStrategies applies all matched strategies concurrently.
func ApplyAllStrategies(wrapper exchangeWrappers.ExchangeWrapper) {
	var wg sync.WaitGroup
	wg.Add(len(appliedTactics))
	for m, s := range appliedTactics {
		go func(s Strategy, m *environment.Market, wg *sync.WaitGroup) {
			defer wg.Done()
			s.Apply(wrapper, m)
		}(s, m, &wg)
	}
	wg.Wait()
}
