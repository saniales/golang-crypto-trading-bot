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
	"github.com/shopspring/decimal"

	"github.com/toorop/go-bittrex"
	api "github.com/toorop/go-bittrex"
)

//package github.com/toorop/go-bittrex
//refer to https://github.com/toorop/go-bittrex/blob/master/examples/bittrex.go

// BittrexWrapper provides a Generic wrapper of the Bittrex API.
type BittrexWrapper struct {
	api                 *api.Bittrex //Represents the helper of the Bittrex API.
	summaries           *SummaryCache
	candles             *CandlesCache
	websocketOn         bool
	unsubscribeChannels map[*environment.Market]chan bool
	depositAddresses    map[string]string
}

// NewBittrexWrapper creates a generic wrapper of the bittrex API.
func NewBittrexWrapper(publicKey string, secretKey string, depositAddresses map[string]string) ExchangeWrapper {
	return &BittrexWrapper{
		api:              api.New(publicKey, secretKey),
		websocketOn:      false,
		summaries:        NewSummaryCache(),
		candles:          NewCandlesCache(),
		depositAddresses: depositAddresses,
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper *BittrexWrapper) Name() string {
	return "bittrex"
}

func (wrapper *BittrexWrapper) String() string {
	return wrapper.Name()
}

// GetMarkets gets all the markets info.
func (wrapper *BittrexWrapper) GetMarkets() ([]*environment.Market, error) {
	bittrexMarkets, err := wrapper.api.GetMarkets()
	if err != nil {
		return nil, err
	}
	wrappedMarkets := make([]*environment.Market, 0, len(bittrexMarkets))
	for _, market := range bittrexMarkets {
		wrappedMarkets = append(wrappedMarkets, &environment.Market{
			Name:           market.Symbol,
			BaseCurrency:   market.BaseCurrencySymbol,
			MarketCurrency: market.QuoteCurrencySymbol,
		})
	}
	return wrappedMarkets, nil
}

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper *BittrexWrapper) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	bittrexOrderBook, err := wrapper.api.GetOrderBook(MarketNameFor(market, wrapper), 5, "both")
	if err != nil {
		return nil, err
	}

	var orderBook environment.OrderBook
	for _, order := range bittrexOrderBook.Bid {
		orderBook.Bids = append(orderBook.Bids, environment.Order{
			Quantity: order.Quantity,
			Value:    order.Rate,
		})
	}
	for _, order := range bittrexOrderBook.Ask {
		orderBook.Asks = append(orderBook.Asks, environment.Order{
			Quantity: order.Quantity,
			Value:    order.Rate,
		})
	}

	return &orderBook, nil
}

// BuyLimit performs a limit buy action.
func (wrapper *BittrexWrapper) BuyLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.CreateOrder(bittrex.CreateOrderParams{
		Type:         bittrex.LIMIT,
		TimeInForce:  bittrex.GOOD_TIL_CANCELLED,
		MarketSymbol: MarketNameFor(market, wrapper),
		Quantity:     decimal.NewFromFloat(amount),
		Limit:        limit,
		Direction:    bittrex.BUY,
	})

	return orderNumber.ID, err
}

// SellLimit performs a limit sell action.
func (wrapper *BittrexWrapper) SellLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.CreateOrder(bittrex.CreateOrderParams{
		Type:         bittrex.LIMIT,
		TimeInForce:  bittrex.GOOD_TIL_CANCELLED,
		MarketSymbol: MarketNameFor(market, wrapper),
		Quantity:     decimal.NewFromFloat(amount),
		Limit:        limit,
		Direction:    bittrex.SELL,
	})

	return orderNumber.ID, err
}

// BuyMarket performs a market buy action.
func (wrapper *BittrexWrapper) BuyMarket(market *environment.Market, amount float64) (string, error) {
	panic("Not supported on bittrex")
}

// SellMarket performs a market sell action.
func (wrapper *BittrexWrapper) SellMarket(market *environment.Market, amount float64) (string, error) {
	panic("Not supported on bittrex")
}

// GetTicker gets the updated ticker for a market.
func (wrapper *BittrexWrapper) GetTicker(market *environment.Market) (*environment.Ticker, error) {
	bittrexTicker, err := wrapper.api.GetTicker(MarketNameFor(market, wrapper))
	if err != nil {
		return nil, err
	}

	return &environment.Ticker{
		Last: bittrexTicker[0].LastTradeRate,
		Bid:  bittrexTicker[0].BidRate,
		Ask:  bittrexTicker[0].AskRate,
	}, nil
}

// GetMarketSummary gets the current market summary.
func (wrapper *BittrexWrapper) GetMarketSummary(market *environment.Market) (*environment.MarketSummary, error) {
	if !wrapper.websocketOn {
		summary, err := wrapper.api.GetMarketSummary(MarketNameFor(market, wrapper))
		if err != nil {
			return nil, err
		}

		ticker, err := wrapper.GetTicker(market)
		if err != nil {
			return nil, err
		}

		wrapper.summaries.Set(market, &environment.MarketSummary{
			High:   summary.High,
			Low:    summary.Low,
			Volume: summary.Volume,
			Bid:    ticker.Bid,
			Ask:    ticker.Ask,
			Last:   ticker.Last,
		})
	}

	val, exists := wrapper.summaries.Get(market)
	if !exists {
		return nil, errors.New("Summary not yet loaded")
	}

	return val, nil
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
func (wrapper *BittrexWrapper) GetCandles(market *environment.Market) ([]environment.CandleStick, error) {
	panic("Not supported in Bittrex V1")
}

// GetBalance gets the balance of the user of the specified currency.
func (wrapper *BittrexWrapper) GetBalance(symbol string) (*decimal.Decimal, error) {
	balance, err := wrapper.api.GetBalance(symbol)
	if err != nil {
		return nil, err
	}

	return &balance.Available, nil
}

// GetDepositAddress gets the deposit address for the specified coin on the exchange.
func (wrapper *BittrexWrapper) GetDepositAddress(coinTicker string) (string, bool) {
	addr, exists := wrapper.depositAddresses[coinTicker]
	return addr, exists
}

// CalculateTradingFees calculates the trading fees for an order on a specified market.
//
//     NOTE: In Bittrex fees are hardcoded due to the inability to obtain them via API before placing an order.
func (wrapper *BittrexWrapper) CalculateTradingFees(market *environment.Market, amount float64, limit float64, orderType TradeType) float64 {
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
func (wrapper *BittrexWrapper) CalculateWithdrawFees(market *environment.Market, amount float64) float64 {
	panic("Not Implemented")
}

// FeedConnect connects to the feed of the exchange.
func (wrapper *BittrexWrapper) FeedConnect(markets []*environment.Market) error {
	return ErrWebsocketNotSupported
}

// SubscribeMarketSummaryFeed subscribes to the Market Summary Feed service.
//
//     NOTE: Not supported on Bittrex v1 API, use *BittrexWrapperV2.
func (wrapper *BittrexWrapper) subscribeMarketSummaryFeed(market *environment.Market) {
	panic(ErrWebsocketNotSupported)
}

// Withdraw performs a withdraw operation from the exchange to a destination address.
func (wrapper *BittrexWrapper) Withdraw(destinationAddress string, coinTicker string, amount float64) error {
	_, err := wrapper.api.Withdraw(destinationAddress, coinTicker, decimal.NewFromFloat(amount), "golang-crypto-trading-bot")
	if err != nil {
		return err
	}
	return nil
}
