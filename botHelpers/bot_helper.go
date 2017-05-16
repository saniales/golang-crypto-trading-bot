package botHelpers

import "github.com/AlessandroSanino1994/gobot/exchangeWrappers"

//ParseArgs parses command line flags and returns initialized variables, or error.
func ParseArgs() {

}

//InitExchange initialize a new ExchangeWrapper binded to the specified exchange provided.
func InitExchange(exchangeName string, publicKey string, secretKey string) {
	if exchangeName == "bittrex" {
		return exchangeWrappers.NewBittrexWrapper(publicKey, secretKey)
	}
}

//InitMarkets uses ExchangeWrapper to find info about markets and initialize them.
func InitMarkets(exchange exchangeWrappers.ExchangeWrapper) {
	exchange.GetMarkets()
}
