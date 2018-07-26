package exchanges

import (
	"fmt"

	"github.com/shopspring/decimal"

	bitfinex "github.com/bitfinexcom/bitfinex-api-go/v1"
	"github.com/saniales/golang-crypto-trading-bot/environment"
)

// BitfinexWrapper provides a Generic wrapper of the Bitfinex API.
type BitfinexWrapper struct {
	api *bitfinex.Client
}

// NewBitfinexWrapper creates a generic wrapper of the bittrex API.
func NewBitfinexWrapper(publicKey string, secretKey string) ExchangeWrapper {
	return BitfinexWrapper{
		api: bitfinex.NewClient().Auth(publicKey, secretKey),
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper BitfinexWrapper) Name() string {
	return "bitfinex"
}

func (wrapper BitfinexWrapper) String() string {
	return wrapper.Name()
}

// GetMarkets gets all the markets info.
func (wrapper BitfinexWrapper) GetMarkets() ([]*environment.Market, error) {
	bitfinexMarkets, err := wrapper.api.Pairs.All()
	if err != nil {
		return nil, err
	}

	wrappedMarkets := make([]*environment.Market, len(bitfinexMarkets))
	for i, pair := range bitfinexMarkets {
		quote, base := pair[0:2], pair[3:5]
		wrappedMarkets[i] = &environment.Market{
			Name:           pair,
			BaseCurrency:   base,
			MarketCurrency: quote,
		}
	}

	return wrappedMarkets, nil
}

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper BitfinexWrapper) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	bitfinexOrderBook, err := wrapper.api.OrderBook.Get(MarketNameFor(market, wrapper), 0, 0, false)
	if err != nil {
		return nil, err
	}

	var orderBook environment.OrderBook
	for _, order := range bitfinexOrderBook.Bids {
		amount, _ := decimal.NewFromString(order.Amount)
		rate, _ := decimal.NewFromString(order.Rate)
		orderBook.Asks = append(orderBook.Asks, environment.Order{
			Quantity: amount,
			Value:    rate,
		})
	}
	for _, order := range bitfinexOrderBook.Asks {
		amount, _ := decimal.NewFromString(order.Amount)
		rate, _ := decimal.NewFromString(order.Rate)
		orderBook.Bids = append(orderBook.Bids, environment.Order{
			Quantity: amount,
			Value:    rate,
		})
	}

	return &orderBook, nil
}

// BuyLimit performs a limit buy action.
//
// NOTE: In bitfinex buy and sell orders behave the same (the go bitfinex api automatically puts it on correct side)
func (wrapper BitfinexWrapper) BuyLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.Orders.Create(MarketNameFor(market, wrapper), amount, limit, bitfinex.OrderTypeExchangeLimit)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.ID), nil
}

// SellLimit performs a limit sell action.
//
// NOTE: In bitfinex buy and sell orders behave the same (the go bitfinex api automatically puts it on correct side)
func (wrapper BitfinexWrapper) SellLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	return wrapper.BuyLimit(market, amount, limit)
}

// GetTicker gets the updated ticker for a market.
func (wrapper BitfinexWrapper) GetTicker(market *environment.Market) (*environment.Ticker, error) {
	bitfinexTicker, err := wrapper.api.Ticker.Get(MarketNameFor(market, wrapper))
	if err != nil {
		return nil, err
	}

	last, _ := decimal.NewFromString(bitfinexTicker.LastPrice)
	ask, _ := decimal.NewFromString(bitfinexTicker.Ask)
	bid, _ := decimal.NewFromString(bitfinexTicker.Bid)

	return &environment.Ticker{
		Last: last,
		Bid:  bid,
		Ask:  ask,
	}, nil
}

// GetMarketSummary gets the current market summary.
func (wrapper BitfinexWrapper) GetMarketSummary(market *environment.Market) (*environment.MarketSummary, error) {
	bitfinexSummary, err := wrapper.api.Ticker.Get(MarketNameFor(market, wrapper))
	if err != nil {
		return nil, err
	}

	high, _ := decimal.NewFromString(bitfinexSummary.High)
	low, _ := decimal.NewFromString(bitfinexSummary.Low)
	volume, _ := decimal.NewFromString(bitfinexSummary.Volume)
	bid, _ := decimal.NewFromString(bitfinexSummary.Bid)
	ask, _ := decimal.NewFromString(bitfinexSummary.Ask)

	return &environment.MarketSummary{
		High:   high,
		Low:    low,
		Volume: volume,
		Bid:    bid,
		Ask:    ask,
		Last:   ask, // TODO: find a better way for last value, if any
	}, nil
}

// CalculateTradingFees calculates the trading fees for an order on a specified market.
//
//     NOTE: In Bitfinex fees are currently hardcoded.
func (wrapper BitfinexWrapper) CalculateTradingFees(market *environment.Market, amount float64, limit float64, orderType TradeType) float64 {
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
func (wrapper BitfinexWrapper) CalculateWithdrawFees(market *environment.Market, amount float64) float64 {
	panic("Not Implemented")
}

func (wrapper BitfinexWrapper) FeedConnect() {
	err := wrapper.api.WebSocket.Connect()
	if err != nil {
		panic(err)
	}

	
}

func (wrapper BitfinexWrapper) SubscribearketSummaryFeed(market *environment.Market, onUpdate func(environment.MarketSummary)) {
	results := make(chan []float64)
	wrapper.api.WebSocket.AddSubscribe("ticker", MarketNameFor(market, wrapper), results)
	for {

	}
}

func (wrapper BitfinexWrapper) UnsubscribeMarketSummaryFeed() {
	wrapper.api.WebSocket.	
}
