package helpers

import (
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
)

//InitExchange initialize a new ExchangeWrapper binded to the specified exchange provided.
func InitExchange(exchangeConfig environment.ExchangeConfig) exchanges.ExchangeWrapper {
	switch exchangeConfig.ExchangeName {
	case "bittrex":
		return exchanges.NewBittrexWrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "bittrexV2":
		return exchanges.NewBittrexV2Wrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "poloniex":
		return exchanges.NewPoloniexWrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "binance":
		return exchanges.NewBinanceWrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "bitfinex":
		return exchanges.NewBitfinexWrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "hitbtc":
		return exchanges.NewHitBtcV2Wrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	default:
		return nil
	}
}
