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

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/pharrisee/poloniex-api"
	"github.com/saniales/golang-crypto-trading-bot/environment"
)

// PoloniexWrapper provides a Generic wrapper of the Poloniex API.
type PoloniexWrapper struct {
	api *poloniex.Poloniex // access to Poloniex API
}

// NewPoloniexWrapper creates a generic wrapper of the poloniex API.
func NewPoloniexWrapper(publicKey string, secretKey string) ExchangeWrapper {
	return PoloniexWrapper{
		api: poloniex.NewWithCredentials(publicKey, secretKey),
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper PoloniexWrapper) Name() string {
	return "poloniex"
}

// GetMarkets gets all the markets info.
func (wrapper PoloniexWrapper) GetMarkets() ([]*environment.Market, error) {
	poloniexMarkets, err := wrapper.api.Currencies()
	if err != nil {
		return nil, err
	}
	wrappedMarkets := make([]*environment.Market, 0, len(poloniexMarkets))
	for _, market := range poloniexMarkets {
		if market.Disabled == 1 {
			name := strings.SplitN(market.Name, "/", 2)
			wrappedMarkets = append(wrappedMarkets, &environment.Market{
				Name:           market.Name,
				BaseCurrency:   name[1],
				MarketCurrency: name[0],
			})
		}
	}
	return wrappedMarkets, nil
}

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper PoloniexWrapper) GetOrderBook(market *environment.Market) error {
	poloniexOrderBook, err := wrapper.api.OrderBook(MarketNameFor(market, wrapper))
	if err != nil {
		return err
	}

	if market.WatchedChart == nil {
		market.WatchedChart = &environment.CandleStickChart{
			// MarketName: market.Name,
		}
	} else {
		market.WatchedChart.OrderBook = nil
	}
	totalLength := len(poloniexOrderBook.Asks) + len(poloniexOrderBook.Bids)
	orders := make([]environment.Order, totalLength)
	for i, order := range poloniexOrderBook.Bids {
		orders[i] = environment.Order{
			Type:     environment.Bid,
			Quantity: decimal.NewFromFloat(order.Amount),
			Value:    decimal.NewFromFloat(order.Rate),
		}
	}
	for i, order := range poloniexOrderBook.Asks {
		orders[i+len(poloniexOrderBook.Asks)] = environment.Order{
			Type:     environment.Ask,
			Quantity: decimal.NewFromFloat(order.Amount),
			Value:    decimal.NewFromFloat(order.Rate),
		}
	}

	return nil
}

// BuyLimit performs a limit buy action.
func (wrapper PoloniexWrapper) BuyLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.Buy(MarketNameFor(market, wrapper), amount, limit)
	return fmt.Sprint(orderNumber.OrderNumber), err
}

// SellLimit performs a limit sell action.
func (wrapper PoloniexWrapper) SellLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.Sell(MarketNameFor(market, wrapper), amount, limit)
	return fmt.Sprint(orderNumber.OrderNumber), err
}

// GetTicker gets the updated ticker for a market.
func (wrapper PoloniexWrapper) GetTicker(market *environment.Market) error {
	poloniexTicker, err := wrapper.api.Ticker()
	if err != nil {
		return err
	}
	ticker, exists := poloniexTicker[MarketNameFor(market, wrapper)]
	if !exists {
		return errors.New("Market not found")
	}
	market.Summary.UpdateFromTicker(environment.Ticker{
		Last: decimal.NewFromFloat(ticker.Last),
		Bid:  decimal.NewFromFloat(ticker.Bid),
		Ask:  decimal.NewFromFloat(ticker.Ask),
	})
	return nil
}

// GetMarketSummaries get the markets summary of all markets
func (wrapper PoloniexWrapper) GetMarketSummaries(markets map[string]*environment.Market) error {
	panic("Not implemented")
	/*
		poloniexSummaries, err := wrapper.api. .GetMarketSummaries()
		if err != nil {
			return err
		}
		for _, summary := range poloniexSummaries {
			markets[summary.MarketName].Summary = environment.MarketSummary{
				High:   summary.High,
				Low:    summary.Low,
				Volume: summary.Volume,
				Bid:    summary.Bid,
				Ask:    summary.Ask,
				Last:   summary.Last,
			}
		}
		return nil
	*/
}

// GetMarketSummary gets the current market summary.
func (wrapper PoloniexWrapper) GetMarketSummary(market *environment.Market) error {
	panic("Not implemented")
	/*
		volume, err := wrapper.api.DailyVolume()
		if err != nil {
			return err
		}

		ticker, err := wrapper.api.Ticker()
		if err != nil {
			return err
		}


		market.Summary = environment.MarketSummary{
			High:   ticker[MarketNameFor(market, wrapper)]. ,
			Low:    summary.Low,
			Volume: decimal.NewFromFloat(volume[MarketNameFor(market, wrapper)]),
			Bid:    summary.Bid,
			Ask:    summary.Ask,
			Last:   summary.Last,
		}
	*/
	return nil
}
