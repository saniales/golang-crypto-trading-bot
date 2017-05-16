package exchangeWrappers

import (
	"time"

	"github.com/AlessandroSanino1994/gobot/environment"
)

//Gets all the markets info.
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

//Gets the candles of a market.
func (wrapper BittrexWrapper) GetCandles(market environment.Market, interval string) error {
	bittrexCandles, err := wrapper.bittrexAPI.GetHisCandles(market.Name, interval)
	candleInterval, err := time.ParseDuration(interval)
	if err != nil {
		return err
	}
	if market.WatchedChart == nil {
		market.WatchedChart = &environment.CandleStickChart{
			MarketName:   market.Name,
			CandleSticks: make([]environment.CandleStick, len(bittrexCandles)),
			CandlePeriod: candleInterval,
			Volume:       0,
			OrderBook:    nil,
		}
	} else {
		market.WatchedChart.Volume = 0
		market.WatchedChart.CandleSticks = make([]environment.CandleStick, len(bittrexCandles))
		market.WatchedChart.CandlePeriod = candleInterval
	}
	if err != nil {
		return err
	}
	for i, candle := range bittrexCandles {
		market.WatchedChart.Volume += candle.Volume
		market.WatchedChart.CandleSticks[i] = convertFromBittrexCandle(candle)
	}
	return nil
}

//Gets the order(ASK + BID) book of a market.
func (wrapper BittrexWrapper) GetOrderBook(market environment.Market) error {
	bittrexOrderBook, err := wrapper.bittrexAPI.GetOrderBook(market.Name, "both", 100)
	if err != nil {
		return err
	}

	if market.WatchedChart == nil {
		market.WatchedChart = &environment.CandleStickChart{
			MarketName: market.Name,
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

//performs a limit buy action.
func (wrapper BittrexWrapper) BuyLimit(market environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.bittrexAPI.BuyLimit(market.Name, amount, limit)
	return orderNumber, err
}

//performs a market buy action.
func (wrapper BittrexWrapper) BuyMarket(market environment.Market, amount float64) (string, error) {
	orderNumber, err := wrapper.bittrexAPI.BuyMarket(market.Name, amount)
	return orderNumber, err
}

//performs a limit sell action.
func (wrapper BittrexWrapper) SellLimit(market environment.Market, amount float64, limit float64) (string, error) {
	orderNumber, err := wrapper.bittrexAPI.SellLimit(market.Name, amount, limit)
	return environment.Order{
		OrderNumber: orderNumber,
		Quantity:    amount,
		Value:       limit,
	}, err
}

//performs a market sell action.
func (wrapper BittrexWrapper) SellMarket(market environment.Market, amount float64) (string, error) {
	orderNumber, err := wrapper.bittrexAPI.SellMarket(market.Name, amount)
	return orderNumber, err
}
