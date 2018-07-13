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

package exchanges

import "github.com/saniales/golang-crypto-trading-bot/environment"

//ExchangeWrapper provides a generic wrapper for exchange services.
type ExchangeWrapper interface {
	Name() string // Gets the name of the exchange.
	//DEPRECATED
	//GetCandles(market *environment.Market, interval string) error // Gets the candles of a market.
	//GetMarkets() ([]*environment.Market, error) //Gets all the markets info.
	GetTicker(market *environment.Market) error        //Gets the updated ticker for a market.
	GetMarketSummary(market *environment.Market) error //Gets the current market summary.
	//GetMarketSummaries(markets map[string]*environment.Market) error                    //Gets the current market summaries.
	GetOrderBook(market *environment.Market) error                                       //Gets the order(ASK + BID) book of a market.
	BuyLimit(market *environment.Market, amount float64, limit float64) (string, error)  //performs a limit buy action.
	SellLimit(market *environment.Market, amount float64, limit float64) (string, error) //performs a limit sell action.
}

// MarketNameFor gets the market name as seen by the exchange.
func MarketNameFor(m *environment.Market, wrapper ExchangeWrapper) string {
	return m.ExchangeNames[wrapper.Name()]
}
