package exchangeWrappers

import (
	"github.com/AlessandroSanino1994/gobot/environment"
	bittrexAPI "github.com/toorop/go-bittrex"
)

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
