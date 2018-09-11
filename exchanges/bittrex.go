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
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/shopspring/decimal"

	api "github.com/toorop/go-bittrex"
)

//package github.com/toorop/go-bittrex
//refer to https://github.com/toorop/go-bittrex/blob/master/examples/bittrex.go

// BittrexWrapper provides a Generic wrapper of the Bittrex API.
type BittrexWrapper struct {
	api                 *api.Bittrex //Represents the helper of the Bittrex API.
	unsubscribeChannels map[*environment.Market]chan bool
}

// NewBittrexWrapper creates a generic wrapper of the bittrex API.
func NewBittrexWrapper(publicKey string, secretKey string) ExchangeWrapper {
	return BittrexWrapper{
		api: api.New(publicKey, secretKey),
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper BittrexWrapper) Name() string {
	return "bittrex"
}

func (wrapper BittrexWrapper) String() string {
	return wrapper.Name()
}

// GetMarkets gets all the markets info.
func (wrapper BittrexWrapper) GetMarkets() ([]*environment.Market, error) {
	bittrexMarkets, err := wrapper.api.GetMarkets()
	if err != nil {
		return nil, err
	}
	wrappedMarkets := make([]*environment.Market, 0, len(bittrexMarkets))
	for _, market := range bittrexMarkets {
		if market.IsActive {
			wrappedMarkets = append(wrappedMarkets, &environment.Market{
				Name:           market.MarketName,
				BaseCurrency:   market.BaseCurrency,
				MarketCurrency: market.MarketCurrency,
			})
		}
	}
	return wrappedMarkets, nil
}

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper BittrexWrapper) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	bittrexOrderBook, err := wrapper.api.GetOrderBook(MarketNameFor(market, wrapper), "both")
	if err != nil {
		return nil, err
	}

	var orderBook environment.OrderBook
	for _, order := range bittrexOrderBook.Buy {
		orderBook.Bids = append(orderBook.Bids, environment.Order{
			Quantity: order.Quantity,
			Value:    order.Rate,
		})
	}
	for _, order := range bittrexOrderBook.Sell {
		orderBook.Asks = append(orderBook.Asks, environment.Order{
			Quantity: order.Quantity,
			Value:    order.Rate,
		})
	}

	return nil, nil
}

// BuyLimit performs a limit buy action.
func (wrapper BittrexWrapper) BuyLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.BuyLimit(MarketNameFor(market, wrapper), decimal.NewFromFloat(amount), decimal.NewFromFloat(limit))
	return orderNumber, err
}

// SellLimit performs a limit sell action.
func (wrapper BittrexWrapper) SellLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.SellLimit(MarketNameFor(market, wrapper), decimal.NewFromFloat(amount), decimal.NewFromFloat(limit))
	return orderNumber, err
}

// GetTicker gets the updated ticker for a market.
func (wrapper BittrexWrapper) GetTicker(market *environment.Market) (*environment.Ticker, error) {
	bittrexTicker, err := wrapper.api.GetTicker(MarketNameFor(market, wrapper))
	if err != nil {
		return nil, err
	}

	return &environment.Ticker{
		Last: bittrexTicker.Last,
		Bid:  bittrexTicker.Bid,
		Ask:  bittrexTicker.Ask,
	}, nil
}

// GetMarketSummary gets the current market summary.
func (wrapper BittrexWrapper) GetMarketSummary(market *environment.Market) (*environment.MarketSummary, error) {
	summaryArray, err := wrapper.api.GetMarketSummary(MarketNameFor(market, wrapper))
	if err != nil {
		return nil, err
	}
	summary := summaryArray[0]

	return &environment.MarketSummary{
		High:   summary.High,
		Low:    summary.Low,
		Volume: summary.Volume,
		Bid:    summary.Bid,
		Ask:    summary.Ask,
		Last:   summary.Last,
	}, nil
}

//convertFromBittrexCandle converts a bittrex candle to a environment.CandleStick.
func convertFromBittrexCandle(candle api.Candle) environment.CandleStick {
	return environment.CandleStick{
		High:  candle.High,
		Open:  candle.Open,
		Close: candle.Close,
		Low:   candle.Low,
	}
}

// GetCandles gets the candle data from the exchange.
func (wrapper BittrexWrapper) GetCandles(market *environment.Market) ([]environment.CandleStick, error) {
	panic("Not supported in Bittrex V1")
}

// GetBalance gets the balance of the user of the specified currency.
func (wrapper BittrexWrapper) GetBalance(symbol string) (*decimal.Decimal, error) {
	panic("Not Implemented")
}

// CalculateTradingFees calculates the trading fees for an order on a specified market.
//
//     NOTE: In Bittrex fees are hardcoded due to the inability to obtain them via API before placing an order.
func (wrapper BittrexWrapper) CalculateTradingFees(market *environment.Market, amount float64, limit float64, orderType TradeType) float64 {
	var feePercentage float64
	if orderType == MakerTrade {
		feePercentage = 0.0025
	} else if orderType == TakerTrade {
		feePercentage = 0.0025
	} else {
		panic("Unknown trade type")
	}

	return amount * limit * feePercentage
}

// CalculateWithdrawFees calculates the withdrawal fees on a specified market.
func (wrapper BittrexWrapper) CalculateWithdrawFees(market *environment.Market, amount float64) float64 {
	panic("Not Implemented")
}

// FeedConnect connects to the feed of the exchange.
func (wrapper BittrexWrapper) FeedConnect() {

}

// SubscribeMarketSummaryFeed subscribes to the Market Summary Feed service.
//
//     NOTE: Not supported on Bittrex v1 API, use BittrexWrapperV2.
func (wrapper BittrexWrapper) SubscribeMarketSummaryFeed(market *environment.Market) {
	results := make(chan api.ExchangeState)

	wrapper.api.SubscribeExchangeUpdate(MarketNameFor(market, wrapper), results, wrapper.unsubscribeChannels[market])
}

// UnsubscribeMarketSummaryFeed unsubscribes from the Market Summary Feed service.
//
//     NOTE: Not supported on Bittrex v1 API, use BittrexWrapperV2.
func (wrapper BittrexWrapper) UnsubscribeMarketSummaryFeed(market *environment.Market) {
	panic("Not supported on bittrex v1 API")
}
