package exchangeWrappers

//ExchangeWrapper provides a generic wrapper for exchange services.
type ExchangeWrapper interface {
	GetCandles(marketName string, time string)                  //Gets the candles of a market.
	GetMarkets()                                                //Gets all the markets info.
	GetTicker(marketName string)                                //Gets a ticker for a market.
	GetSellBook(marketName string)                              //Gets the sell(ASK) book of a market.
	GetBuyBook(marketName string)                               //Gets the buy(BID) book of a market.
	GetOrderBook(marketName string)                             //Gets the order(ASK + BID) book of a market.
	BuyLimit(marketName string, amount float64, limit float64)  //performs a limit buy action.
	BuyMarket(marketName string, amount float64)                //performs a market buy action.
	SellLimit(marketName string, amount float64, limit float64) //performs a limit sell action.
	SellMarket(marketName string, amount float64)               //performs a market sell action.
}
