package helpers

import (
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/saniales/golang-crypto-trading-bot/exchanges"
	"github.com/shopspring/decimal"
)

//InitExchange initialize a new ExchangeWrapper binded to the specified exchange provided.
func InitExchange(exchangeConfig environment.ExchangeConfig, simulatedMode bool, fakeBalances map[string]decimal.Decimal) exchanges.ExchangeWrapper {
	var exch exchanges.ExchangeWrapper
	switch exchangeConfig.ExchangeName {
	case "bittrex":
		exch = exchanges.NewBittrexWrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "bittrexV2":
		exch = exchanges.NewBittrexV2Wrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "poloniex":
		exch = exchanges.NewPoloniexWrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "binance":
		exch = exchanges.NewBinanceWrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "bitfinex":
		exch = exchanges.NewBitfinexWrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	case "hitbtc":
		exch = exchanges.NewHitBtcV2Wrapper(exchangeConfig.PublicKey, exchangeConfig.SecretKey)
	default:
		return nil
	}

	if simulatedMode {
		if fakeBalances == nil {
			return nil
		}
		exch = exchanges.NewExchangeWrapperSimulator(exch, fakeBalances)
	}

	return exch
}
