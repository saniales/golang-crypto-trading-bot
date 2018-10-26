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

	"github.com/dgrr/kucoin-go"
	"github.com/dgrr/kucoin-go/websocket"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/shopspring/decimal"
)

// KucoinWrapper wrapsKucoin
type KucoinWrapper struct {
	api         *kucoin.Kucoin
	ws          *kucoin.WSClient
	websocketOn bool
	summaries   *SummaryCache
	orderbook   *OrderbookCache
}

// NewKucoinWrapper creates a generic wrapper of theKucoin
func NewKucoinWrapper(publicKey string, secretKey string) ExchangeWrapper {
	ws, _ := kucoin.NewWSClient()
	return &KucoinWrapper{
		api:         kucoin.New(publicKey, secretKey),
		ws:          ws,
		websocketOn: false,
		summaries:   NewSummaryCache(),
		orderbook:   NewOrderbookCache(),
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper *KucoinWrapper) Name() string {
	return "kucoin"
}

func (wrapper *KucoinWrapper) String() string {
	return wrapper.Name()
}

// GetMarkets gets all the markets info.
func (wrapper *KucoinWrapper) GetMarkets() ([]*environment.Market, error) {
	KucoinMarkets, err := wrapper.api.GetSymbols()

	if err != nil {
		return nil, err
	}

	wrappedMarkets := make([]*environment.Market, 0, len(KucoinMarkets))
	for _, market := rangeKucoinMarkets {
		wrappedMarkets = append(wrappedMarkets, &environment.Market{
			Name:           market.Symbol,
			BaseCurrency:   market.CoinType,
			MarketCurrency: market.CoinTypePair,
		})
	}

	return wrappedMarkets, nil
}

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper *KucoinWrapper) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	ret, exists := wrapper.orderbook.Get(market)
	if !wrapper.websocketOn {
		kucoinOrderBook, err := wrapper.api.OrdersBook(MarketNameFor(market, wrapper))

		if err != nil {
			return nil, err
		}

		ret = &environment.OrderBook{}
		for _, order := range kucoinOrderBook.BUY {
			amount := order[1]
			rate := order[0]
			ret.Bids = append(ret.Bids, environment.Order{
				Quantity: amount,
				Value:    rate,
			})
		}
		for _, order := range kucoinOrderBook.SELL {
			amount := order[1]
			rate := order[0]
			ret.Asks = append(ret.Asks, environment.Order{
				Quantity: amount,
				Value:    rate,
			})
		}

		wrapper.orderbook.Set(market, ret)
		return ret, nil
	}

	if !exists {
		return nil, errors.New("Orderbook not loaded")
	}

	return ret, nil
}

// BuyLimit performs a limit buy action.
func (wrapper *KucoinWrapper) BuyLimit(market *environment.Market, amount, limit float64) (string, error) {
	orderOid, err := wrapper.api.CreateOrder(MarketNameFor(market, wrapper), "BUY", limit, amount) 
	
	if err != nil {
		return "", err
	}

	return fmt.Sprint(orderOid), nil
}

// BuyMarket performs a market buy action.
func (wrapper *KucoinWrapper) BuyMarket(market *environment.Market, amount float64) (string, error) {
	panic("Not Implemented")
}

// SellLimit performs a limit sell action.
func (wrapper *KucoinWrapper) SellLimit(market *environment.Market, amount, limit float64) (string, error) {
	orderOid, err := wrapper.api.CreateOrder(MarketNameFor(market, wrapper), "SELL", limit, amount) 
	
	if err != nil {
		return "", err
	}

	return fmt.Sprint(orderOid), nil
}

// SellMarket performs a market sell action.
func (wrapper *KucoinWrapper) SellMarket(market *environment.Market, amount float64) (string, error) {
	panic("Not Implemented")
}

// GetTicker gets the updated ticker for a market.
func (wrapper *KucoinWrapper) GetTicker(market *environment.Market) (*environment.Ticker, error) {

	kucoinTicker, err := wrapper.api.GetSymbol(MarketNameFor(market, wrapper))
	if err != nil {
		return nil, err
	}

	ask := decimal.NewFromFloat(kucoinTicker.Sell)
	bid := decimal.NewFromFloat(kucoinTicker.Buy)

	return &environment.Ticker{
		Last: ask,
		Ask:  ask,
		Bid:  bid,
	}, nil
}

// GetMarketSummary gets the current market summary.
func (wrapper *KucoinWrapper) GetMarketSummary(market *environment.Market) (*environment.MarketSummary, error) {
	ret, exists := wrapper.summaries.Get(market)
	if !wrapper.websocketOn {
		kucoinSummary, err := wrapper.api.GetSymbol(MarketNameFor(market, wrapper))
		if err != nil {
			return nil, err
		}

		ask := decimal.NewFromFloat(kucoinSummary.Sell)
		bid := decimal.NewFromFloat(kucoinSummary.Buy)
		high := decimal.NewFromFloat(kucoinSummary.High)
		low := decimal.NewFromFloat(kucoinSummary.Low)
		last := decimal.NewFromFloat(kucoinSummary.LastDealPrice)
		volume := decimal.NewFromFloat(kucoinSummary.VolValue)

		ret = &environment.MarketSummary{
			Last:   last,
			Ask:    ask,
			Bid:    bid,
			High:   high,
			Low:    low,
			Volume: volume,
		}

		wrapper.summaries.Set(market, ret)
		return ret, nil
	}

	if !exists {
		return nil, errors.New("Summary not loaded")
	}

	return ret, nil
}

