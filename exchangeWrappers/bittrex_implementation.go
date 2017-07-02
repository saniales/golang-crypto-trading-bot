package exchangeWrappers

import (
	"fmt"

	"github.com/AlessandroSanino1994/gobot/environment"

	bittrexAPI "github.com/toorop/go-bittrex"
)

// GetMarkets gets all the markets info.
func (wrapper BittrexWrapper) GetMarkets() ([]environment.Market, error) {
	bittrexMarkets, err := wrapper.bittrexAPI.GetMarkets()
	if err != nil {
		return nil, err
	}
	wrappedMarkets := make([]environment.Market, 0, len(bittrexMarkets))
	for _, market := range bittrexMarkets {
		if market.IsActive {
			wrappedMarkets = append(wrappedMarkets, convertFromBittrexMarket(market))
		}
	}
	return wrappedMarkets, nil
}

// GetCandles gets the candles of a market.
func (wrapper BittrexWrapper) GetCandles(market *environment.Market, interval string) error {
	bittrexCandles, err := wrapper.bittrexAPI.GetHisCandles(market.Name, interval)
	if err != nil {
		return err
	}
	/*
		candleInterval, err := time.ParseDuration(interval)
		if err != nil {
			return err
		}
	*/
	if market.WatchedChart == nil {
		market.WatchedChart = &environment.CandleStickChart{
			CandleSticks: make([]environment.CandleStick, len(bittrexCandles)),
			//CandlePeriod: candleInterval,
			OrderBook: nil,
		}
	} else {
		market.WatchedChart.CandleSticks = make([]environment.CandleStick, len(bittrexCandles))
		//market.WatchedChart.CandlePeriod = candleInterval
	}
	if err != nil {
		return err
	}
	for i, candle := range bittrexCandles {
		// market.WatchedChart.Volume += candle.Volume
		market.WatchedChart.CandleSticks[i] = convertFromBittrexCandle(candle)
	}
	return nil
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
		orders[i] = convertFromBittrexOrder(environment.Bid, order)
	}
	for i, order := range bittrexOrderBook.Sell {
		orders[i+len(bittrexOrderBook.Buy)] = convertFromBittrexOrder(environment.Ask, order)
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
	market.Summary.UpdateFromTicker(convertFromBittrexTicker(bittrexTicker))
	return nil
}

// GetMarketSummary gets the current market summary.
func (wrapper BittrexWrapper) GetMarketSummary(market *environment.Market) error {
	bittrexSummary, err := wrapper.bittrexAPI.GetMarketSummary(market.Name)
	if err != nil {
		return err
	}
	fmt.Printf("%v", bittrexSummary[0])
	market.Summary = convertFromBittrexMarketSummary(bittrexSummary)
	return nil
}

//package github.com/toorop/go-bittrex
//refer to https://github.com/toorop/go-bittrex/blob/master/examples/bittrex.go

//BittrexWrapper provides a Generic wrapper of the Bittrex API.
type BittrexWrapper struct {
	bittrexAPI *bittrexAPI.Bittrex //Represents the helper of the Bittrex API.
}

//NewBittrexWrapper creates a generic wrapper of the bittrex API.
func NewBittrexWrapper(publicKey string, secretKey string) ExchangeWrapper {
	return BittrexWrapper{
		bittrexAPI: bittrexAPI.New(publicKey, secretKey),
	}
}

//convertFromBittrexMarket converts a bittrex market to a environment.Market.
func convertFromBittrexMarket(market bittrexAPI.Market) environment.Market {
	return environment.Market{
		Name:           market.MarketName,
		BaseCurrency:   market.BaseCurrency,
		MarketCurrency: market.MarketCurrency,
	}
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

//convertFromBittrexOrder converts a bittrex order to a environment.Order.
func convertFromBittrexOrder(typo environment.OrderType, order bittrexAPI.Orderb) environment.Order {
	return environment.Order{
		Type:        typo,
		Quantity:    order.Quantity,
		Value:       order.Rate,
		OrderNumber: "",
	}
}

//convertFromBittrexMarketSummary converts a bittrex Market Summary to a environment.MarketSummary.
func convertFromBittrexMarketSummary(summary []bittrexAPI.MarketSummary) environment.MarketSummary {
	return environment.MarketSummary{
		High:   summary[0].High,
		Low:    summary[0].Low,
		Volume: summary[0].Volume,
		Bid:    summary[0].Bid,
		Ask:    summary[0].Ask,
		Last:   summary[0].Last,
	}
}

//convertFromBittrexTicker converts a bittrex ticker to a environment.Ticker.
func convertFromBittrexTicker(ticker bittrexAPI.Ticker) environment.Ticker {
	return environment.Ticker{
		Last: ticker.Last,
		Bid:  ticker.Bid,
		Ask:  ticker.Ask,
	}
}
