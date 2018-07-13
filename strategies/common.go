// Copyright © 2017 Alessandro Sanino <saninoale@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package strategies

import (
	"fmt"
	"sync"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
)

var available map[string]Strategy                   //mapped name -> strategy
var appliedTactics map[*environment.Market]Strategy //mapped strategy -> marketName

// Strategy represents a generic strategy.
type Strategy interface {
	Name() string                                           // Name returns the name of the strategy.
	Apply([]exchanges.ExchangeWrapper, *environment.Market) // Apply applies the strategy when called, using the specified wrapper.
}

//StrategyModel represents a strategy model used by strategies.
type StrategyModel struct {
	Name     string
	Setup    func([]exchanges.ExchangeWrapper, *environment.Market) error
	TearDown func([]exchanges.ExchangeWrapper, *environment.Market) error
	OnUpdate func([]exchanges.ExchangeWrapper, *environment.Market) error
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
func ApplyAllStrategies(wrappers []exchanges.ExchangeWrapper) {
	var wg sync.WaitGroup
	wg.Add(len(appliedTactics))
	for m, s := range appliedTactics {
		go func(s Strategy, m *environment.Market, wg *sync.WaitGroup) {
			defer wg.Done()
			s.Apply(wrappers, m)
		}(s, m, &wg)
	}
	wg.Wait()
}
