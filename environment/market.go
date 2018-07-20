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

package environment

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

//Ticker provides data incoming from API Tickers, which have little amount of information
//regarding very last updates from a market.
type Ticker struct {
	Ask  decimal.Decimal `json:"Ask"`  //Represents ASK value from the ticker.
	Bid  decimal.Decimal `json:"Bid"`  //Represents BID value from the ticker.
	Last decimal.Decimal `json:"Last"` //Represents LAST trade value from the ticker.
}

//Market represents the environment the bot is trading in.
type Market struct {
	Name           string            `json:"name,required"`            //Represents the name of the market as defined in general (e.g. ETH-BTC).
	BaseCurrency   string            `json:"baseCurrency,omitempty"`   //Represents the base currency of the market.
	MarketCurrency string            `json:"marketCurrency,omitempty"` //Represents the currency to exchange by using base currency.
	ExchangeNames  map[string]string `json:"-"`                        // Represents the various names of the market on various exchanges.
}

func (m Market) String() string {
	ret := fmt.Sprintln("Market", m.Name)
	return strings.TrimSpace(ret)
}

//MarketSummary represents the summary data of a market.
type MarketSummary struct {
	High   decimal.Decimal `json:"high,required"`   //Represents the 24 hours maximum peak of this market.
	Low    decimal.Decimal `json:"low,required"`    //Represents the 24 hours minimum peak of this market.
	Volume decimal.Decimal `json:"volume,required"` //Represents the 24 volume peak of this market.
	Ask    decimal.Decimal `json:"ask,required"`    //Represents the current ASK value.
	Bid    decimal.Decimal `json:"bid,required"`    //Represents the current BID value.
	Last   decimal.Decimal `json:"last,required"`   //Represents the value of the last trade.
}

func (ms MarketSummary) String() string {
	return fmt.Sprintln("  Last: ", ms.Last) +
		fmt.Sprintln("  ASK: ", ms.Ask) +
		fmt.Sprintln("  BID: ", ms.Bid) +
		fmt.Sprintln("  Volume: ", ms.Volume) +
		fmt.Sprintln("  High: ", ms.High) +
		fmt.Sprintln("  Low: ", ms.Low)
}

//UpdateFromTicker updates the values of the market summary from a Ticker Data.
func (ms *MarketSummary) UpdateFromTicker(ticker Ticker) {
	ms.Ask = ticker.Ask
	ms.Bid = ticker.Bid
	ms.Last = ticker.Last
}
