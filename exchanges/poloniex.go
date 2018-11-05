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
	api              *poloniex.Poloniex // access to Poloniex API
	bindedTickers    map[string]bool    // if true, i am subscribing to market ticker.
	summaries        *SummaryCache
	candles          *CandlesCache
	depositAddresses map[string]string
	websocketOn      bool
}

// NewPoloniexWrapper creates a generic wrapper of the poloniex API.
func NewPoloniexWrapper(publicKey string, secretKey string, depositAddresses map[string]string) ExchangeWrapper {
	return &PoloniexWrapper{
		api:              poloniex.NewWithCredentials(publicKey, secretKey),
		bindedTickers:    make(map[string]bool),
		summaries:        NewSummaryCache(),
		candles:          NewCandlesCache(),
		depositAddresses: depositAddresses,
		websocketOn:      false,
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper *PoloniexWrapper) Name() string {
	return "poloniex"
}

func (wrapper *PoloniexWrapper) String() string {
	return wrapper.Name()
}

// GetMarkets gets all the markets info.
func (wrapper *PoloniexWrapper) GetMarkets() ([]*environment.Market, error) {
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

// GetCandles gets the candle data from the exchange.
func (wrapper *PoloniexWrapper) GetCandles(market *environment.Market) ([]environment.CandleStick, error) {
	if !wrapper.websocketOn {
		poloniesCandles, err := wrapper.api.ChartData(MarketNameFor(market, wrapper))
		if err != nil {
			return nil, err
		}

		ret := make([]environment.CandleStick, len(poloniesCandles))

		for i, poloniexCandle := range poloniesCandles {
			ret[i] = environment.CandleStick{
				High:   decimal.NewFromFloat(poloniexCandle.High),
				Open:   decimal.NewFromFloat(poloniexCandle.Open),
				Close:  decimal.NewFromFloat(poloniexCandle.Close),
				Low:    decimal.NewFromFloat(poloniexCandle.Low),
				Volume: decimal.NewFromFloat(poloniexCandle.Volume),
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

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper *PoloniexWrapper) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	poloniexOrderBook, err := wrapper.api.OrderBook(MarketNameFor(market, wrapper))
	if err != nil {
		return nil, err
	}

	var orderBook environment.OrderBook
	for _, order := range poloniexOrderBook.Bids {
		orderBook.Bids = append(orderBook.Bids, environment.Order{
			Quantity: decimal.NewFromFloat(order.Amount),
			Value:    decimal.NewFromFloat(order.Rate),
		})
	}
	for _, order := range poloniexOrderBook.Asks {
		orderBook.Asks = append(orderBook.Asks, environment.Order{
			Quantity: decimal.NewFromFloat(order.Amount),
			Value:    decimal.NewFromFloat(order.Rate),
		})
	}

	return &orderBook, nil
}

// BuyLimit performs a limit buy action.
func (wrapper *PoloniexWrapper) BuyLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.Buy(MarketNameFor(market, wrapper), amount, limit)
	return fmt.Sprint(orderNumber.OrderNumber), err
}

// SellLimit performs a limit sell action.
func (wrapper *PoloniexWrapper) SellLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.Sell(MarketNameFor(market, wrapper), amount, limit)
	return fmt.Sprint(orderNumber.OrderNumber), err
}

// BuyMarket performs a market buy action.
func (wrapper *PoloniexWrapper) BuyMarket(market *environment.Market, amount float64) (string, error) {
	panic("Not supported on poloniex")
}

// SellMarket performs a market sell action.
func (wrapper *PoloniexWrapper) SellMarket(market *environment.Market, amount float64) (string, error) {
	panic("Not supported on poloniex")
}

// GetTicker gets the updated ticker for a market.
func (wrapper *PoloniexWrapper) GetTicker(market *environment.Market) (*environment.Ticker, error) {
	poloniexTicker, err := wrapper.api.Ticker()
	if err != nil {
		return nil, err
	}
	ticker, exists := poloniexTicker[MarketNameFor(market, wrapper)]
	if !exists {
		return nil, errors.New("Market not found")
	}

	return &environment.Ticker{
		Last: decimal.NewFromFloat(ticker.Last),
		Bid:  decimal.NewFromFloat(ticker.Bid),
		Ask:  decimal.NewFromFloat(ticker.Ask),
	}, nil
}

// GetMarketSummary gets the current market summary.
func (wrapper *PoloniexWrapper) GetMarketSummary(market *environment.Market) (*environment.MarketSummary, error) {
	if !wrapper.websocketOn {
		poloniexSummaries, err := wrapper.api.Ticker()
		if err != nil {
			return nil, err
		}

		for pair, poloniexSummary := range poloniexSummaries {
			if pair == MarketNameFor(market, wrapper) {
				wrapper.summaries.Set(market, &environment.MarketSummary{
					Ask:    decimal.NewFromFloat(poloniexSummary.Ask),
					Bid:    decimal.NewFromFloat(poloniexSummary.Bid),
					Last:   decimal.NewFromFloat(poloniexSummary.Last),
					Volume: decimal.NewFromFloat(poloniexSummary.BaseVolume),
				})
				break
			}
		}
	}

	ret, exists := wrapper.summaries.Get(market)
	if !exists {
		return nil, errors.New("Market not found")
	}

	return ret, nil
}

// GetBalance gets the balance of the user of the specified currency.
func (wrapper *PoloniexWrapper) GetBalance(symbol string) (*decimal.Decimal, error) {
	poloniexBalances, err := wrapper.api.Balances()
	if err != nil {
		return nil, err
	}

	for asset, poloniexBalance := range poloniexBalances {
		if asset == symbol {
			ret := decimal.NewFromFloat(poloniexBalance.Available)
			return &ret, nil
		}
	}

	return nil, errors.New("Symbol not found")
}

// GetDepositAddress gets the deposit address for the specified coin on the exchange.
func (wrapper *PoloniexWrapper) GetDepositAddress(coinTicker string) (string, bool) {
	addr, exists := wrapper.depositAddresses[coinTicker]
	return addr, exists
}

// CalculateTradingFees calculates the trading fees for an order on a specified market.
//
//     NOTE: In Binance fees are currently hardcoded.
func (wrapper *PoloniexWrapper) CalculateTradingFees(market *environment.Market, amount float64, limit float64, orderType TradeType) float64 {
	// NOTE: possibility to use wrapper FeesInfo function.
	var feePercentage float64
	if orderType == MakerTrade {
		feePercentage = 0.0010
	} else if orderType == TakerTrade {
		feePercentage = 0.0020
	} else {
		panic("Unknown trade type")
	}

	return amount * limit * feePercentage
}

// CalculateWithdrawFees calculates the withdrawal fees on a specified market.
func (wrapper *PoloniexWrapper) CalculateWithdrawFees(market *environment.Market, amount float64) float64 {
	panic("Not Implemented")
}

// FeedConnect connects to the feed of the poloniex websocket.
func (wrapper *PoloniexWrapper) FeedConnect(markets []*environment.Market) error {
	wrapper.api.StartWS()
	wrapper.websocketOn = true

	for _, m := range markets {
		wrapper.subscribeMarketSummaryFeed(m)
	}

	return nil
}

// SubscribeMarketSummaryFeed subscribes to the Market Summary Feed service.
func (wrapper *PoloniexWrapper) subscribeMarketSummaryFeed(market *environment.Market) {
	if wrapper.websocketOn {
		subTicker := fmt.Sprintf("ticker:%s", MarketNameFor(market, wrapper))
		if len(wrapper.bindedTickers) == 0 {
			wrapper.api.Subscribe("ticker")

			wrapper.api.On("ticker", func(t poloniex.WSTicker) {
				for bindedTicker := range wrapper.bindedTickers {
					if bindedTicker == t.Pair {
						wrapper.api.Emit(subTicker, t)
					}
				}
			})
		}

		if _, exists := wrapper.bindedTickers[MarketNameFor(market, wrapper)]; !exists {
			wrapper.bindedTickers[MarketNameFor(market, wrapper)] = true

			wrapper.api.On(subTicker, func(t poloniex.WSTicker) {
				wrapper.summaries.Set(market, &environment.MarketSummary{
					High:   decimal.NewFromFloat(t.DailyHigh),
					Low:    decimal.NewFromFloat(t.DailyLow),
					Last:   decimal.NewFromFloat(t.Last),
					Ask:    decimal.NewFromFloat(t.Ask),
					Bid:    decimal.NewFromFloat(t.Bid),
					Volume: decimal.NewFromFloat(t.BaseVolume),
				})
			})
		}
	}
}

// Withdraw performs a withdraw operation from the exchange to a destination address.
func (wrapper *PoloniexWrapper) Withdraw(destinationAddress string, coinTicker string, amount float64) error {
	_, err := wrapper.api.Withdraw(coinTicker, amount, destinationAddress)
	if err != nil {
		return err
	}

	return nil
}
