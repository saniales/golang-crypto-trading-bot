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
	"errors"
	"time"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
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
func (is IntervalStrategy) Apply(wrappers []exchanges.ExchangeWrapper, markets []*environment.Market) {
	var err error

	hasSetupFunc := is.Model.Setup != nil
	hasTearDownFunc := is.Model.TearDown != nil
	hasUpdateFunc := is.Model.OnUpdate != nil
	hasErrorFunc := is.Model.OnError != nil

	if hasSetupFunc {
		err = is.Model.Setup(wrappers, markets)
		if err != nil && hasErrorFunc {
			is.Model.OnError(err)
		}
	}

	if !hasUpdateFunc {
		_err := errors.New("OnUpdate func cannot be empty")
		if hasErrorFunc {
			is.Model.OnError(_err)
		} else {
			panic(_err)
		}
	}
	for err == nil {
		err = is.Model.OnUpdate(wrappers, markets)
		if err != nil && hasErrorFunc {
			is.Model.OnError(err)
		}
		time.Sleep(is.Interval)
	}
	if hasTearDownFunc {
		err = is.Model.TearDown(wrappers, markets)
		if err != nil && hasErrorFunc {
			is.Model.OnError(err)
		}
	}
}
