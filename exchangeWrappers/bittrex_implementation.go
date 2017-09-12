package exchangeWrappers

import (
	"github.com/saniales/golang-crypto-trading-bot/environment"

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
func (wrapper BittrexWrapper) GetOrderBook(market *environment.Market) error {
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
		orders[i] = environment.Order{
			Type:     environment.Bid,
			Quantity: order.Quantity,
			Value:    order.Rate,
		}
	}
	for i, order := range bittrexOrderBook.Sell {
		orders[i+len(bittrexOrderBook.Buy)] = environment.Order{
			Type:     environment.Ask,
			Quantity: order.Quantity,
			Value:    order.Rate,
		}
	}

	return nil
}

// BuyLimit performs a limit buy action.
func (wrapper BittrexWrapper) BuyLimit(market environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.bittrexAPI.BuyLimit(market.Name, amount, limit)
	return orderNumber, err
}

// BuyMarket performs a market buy action.
func (wrapper BittrexWrapper) BuyMarket(market environment.Market, amount float64) (string, error) {
	orderNumber, err := wrapper.bittrexAPI.BuyMarket(market.Name, amount)
	return orderNumber, err
}

// SellLimit performs a limit sell action.
func (wrapper BittrexWrapper) SellLimit(market environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.bittrexAPI.SellLimit(market.Name, amount, limit)
	return orderNumber, err
}

// SellMarket performs a market sell action.
func (wrapper BittrexWrapper) SellMarket(market environment.Market, amount float64) (string, error) {
	orderNumber, err := wrapper.bittrexAPI.SellMarket(market.Name, amount)
	return orderNumber, err
}

// GetTicker gets the updated ticker for a market.
func (wrapper BittrexWrapper) GetTicker(market *environment.Market) error {
	bittrexTicker, err := wrapper.bittrexAPI.GetTicker(market.Name)
	if err != nil {
		return err
	}
	market.Summary.UpdateFromTicker(environment.Ticker{
		Last: bittrexTicker.Last,
		Bid:  bittrexTicker.Bid,
		Ask:  bittrexTicker.Ask,
	})
	return nil
}

// GetMarketSummaries get the markets summary of all markets
func (wrapper BittrexWrapper) GetMarketSummaries(markets map[string]*environment.Market) error {
	bittrexSummaries, err := wrapper.bittrexAPI.GetMarketSummaries()
	if err != nil {
		return err
	}
	for _, summary := range bittrexSummaries {
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
func (wrapper BittrexWrapper) GetMarketSummary(market *environment.Market) error {
	summaryArray, err := wrapper.bittrexAPI.GetMarketSummary(market.Name)
	if err != nil {
		return err
	}
	summary := summaryArray[0]

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

//convertFromBittrexCandle converts a bittrex candle to a environment.CandleStick.
func convertFromBittrexCandle(candle bittrexAPI.Candle) environment.CandleStick {
	return environment.CandleStick{
		High:  candle.High,
		Open:  candle.Open,
		Close: candle.Close,
		Low:   candle.Low,
	}
}
