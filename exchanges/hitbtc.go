// Copyright © 2017 Alessandro Sanino <saninoale@gmail.com>
// Copyright © 2018 Mangrovia Solutions
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

	"github.com/saniales/go-hitbtc"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/shopspring/decimal"
)

// HitBtcWrapperV2 wraps HitBtc API v2.0
type HitBtcWrapperV2 struct {
	api         *hitbtc.HitBtc
	ws          *hitbtc.WSClient
	websocketOn bool
	summaries   *SummaryCache
}

// NewHitBtcV2Wrapper creates a generic wrapper of the HitBtc API v2.0.
func NewHitBtcV2Wrapper(publicKey string, secretKey string) ExchangeWrapper {
	ws, _ := hitbtc.NewWSClient()
	return &HitBtcWrapperV2{
		api:         hitbtc.New(publicKey, secretKey),
		ws:          ws,
		websocketOn: false,
		summaries:   NewSummaryCache(),
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper *HitBtcWrapperV2) Name() string {
	return "hitbtc"
}

func (wrapper *HitBtcWrapperV2) String() string {
	return wrapper.Name()
}

// GetMarkets gets all the markets info.
func (wrapper *HitBtcWrapperV2) GetMarkets() ([]*environment.Market, error) {
	HitBtcMarkets, err := wrapper.api.GetSymbols()

	if err != nil {
		return nil, err
	}

	wrappedMarkets := make([]*environment.Market, 0, len(HitBtcMarkets))
	for _, market := range HitBtcMarkets {
		wrappedMarkets = append(wrappedMarkets, &environment.Market{
			Name:           market.Id,
			BaseCurrency:   market.BaseCurrency,
			MarketCurrency: market.QuoteCurrency,
		})
	}

	return wrappedMarkets, nil
}

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper *HitBtcWrapperV2) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	hitbtcOrderBook, err := wrapper.api.GetOrderbook(MarketNameFor(market, wrapper))

	if err != nil {
		return nil, err
	}

	var orderBook environment.OrderBook
	for _, order := range hitbtcOrderBook.Bid {
		amount := decimal.NewFromFloat(order.Size)
		rate := decimal.NewFromFloat(order.Price)
		orderBook.Bids = append(orderBook.Bids, environment.Order{
			Quantity: amount,
			Value:    rate,
		})
	}
	for _, order := range hitbtcOrderBook.Ask {
		amount := decimal.NewFromFloat(order.Size)
		rate := decimal.NewFromFloat(order.Price)
		orderBook.Asks = append(orderBook.Asks, environment.Order{
			Quantity: amount,
			Value:    rate,
		})
	}

	return &orderBook, nil
}

// BuyLimit performs a limit buy action.
func (wrapper *HitBtcWrapperV2) BuyLimit(market *environment.Market, amount float64, limit float64) (string, error) {

	requestOrder := hitbtc.Order{
		Symbol:   MarketNameFor(market, wrapper),
		Side:     "buy",
		Status:   "new",
		Type:     "limit",
		Quantity: amount,
		Price:    limit,
	}

	orderNumber, err := wrapper.api.PlaceOrder(requestOrder)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.ClientOrderId), nil
}

// BuyMarket performs a market buy action.
func (wrapper *HitBtcWrapperV2) BuyMarket(market *environment.Market, amount float64) (string, error) {
	requestOrder := hitbtc.Order{
		Symbol:   MarketNameFor(market, wrapper),
		Side:     "buy",
		Status:   "new",
		Type:     "market",
		Quantity: amount,
	}

	orderNumber, err := wrapper.api.PlaceOrder(requestOrder)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.ClientOrderId), nil
}

// SellLimit performs a limit sell action.
func (wrapper *HitBtcWrapperV2) SellLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	requestOrder := hitbtc.Order{
		Symbol:   MarketNameFor(market, wrapper),
		Side:     "sell",
		Status:   "new",
		Type:     "limit",
		Quantity: amount,
		Price:    limit,
	}

	orderNumber, err := wrapper.api.PlaceOrder(requestOrder)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.ClientOrderId), nil
}

// SellMarket performs a market sell action.
func (wrapper *HitBtcWrapperV2) SellMarket(market *environment.Market, amount float64) (string, error) {
	requestOrder := hitbtc.Order{
		Symbol:   MarketNameFor(market, wrapper),
		Side:     "sell",
		Status:   "new",
		Type:     "market",
		Quantity: amount,
	}

	orderNumber, err := wrapper.api.PlaceOrder(requestOrder)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.ClientOrderId), nil
}

