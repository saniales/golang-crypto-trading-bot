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
	"time"

	"github.com/beldur/kraken-go-api-client"
	"github.com/fatih/structs"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/shopspring/decimal"
)

// NOTE: https://www.kraken.com/help/api

// KrakenWrapper provides a Generic wrapper of the Kraken API.
type KrakenWrapper struct {
	api              *krakenapi.KrakenApi
	summaries        *SummaryCache
	candles          *CandlesCache
	depositAddresses map[string]string
	websocketOn      bool
}

// NewKrakenWrapper creates a generic wrapper of the poloniex API.
func NewKrakenWrapper(publicKey string, secretKey string, depositAddresses map[string]string) ExchangeWrapper {
	return &KrakenWrapper{
		api:              krakenapi.New(publicKey, secretKey),
		summaries:        NewSummaryCache(),
		candles:          NewCandlesCache(),
		depositAddresses: depositAddresses,
		websocketOn:      false,
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper *KrakenWrapper) Name() string {
	return "kraken"
}

func (wrapper *KrakenWrapper) String() string {
	return wrapper.Name()
}

// GetMarkets gets all the markets info.
func (wrapper *KrakenWrapper) GetMarkets() ([]*environment.Market, error) {
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
func (wrapper *KrakenWrapper) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	krakenOrderBook, err := wrapper.api.Depth(MarketNameFor(market, wrapper), 0)
	if err != nil {
		return nil, err
	}

	var orderBook environment.OrderBook
	for _, order := range krakenOrderBook.Bids {
		amount := decimal.NewFromFloat(order.Amount)
		rate := decimal.NewFromFloat(order.Price)
		orderBook.Bids = append(orderBook.Bids, environment.Order{
			Quantity:  amount,
			Value:     rate,
			Timestamp: time.Unix(order.Ts, 0),
		})
	}
	for _, order := range krakenOrderBook.Asks {
		amount := decimal.NewFromFloat(order.Amount)
		rate := decimal.NewFromFloat(order.Price)
		orderBook.Asks = append(orderBook.Asks, environment.Order{
			Quantity:  amount,
			Value:     rate,
			Timestamp: time.Unix(order.Ts, 0),
		})
	}

	return &orderBook, nil
}

// BuyLimit performs a limit buy action.
func (wrapper *KrakenWrapper) BuyLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.AddOrder(MarketNameFor(market, wrapper), "buy", "limit", fmt.Sprint(amount), map[string]string{"price": fmt.Sprint(limit)})
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.TransactionIds), nil
}

// SellLimit performs a limit sell action.
//
// NOTE: In kraken buy and sell orders behave the same (the go kraken api automatically puts it on correct side)
func (wrapper *KrakenWrapper) SellLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.AddOrder(MarketNameFor(market, wrapper), "sell", "limit", fmt.Sprint(amount), map[string]string{"price": fmt.Sprint(limit)})
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.TransactionIds), nil
}

// BuyMarket performs a market buy action.
func (wrapper *KrakenWrapper) BuyMarket(market *environment.Market, amount float64) (string, error) {
	orderNumber, err := wrapper.api.AddOrder(MarketNameFor(market, wrapper), "buy", "market", fmt.Sprint(amount), map[string]string{})
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.TransactionIds), nil
}

// SellMarket performs a market sell action.
func (wrapper *KrakenWrapper) SellMarket(market *environment.Market, amount float64) (string, error) {
	orderNumber, err := wrapper.api.AddOrder(MarketNameFor(market, wrapper), "sell", "market", fmt.Sprint(amount), map[string]string{})
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.TransactionIds), nil
}

// GetTicker gets the updated ticker for a market.
func (wrapper *KrakenWrapper) GetTicker(market *environment.Market) (*environment.Ticker, error) {
	krakenTicker, err := wrapper.api.Ticker(MarketNameFor(market, wrapper))
	if err != nil {
		return nil, err
	}

	ticker := krakenTicker.GetPairTickerInfo(MarketNameFor(market, wrapper))

	last, _ := decimal.NewFromString(ticker.Close[0])
	ask, _ := decimal.NewFromString(ticker.Ask[0])
	bid, _ := decimal.NewFromString(ticker.Bid[0])

	return &environment.Ticker{
		Last: last,
		Bid:  bid,
		Ask:  ask,
	}, nil
}

