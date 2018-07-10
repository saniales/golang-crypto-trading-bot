package exchanges

import (
	"fmt"
	"sync"

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
func (wrapper BitfinexWrapper) GetOrderBook(market *environment.Market) error {
	bitfinexOrderBook, err := wrapper.api.OrderBook.Get(market.Name, 0, 0, false)
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
	totalLength := len(bitfinexOrderBook.Asks) + len(bitfinexOrderBook.Bids)
	orders := make([]environment.Order, 0, totalLength)
	for _, order := range bitfinexOrderBook.Bids {
		amount, _ := decimal.NewFromString(order.Amount)
		rate, _ := decimal.NewFromString(order.Rate)
		orders = append(orders, environment.Order{
			Type:     environment.Bid,
			Quantity: amount,
			Value:    rate,
		})
	}
	for _, order := range bitfinexOrderBook.Asks {
		amount, _ := decimal.NewFromString(order.Amount)
		rate, _ := decimal.NewFromString(order.Rate)
		orders = append(orders, environment.Order{
			Type:     environment.Ask,
			Quantity: amount,
			Value:    rate,
		})
	}

	return nil
}

// BuyLimit performs a limit buy action.
//
// NOTE: In bitfinex buy and sell orders behave the same (the go bitfinex api automatically puts it on correct side)
func (wrapper BitfinexWrapper) BuyLimit(market environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.api.Orders.Create(market.Name, amount, limit, bitfinex.OrderTypeExchangeLimit)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(orderNumber.ID), nil
}

// SellLimit performs a limit sell action.
//
// NOTE: In bitfinex buy and sell orders behave the same (the go bitfinex api automatically puts it on correct side)
func (wrapper BitfinexWrapper) SellLimit(market environment.Market, amount float64, limit float64) (string, error) {
	return wrapper.BuyLimit(market, amount, limit)
}

// GetTicker gets the updated ticker for a market.
func (wrapper BitfinexWrapper) GetTicker(market *environment.Market) error {
	bitfinexTicker, err := wrapper.api.Ticker.Get(market.Name)
	if err != nil {
		return err
	}

	last, _ := decimal.NewFromString(bitfinexTicker.LastPrice)
	ask, _ := decimal.NewFromString(bitfinexTicker.Ask)
	bid, _ := decimal.NewFromString(bitfinexTicker.Bid)

	market.Summary.UpdateFromTicker(environment.Ticker{
		Last: last,
		Bid:  bid,
		Ask:  ask,
	})
	return nil
}

// GetMarketSummaries get the markets summary of all markets
//
// WARNING: it panics on error, must be handled by a recover func somewhere
func (wrapper BitfinexWrapper) GetMarketSummaries(markets map[string]*environment.Market) error {
	var wg sync.WaitGroup
	wg.Add(len(markets))
	for _, market := range markets {
		go func(wg *sync.WaitGroup, wrapper ExchangeWrapper, market *environment.Market) {
			err := wrapper.GetMarketSummary(market)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(&wg, wrapper, market)
	}
	wg.Wait()
	return nil
}

// GetMarketSummary gets the current market summary.
func (wrapper BitfinexWrapper) GetMarketSummary(market *environment.Market) error {
	bitfinexSummary, err := wrapper.api.Ticker.Get(market.Name)
	if err != nil {
		return err
	}

	high, _ := decimal.NewFromString(bitfinexSummary.High)
	low, _ := decimal.NewFromString(bitfinexSummary.Low)
	volume, _ := decimal.NewFromString(bitfinexSummary.Volume)
	bid, _ := decimal.NewFromString(bitfinexSummary.Bid)
	ask, _ := decimal.NewFromString(bitfinexSummary.Ask)

	market.Summary = environment.MarketSummary{
		High:   high,
		Low:    low,
		Volume: volume,
		Bid:    bid,
		Ask:    ask,
		Last:   ask, // TODO: find a better way for last value, if any
	}
	return nil
}
