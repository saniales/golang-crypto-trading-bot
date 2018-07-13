package botHelpers

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
		return nil
	case "yobit":
		return nil
	case "cryptopia":
		return nil
	default:
		return nil
	}
}