// GetBalance gets the balance of the user of the specified currency.
func (wrapper *KucoinWrapper) GetBalance(symbol string) (*decimal.Decimal, error) {
	kucoinBalance, err := wrapper.api.GetCoinBalance(symbol)

	if err != nil {
		return nil, errors.New("Symbol not found")
	}

	ret, err := decimal.NewFromFloat(kucoinBalance.Balance)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// CalculateTradingFees calculates the trading fees for an order on a specified market.
func (wrapper *KucoinWrapper) CalculateTradingFees(market *environment.Market, amount float64, limit float64, orderType TradeType) float64 {
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
func (wrapper *KucoinWrapper) CalculateWithdrawFees(market *environment.Market, amount float64) float64 {
	panic("Not Implemented")
}

// GetCandles gets the candle data from the exchange.
func (wrapper *KucoinWrapper) GetCandles(market *environment.Market) ([]environment.CandleStick, error) {
	panic("Not Implemented")
}

// FeedConnect connects to the feed of the exchange.
func (wrapper *KucoinWrapper) FeedConnect(markets []*environment.Market) error {
	wrapper.websocketOn = true
	for _, m := range markets {
		err := wrapper.subscribeFeeds(m)
		if err != nil {
			return err
		}
	}

	return nil
}

// subscribeFeeds subscribes to the Market Summary Feed service.
func (wrapper *KucoinWrapper) subscribeFeeds(market *environment.Market) error {
	handleTicker := func(wrapper *KucoinWrapper, summaryChannel <-chan kucoin.WSNotificationTickerResponse, m *environment.Market) {
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
	}

	handleOrderbook := func(wrapper *KucoinWrapper, bookSnapshotChannel <-chan kucoin.WSNotificationOrderbookSnapshot, bookUpdateChannel <-chan kucoin.WSNotificationOrderbookUpdate, m *environment.Market) {
		var currentSequence int64

		for {
			select {
			case snap, stillOpen := <-bookSnapshotChannel:
				if !stillOpen {
					return
				}
				if currentSequence > snap.Sequence { // my snapshot is more recent than the one provided
					continue
				}

				orderbook := new(environment.OrderBook)

				for _, item := range snap.Ask {
					price, _ := decimal.NewFromString(item.Price)
					size, _ := decimal.NewFromString(item.Size)

					orderbook.Asks = append(orderbook.Asks, environment.Order{
						Value:    price,
						Quantity: size,
					})
				}
				for _, item := range snap.Bid {
					price, _ := decimal.NewFromString(item.Price)
					size, _ := decimal.NewFromString(item.Size)

					orderbook.Bids = append(orderbook.Bids, environment.Order{
						Value:    price,
						Quantity: size,
					})
				}
				wrapper.orderbook.Set(market, orderbook)
			case update, stillOpen := <-bookUpdateChannel:
				if !stillOpen {
					return
				}

				if currentSequence > update.Sequence {
					continue // my snapshot is more recent than the one provided
				}

				orderbook, exists := wrapper.orderbook.Get(m)
				if !exists {
					continue // wait for snapshot
				}

				orderbook.Asks = updateBook(orderbook.Asks, update.Ask, false)
				orderbook.Bids = updateBook(orderbook.Bids, update.Bid, true)

				wrapper.orderbook.Set(market, orderbook)
			}
		}
	}

	summaryChannel, err := wrapper.ws.SubscribeTicker(MarketNameFor(market, wrapper))
	if err != nil {
		return err
	}

	bookUpdateChannel, bookSnapshotChannel, err := wrapper.ws.SubscribeOrderbook(MarketNameFor(market, wrapper))
	if err != nil {
		return err
	}

	go handleTicker(wrapper, summaryChannel, market)
	go handleOrderbook(wrapper, bookSnapshotChannel, bookUpdateChannel, market)
	return nil
}

func updateBook(ordersToUpdate []environment.Order, newOrders []kucoin.WSSubtypeTrade, reverseOrdering bool) []environment.Order {
	N := len(ordersToUpdate)

	for _, item := range newOrders {
		// replace values
		price, _ := decimal.NewFromString(item.Price)
		size, _ := decimal.NewFromString(item.Size)

		newOrder := environment.Order{
			Value:    price,
			Quantity: size,
		}

		i := sort.Search(N, func(i int) bool {
			if reverseOrdering {
				return ordersToUpdate[i].Value.LessThanOrEqual(price)
			}
			return ordersToUpdate[i].Value.GreaterThanOrEqual(price)
		})
		if size.Equals(decimal.Zero) { //remove it
			if i == N-1 {
				ordersToUpdate = ordersToUpdate[:i-1]
				N--
			} else { // i < N - 1
				ordersToUpdate = append(ordersToUpdate[:i], ordersToUpdate[i+1:]...)
				N--
			}
		} else if i == N { // not found, append
			ordersToUpdate = append(ordersToUpdate, newOrder)
			N++
		} else if price.Equals(ordersToUpdate[i].Value) {
			// replace it
			ordersToUpdate[i] = newOrder
		} else if i == 0 { // prepend it
			ordersToUpdate = append([]environment.Order{newOrder}, ordersToUpdate...)
			N++
		} else { // 0 < i < N, so put new order in the middle
			orders := ordersToUpdate[:i-1]
			orders = append(orders, newOrder)
			orders = append(orders, ordersToUpdate[i-1:]...)
			ordersToUpdate = orders
			N++
		}
	}

	return ordersToUpdate
}

// Withdraw performs a withdraw operation from the exchange to a destination address.
func (wrapper *KucoinWrapper) Withdraw(destinationAddress string, coinTicker string, amount float64) error {
	_, err := wrapper.api.CreateWithdrawalApply(coinTicker, destinationAddress, decimal.NewFromFloat(amount))
	if err != nil {
		return err
	}

	return nil	
}
