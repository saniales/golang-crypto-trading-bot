package exchangeWrappers

import (
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/thebotguys/golang-bittrex-api/bittrex"
)

type BittrexWrapperV2 struct {
}

// GetMarkets gets all the markets info.
func (wrapper BittrexWrapperV2) GetMarkets() ([]*environment.Market, error) {
	bittrexMarkets, err := bittrex.GetMarkets()
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

/* GetOrderBook gets the order(ASK + BID) book of a market.
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
func (wrapper BittrexWrapperV2) GetTicker(market *environment.Market) error {
	bittrexTicker, err := wrapper.bittrexAPI.GetTicker(market.Name)
	if err != nil {
		return err
	}
	market.Summary.UpdateFromTicker(convertFromBittrexTicker(bittrexTicker))
	return nil
}*/

// GetMarketSummaries get the markets summary of all markets
func (wrapper BittrexWrapperV2) GetMarketSummaries(markets map[string]*environment.Market) error {
	bittrexSummaries, err := bittrex.GetMarketSummaries()
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
func (wrapper BittrexWrapperV2) GetMarketSummary(market *environment.Market) error {
	summary, err := bittrex.GetMarketSummary(market.Name)
	if err != nil {
		return err
	}

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