// GetTicker gets the updated ticker for a market.
func (wrapper *HitBtcWrapperV2) GetTicker(market *environment.Market) (*environment.Ticker, error) {

	hitbtcTicker, err := wrapper.api.GetTicker(MarketNameFor(market, wrapper))
	if err != nil {
		return nil, err
	}

	ask := decimal.NewFromFloat(hitbtcTicker.Ask)
	bid := decimal.NewFromFloat(hitbtcTicker.Bid)

	return &environment.Ticker{
		Last: ask,
		Ask:  ask,
		Bid:  bid,
	}, nil
}

// GetMarketSummary gets the current market summary.
func (wrapper *HitBtcWrapperV2) GetMarketSummary(market *environment.Market) (*environment.MarketSummary, error) {
	ret, exists := wrapper.summaries.Get(market)
	if !exists || !wrapper.websocketOn {
		hilo, err := wrapper.api.GetAllTicker()
		if err != nil {
			return nil, err
		}

		var hitbtcSummary *hitbtc.Ticker

		for _, val := range hilo {
			if val.Symbol == MarketNameFor(market, wrapper) {
				hitbtcSummary = &val
				break
			}
		}

		if hitbtcSummary == nil {
			return nil, errors.New("Symbol not found")
		}

		ask := decimal.NewFromFloat(hitbtcSummary.Ask)
		bid := decimal.NewFromFloat(hitbtcSummary.Bid)
		high := decimal.NewFromFloat(hitbtcSummary.High)
		low := decimal.NewFromFloat(hitbtcSummary.Low)
		last := decimal.NewFromFloat(hitbtcSummary.Last)
		volume := decimal.NewFromFloat(hitbtcSummary.Volume)

		ret = &environment.MarketSummary{
			Last:   last,
			Ask:    ask,
			Bid:    bid,
			High:   high,
			Low:    low,
			Volume: volume,
		}

		wrapper.summaries.Set(market, ret)
	}

	return ret, nil
}

// GetBalance gets the balance of the user of the specified currency.
func (wrapper *HitBtcWrapperV2) GetBalance(symbol string) (*decimal.Decimal, error) {
	Hitbtcbalance, err := wrapper.api.GetBalances()

	if err != nil {
		return nil, err
	}

	for _, hitbtcBalance := range Hitbtcbalance {
		if hitbtcBalance.Currency == symbol {
			ret, err := decimal.NewFromString(hitbtcBalance.Currency)
			if err != nil {
				return nil, err
			}

			return &ret, nil
		}
	}

	return nil, errors.New("Symbol not found")
}

// CalculateTradingFees calculates the trading fees for an order on a specified market.
func (wrapper *HitBtcWrapperV2) CalculateTradingFees(market *environment.Market, amount float64, limit float64, orderType TradeType) float64 {
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
func (wrapper *HitBtcWrapperV2) CalculateWithdrawFees(market *environment.Market, amount float64) float64 {
	panic("Not Implemented")
}

// GetCandles gets the candle data from the exchange.
func (wrapper *HitBtcWrapperV2) GetCandles(market *environment.Market) ([]environment.CandleStick, error) {
	panic("Not Implemented")
}

// FeedConnect connects to the feed of the exchange.
func (wrapper *HitBtcWrapperV2) FeedConnect(markets []*environment.Market) error {
	wrapper.websocketOn = true

	for _, m := range markets {
		err := wrapper.subscribeMarketSummaryFeed(m)
		if err != nil {
			return err
		}
	}

	return nil
}

// SubscribeMarketSummaryFeed subscribes to the Market Summary Feed service.
func (wrapper *HitBtcWrapperV2) subscribeMarketSummaryFeed(market *environment.Market) error {
	summaryChannel, err := wrapper.ws.SubscribeTicker(MarketNameFor(market, wrapper))
	if err != nil {
		return err
	}

	go func(wrapper *HitBtcWrapperV2, summaryChannel <-chan hitbtc.WSNotificationTickerResponse, m *environment.Market) {
		for {
			summary, stillOpen := <-summaryChannel
			if !stillOpen {
				return
			}

			high, _ := decimal.NewFromString(summary.High)
			low, _ := decimal.NewFromString(summary.Low)
			ask, _ := decimal.NewFromString(summary.Ask)
			bid, _ := decimal.NewFromString(summary.Bid)
			last, _ := decimal.NewFromString(summary.Last)
			volume, _ := decimal.NewFromString(summary.Volume)

			sum := &environment.MarketSummary{
				High:   high,
				Low:    low,
				Last:   last,
				Volume: volume,
				Ask:    ask,
				Bid:    bid,
			}

			wrapper.summaries.Set(m, sum)
		}
	}(wrapper, summaryChannel, market)

	return nil
}
