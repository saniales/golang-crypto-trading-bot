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