// GetMarketSummary gets the current market summary.
func (wrapper *KrakenWrapper) GetMarketSummary(market *environment.Market) (*environment.MarketSummary, error) {
	krakenSummary, err := wrapper.api.Ticker(MarketNameFor(market, wrapper))
	if err != nil {
		return nil, err
	}

	sum := krakenSummary.GetPairTickerInfo(MarketNameFor(market, wrapper))

	high, _ := decimal.NewFromString(sum.High[0])
	low, _ := decimal.NewFromString(sum.Low[0])
	volume, _ := decimal.NewFromString(sum.Volume[0])
	bid, _ := decimal.NewFromString(sum.Bid[0])
	ask, _ := decimal.NewFromString(sum.Ask[0])

	return &environment.MarketSummary{
		High:   high,
		Low:    low,
		Volume: volume,
		Bid:    bid,
		Ask:    ask,
		Last:   ask, // TODO: find a better way for last value, if any
	}, nil
}

// GetCandles gets the candle data from the exchange.
func (wrapper *KrakenWrapper) GetCandles(market *environment.Market) ([]environment.CandleStick, error) {
	if !wrapper.websocketOn {
		now := time.Now()

		krakenTrades, err := wrapper.api.Trades(MarketNameFor(market, wrapper), now.Add(-time.Hour*24).Unix())
		if err != nil {
			return nil, err
		}

		trades := krakenTrades.Trades

		for lastTradeTime := time.Unix(krakenTrades.Last, 0); lastTradeTime.Before(now); {
			krakenTrades, err = wrapper.api.Trades(MarketNameFor(market, wrapper), now.Add(-time.Hour*24).Unix())
			if err != nil {
				return nil, err
			}

			trades = append(trades, krakenTrades.Trades...)
		}

		ret := make([]environment.CandleStick, 0, 50)

		step := time.Minute * 30
		start := time.Unix(krakenTrades.Trades[0].Time, 0)

		var open = decimal.NewFromFloat(krakenTrades.Trades[0].PriceFloat)
		var high = open
		var low = open
		var close decimal.Decimal

		N := len(trades)
		for i := 1; i < N; i++ {
			currentTrade := trades[i]
			candleTime := time.Unix(currentTrade.Time, 0)
			isLastTrade := i == N-1

			if candleTime.Before(start.Add(step)) || isLastTrade {
				// aggregate candles from trades.
				currentPrice := decimal.NewFromFloat(currentTrade.PriceFloat)
				high = decimal.Max(high, currentPrice)
				low = decimal.Min(low, currentPrice)
			} else {
				// add candle with aggregate data and reset.
				previousTrade := trades[i-1]
				close = decimal.NewFromFloat(previousTrade.PriceFloat)
				ret = append(ret, environment.CandleStick{
					High:  high,
					Open:  open,
					Close: close,
					Low:   low,
				})
				open = decimal.NewFromFloat(currentTrade.PriceFloat)
				high = decimal.NewFromFloat(0)
				low = decimal.NewFromFloat(999999999)
			}
		}

		wrapper.candles.Set(market, ret)
	}

	ret, candleLoaded := wrapper.candles.Get(market)
	if !candleLoaded {
		return nil, errors.New("No candle data yet")
	}

	return ret, nil
}

// GetBalance gets the balance of the user of the specified currency.
func (wrapper *KrakenWrapper) GetBalance(symbol string) (*decimal.Decimal, error) {
	_, err := wrapper.api.Balance()
	if err != nil {
		return nil, err
	}

	panic("Help me integrate this feature!")
}

// GetDepositAddress gets the deposit address for the specified coin on the exchange.
func (wrapper *KrakenWrapper) GetDepositAddress(coinTicker string) (string, bool) {
	addr, exists := wrapper.depositAddresses[coinTicker]
	return addr, exists
}

// CalculateTradingFees calculates the trading fees for an order on a specified market.
//
//     NOTE: In Kraken fees are currently hardcoded.
func (wrapper *KrakenWrapper) CalculateTradingFees(market *environment.Market, amount float64, limit float64, orderType TradeType) float64 {
	var feePercentage float64
	if orderType == MakerTrade {
		feePercentage = 0.0016
	} else if orderType == TakerTrade {
		feePercentage = 0.0026
	} else {
		panic("Unknown trade type")
	}

	return amount * limit * feePercentage
}

// CalculateWithdrawFees calculates the withdrawal fees on a specified market.
func (wrapper *KrakenWrapper) CalculateWithdrawFees(market *environment.Market, amount float64) float64 {
	panic("Not Implemented")
}

// FeedConnect connects to the feed of the exchange.
func (wrapper *KrakenWrapper) FeedConnect(markets []*environment.Market) error {
	return ErrWebsocketNotSupported
}

// SubscribeMarketSummaryFeed subscribes to the Market Summary Feed service.
func (wrapper *KrakenWrapper) subscribeMarketSummaryFeed(market *environment.Market) {
	panic("Websocket Not Supported")
}

// Withdraw performs a withdraw operation from the exchange to a destination address.
func (wrapper *KrakenWrapper) Withdraw(destinationAddress string, coinTicker string, amount float64) error {
	panic("Not Supported")
}
