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

package exchangeWrappers

import (
	"context"
	"errors"
	"fmt"

	"github.com/adshao/go-binance"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/shopspring/decimal"
)

// BinanceWrapper represents the wrapper for the Binance exchange.
type BinanceWrapper struct {
	api *binance.Client
}

// NewBinanceWrapper creates a generic wrapper of the binance API.
func NewBinanceWrapper(publicKey string, secretKey string) ExchangeWrapper {
	client := binance.NewClient(publicKey, secretKey)
	return BinanceWrapper{
		api: client,
	}
}

// GetMarkets Gets all the markets info.
func (wrapper BinanceWrapper) GetMarkets() ([]*environment.Market, error) {
	binanceMarkets, err := wrapper.api.NewListPricesService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	ret := make([]*environment.Market, len(binanceMarkets))

	for i, market := range binanceMarkets {
		if len(market.Symbol) == 6 {
			quote := market.Symbol[0:2]
			base := market.Symbol[3:5]
			ret[i] = &environment.Market{
				Name:           market.Symbol,
				BaseCurrency:   base,
				MarketCurrency: quote,
			}
		} else {
			panic("Handle this case")
		}
	}

	return ret, nil
}

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper BinanceWrapper) GetOrderBook(market *environment.Market) error {
	binanceOrderBook, err := wrapper.api.NewListOrdersService().Symbol(market.Name).Do(context.Background())
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

	totalLength := len(binanceOrderBook)
	orders := make([]environment.Order, totalLength)
	for i, order := range binanceOrderBook {
		qty, err := decimal.NewFromString(order.ExecutedQuantity)
		if err != nil {
			return err
		}

		value, err := decimal.NewFromString(order.Price)
		if err != nil {
			return err
		}

		if order.Type == "ASK" {
			orders[i] = environment.Order{
				Type:     environment.Ask,
				Quantity: qty,
				Value:    value,
			}
		} else if order.Type == "BID" {
			orders[i] = environment.Order{
				Type:     environment.Bid,
				Quantity: qty,
				Value:    value,
			}
		}
	}

	market.WatchedChart.OrderBook = orders
	return nil
}

// BuyLimit performs a limit buy action.
func (wrapper BinanceWrapper) BuyLimit(market environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.NewCreateOrderService().Type(binance.OrderTypeLimit).Side(binance.SideTypeBuy).Symbol(market.Name).Price(fmt.Sprint(limit)).Quantity(fmt.Sprint(amount)).Do(context.Background())
	return fmt.Sprint(orderNumber.ClientOrderID), err
}

// SellLimit performs a limit sell action.
func (wrapper BinanceWrapper) SellLimit(market environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.NewCreateOrderService().Type(binance.OrderTypeLimit).Side(binance.SideTypeSell).Symbol(market.Name).Price(fmt.Sprint(limit)).Quantity(fmt.Sprint(amount)).Do(context.Background())
	return fmt.Sprint(orderNumber.ClientOrderID), err
}

// GetTicker gets the updated ticker for a market.
func (wrapper BinanceWrapper) GetTicker(market *environment.Market) error {
	binanceTicker, err := wrapper.api.NewBookTickerService().Symbol(market.Name).Do(context.Background())
	if err != nil {
		return err
	}

	ask, _ := decimal.NewFromString(binanceTicker.AskPrice)
	bid, _ := decimal.NewFromString(binanceTicker.BidPrice)

	market.Summary.UpdateFromTicker(environment.Ticker{
		Last: ask,
		Ask:  ask,
		Bid:  bid,
	})

	return nil
}

// GetMarketSummaries get the markets summary of all markets
func (wrapper BinanceWrapper) GetMarketSummaries(markets map[string]*environment.Market) error {
	binanceSummaries, err := wrapper.api.NewListPriceChangeStatsService().Do(context.Background())
	if err != nil {
		return err
	}

	for _, binanceSummary := range binanceSummaries {
		go func() { // we can assume to have never same symbol twice so its safe
			ask, _ := decimal.NewFromString(binanceSummary.AskPrice)
			bid, _ := decimal.NewFromString(binanceSummary.BidPrice)
			high, _ := decimal.NewFromString(binanceSummary.HighPrice)
			low, _ := decimal.NewFromString(binanceSummary.LowPrice)
			volume, _ := decimal.NewFromString(binanceSummary.Volume)
			markets[binanceSummary.Symbol].Summary = environment.MarketSummary{
				Last:   ask,
				Ask:    ask,
				Bid:    bid,
				High:   high,
				Low:    low,
				Volume: volume,
			}
		}()
	}

	return nil
}

// GetMarketSummary gets the current market summary.
func (wrapper BinanceWrapper) GetMarketSummary(market *environment.Market) error {
	hilo, err := wrapper.api.NewListPriceChangeStatsService().Do(context.Background())
	if err != nil {
		return err
	}

	var binanceSummary *binance.PriceChangeStats

	for _, val := range hilo {
		if val.Symbol == market.Name {
			binanceSummary = val
			break
		}
	}

	if binanceSummary == nil {
		return errors.New("Symbol not found")
	}

	ask, _ := decimal.NewFromString(binanceSummary.AskPrice)
	bid, _ := decimal.NewFromString(binanceSummary.BidPrice)
	high, _ := decimal.NewFromString(binanceSummary.HighPrice)
	low, _ := decimal.NewFromString(binanceSummary.LowPrice)
	volume, _ := decimal.NewFromString(binanceSummary.Volume)

	market.Summary = environment.MarketSummary{
		Last:   ask,
		Ask:    ask,
		Bid:    bid,
		High:   high,
		Low:    low,
		Volume: volume,
	}

	return nil
}
