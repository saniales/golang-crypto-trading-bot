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
	"fmt"
	"sync"

	"github.com/beldur/kraken-go-api-client"
	"github.com/fatih/structs"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/shopspring/decimal"
)

// NOTE: https://www.kraken.com/help/api

// KrakenWrapper provides a Generic wrapper of the Kraken API.
type KrakenWrapper struct {
	api *krakenapi.KrakenApi
}

// NewKrakenWrapper creates a generic wrapper of the poloniex API.
func NewKrakenWrapper(publicKey string, secretKey string) ExchangeWrapper {
	return KrakenWrapper{
		api: krakenapi.New(publicKey, secretKey),
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper KrakenWrapper) Name() string {
	return "kraken"
}

// GetMarkets gets all the markets info.
func (wrapper KrakenWrapper) GetMarkets() ([]*environment.Market, error) {
	krakenMarkets, err := wrapper.api.AssetPairs()
	if err != nil {
		return nil, err
	}

	markets := structs.Map(krakenMarkets)

	wrappedMarkets := make([]*environment.Market, len(markets))
	i := 0
	for name, pair := range markets {
		p := pair.(krakenapi.AssetPairInfo)
		wrappedMarkets[i] = &environment.Market{
			Name:           name,
			BaseCurrency:   p.Base,
			MarketCurrency: p.Quote,
		}
		i++
	}

	return wrappedMarkets, nil
}

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper KrakenWrapper) GetOrderBook(market *environment.Market) error {
	krakenOrderBook, err := wrapper.api.Depth(market.Name, 0)
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
	totalLength := len(krakenOrderBook.Asks) + len(krakenOrderBook.Bids)
	orders := make([]environment.Order, 0, totalLength)
	for _, order := range krakenOrderBook.Bids {
		amount := decimal.NewFromFloat(order.Amount)
		rate := decimal.NewFromFloat(order.Price)
		orders = append(orders, environment.Order{
			Type:     environment.Bid,
			Quantity: amount,
			Value:    rate,
		})
	}
	for _, order := range krakenOrderBook.Asks {
		amount := decimal.NewFromFloat(order.Amount)
		rate := decimal.NewFromFloat(order.Price)
		orders = append(orders, environment.Order{
			Type:     environment.Ask,
			Quantity: amount,
			Value:    rate,
		})
	}

	return nil
}

// BuyLimit performs a limit buy action.
func (wrapper KrakenWrapper) BuyLimit(market environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.AddOrder(market.Name, "buy", "limit", fmt.Sprint(amount), map[string]string{"price": fmt.Sprint(limit)})
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.TransactionIds), nil
}

// SellLimit performs a limit sell action.
//
// NOTE: In kraken buy and sell orders behave the same (the go kraken api automatically puts it on correct side)
func (wrapper KrakenWrapper) SellLimit(market environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.AddOrder(market.Name, "sell", "limit", fmt.Sprint(amount), map[string]string{"price": fmt.Sprint(limit)})
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.TransactionIds), nil
}

// GetTicker gets the updated ticker for a market.
func (wrapper KrakenWrapper) GetTicker(market *environment.Market) error {
	krakenTicker, err := wrapper.api.Ticker(market.Name)
	if err != nil {
		return err
	}

	ticker := krakenTicker.GetPairTickerInfo(market.Name)

	last, _ := decimal.NewFromString(ticker.Close[0])
	ask, _ := decimal.NewFromString(ticker.Ask[0])
	bid, _ := decimal.NewFromString(ticker.Bid[0])

	market.Summary.UpdateFromTicker(environment.Ticker{
		Last: last,
		Bid:  bid,
		Ask:  ask,
	})
	return nil
}

// GetMarketSummaries get the markets summary of all markets
//
// WARNING: it panics on error, must be handled by a recover func somewhere
func (wrapper KrakenWrapper) GetMarketSummaries(markets map[string]*environment.Market) error {
	var wg sync.WaitGroup
	wg.Add(len(markets))
	for _, market := range markets {
		go func(wg *sync.WaitGroup, wrapper ExchangeWrapper, market *environment.Market) {
			err := wrapper.GetMarketSummary(market)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(&wg, wrapper, market)
	}
	wg.Wait()
	return nil
}

// GetMarketSummary gets the current market summary.
func (wrapper KrakenWrapper) GetMarketSummary(market *environment.Market) error {
	krakenSummary, err := wrapper.api.Ticker(market.Name)
	if err != nil {
		return err
	}

	sum := krakenSummary.GetPairTickerInfo(market.Name)

	high, _ := decimal.NewFromString(sum.High[0])
	low, _ := decimal.NewFromString(sum.Low[0])
	volume, _ := decimal.NewFromString(sum.Volume[0])
	bid, _ := decimal.NewFromString(sum.Bid[0])
	ask, _ := decimal.NewFromString(sum.Ask[0])

	market.Summary = environment.MarketSummary{
		High:   high,
		Low:    low,
		Volume: volume,
		Bid:    bid,
		Ask:    ask,
		Last:   ask, // TODO: find a better way for last value, if any
	}
	return nil
}
