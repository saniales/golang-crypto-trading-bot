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

	bittrexAPI "github.com/toorop/go-bittrex"
)

//package github.com/toorop/go-bittrex
//refer to https://github.com/toorop/go-bittrex/blob/master/examples/bittrex.go

// BittrexWrapper provides a Generic wrapper of the Bittrex API.
type BittrexWrapper struct {
	bittrexAPI *bittrexAPI.Bittrex //Represents the helper of the Bittrex API.
}

// NewBittrexWrapper creates a generic wrapper of the bittrex API.
func NewBittrexWrapper(publicKey string, secretKey string) ExchangeWrapper {
	return BittrexWrapper{
		bittrexAPI: bittrexAPI.New(publicKey, secretKey),
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper BittrexWrapper) Name() string {
	return "bittrex"
}

// GetMarkets gets all the markets info.
func (wrapper BittrexWrapper) GetMarkets() ([]*environment.Market, error) {
	bittrexMarkets, err := wrapper.bittrexAPI.GetMarkets()
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
	bittrexOrderBook, err := wrapper.bittrexAPI.GetOrderBook(market.Name, "both")
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
	orderNumber, err := wrapper.bittrexAPI.BuyLimit(market.Name, decimal.NewFromFloat(amount), decimal.NewFromFloat(limit))
	return orderNumber, err
}

// SellLimit performs a limit sell action.
func (wrapper BittrexWrapper) SellLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.bittrexAPI.SellLimit(market.Name, decimal.NewFromFloat(amount), decimal.NewFromFloat(limit))
	return orderNumber, err
}

// GetTicker gets the updated ticker for a market.
func (wrapper BittrexWrapper) GetTicker(market *environment.Market) (*environment.Ticker, error) {
	bittrexTicker, err := wrapper.bittrexAPI.GetTicker(market.Name)
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
	summaryArray, err := wrapper.bittrexAPI.GetMarketSummary(market.Name)
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
func convertFromBittrexCandle(candle bittrexAPI.Candle) environment.CandleStick {
	return environment.CandleStick{
		High:  candle.High,
		Open:  candle.Open,
		Close: candle.Close,
		Low:   candle.Low,
	}
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
