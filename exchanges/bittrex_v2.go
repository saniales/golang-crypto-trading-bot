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

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/thebotguys/golang-bittrex-api/bittrex"
)

// BittrexWrapperV2 wraps Bittrex API v2.0
type BittrexWrapperV2 struct {
	PublicKey string
	SecretKey string
}

// NewBittrexV2Wrapper creates a generic wrapper of the bittrex API v2.0.
func NewBittrexV2Wrapper(publicKey string, secretKey string) ExchangeWrapper {
	return BittrexWrapperV2{
		PublicKey: publicKey,
		SecretKey: secretKey,
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper BittrexWrapperV2) Name() string {
	return "bittrex"
}

// GetMarkets gets all the markets info.
func (wrapper BittrexWrapperV2) GetMarkets() ([]*environment.Market, error) {
	bittrexMarkets, err := bittrex.GetMarkets()
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
func (wrapper BittrexWrapperV2) GetOrderBook(market *environment.Market) error {
	return errors.New("GetOrderBook not implemented")
	/*
		bittrexOrderBook, err := wrapper.bittrexAPI.GetOrderBook(market.Name, "both", 100)
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
		totalLength := len(bittrexOrderBook.Buy) + len(bittrexOrderBook.Sell)
		orders := make([]environment.Order, totalLength)
		for i, order := range bittrexOrderBook.Buy {
			orders[i] = convertFromBittrexOrder(environment.Bid, order)
		}
		for i, order := range bittrexOrderBook.Sell {
			orders[i+len(bittrexOrderBook.Buy)] = convertFromBittrexOrder(environment.Ask, order)
		}

		return nil
	*/
}

// BuyLimit performs a limit buy action.
func (wrapper BittrexWrapperV2) BuyLimit(market environment.Market, amount float64, limit float64) (string, error) {
	return "", errors.New("BuyLimit not implemented")
	/*
		orderNumber, err := wrapper.bittrexAPI.BuyLimit(market.Name, amount, limit)
		return orderNumber, err
	*/
}

// BuyMarket performs a market buy action.
func (wrapper BittrexWrapperV2) BuyMarket(market environment.Market, amount float64) (string, error) {
	return "", errors.New("BuyMarket not implemented")
	/*
		orderNumber, err := wrapper.bittrexAPI.BuyMarket(market.Name, amount)
		return orderNumber, err
	*/
}

// SellLimit performs a limit sell action.
func (wrapper BittrexWrapperV2) SellLimit(market environment.Market, amount float64, limit float64) (string, error) {
	return "", errors.New("SellLimit not implemented")
	/*
		orderNumber, err := wrapper.bittrexAPI.SellLimit(market.Name, amount, limit)
		return orderNumber, err
	*/
}

// SellMarket performs a market sell action.
func (wrapper BittrexWrapperV2) SellMarket(market environment.Market, amount float64) (string, error) {
	return "", errors.New("SellMarket not implemented")
	/*
		orderNumber, err := wrapper.bittrexAPI.SellMarket(market.Name, amount)
		return orderNumber, err
	*/
}

// GetTicker gets the updated ticker for a market.
func (wrapper BittrexWrapperV2) GetTicker(market *environment.Market) error {
	return errors.New("GetTicker not implemented")
	/*
		bittrexTicker, err := wrapper.bittrexAPI.GetTicker(market.Name)
		if err != nil {
			return err
		}
		market.Summary.UpdateFromTicker(convertFromBittrexTicker(bittrexTicker))
		return nil
	*/
}

// GetMarketSummaries get the markets summary of all markets
func (wrapper BittrexWrapperV2) GetMarketSummaries(markets map[string]*environment.Market) error {
	bittrexSummaries, err := bittrex.GetMarketSummaries()
	if err != nil {
		return err
	}
	for _, summary := range bittrexSummaries {
		if _, exists := markets[summary.MarketName]; !exists {
			markets[summary.MarketName] = &environment.Market{
				Name: summary.MarketName,
			}
		}
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
}

// GetMarketSummary gets the current market summary.
func (wrapper BittrexWrapperV2) GetMarketSummary(market *environment.Market) error {
	summary, err := bittrex.GetMarketSummary(market.Name)
	if err != nil {
		return err
	}

	market.Summary = environment.MarketSummary{
		High:   summary.High,
		Low:    summary.Low,
		Volume: summary.Volume,
		Bid:    summary.Bid,
		Ask:    summary.Ask,
		Last:   summary.Last,
	}
	return nil
}
