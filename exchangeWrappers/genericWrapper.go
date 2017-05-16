package exchangeWrappers

import "github.com/AlessandroSanino1994/gobot/environment"

//ExchangeWrapper provides a generic wrapper for exchange services.
type ExchangeWrapper interface {
	GetCandles(market environment.Market, interval string) error //Gets the candles of a market.
	GetMarkets() ([]environment.Market, error)                   //Gets all the markets info.
	//GetTicker(market environment.Market)                                                         //Gets a ticker for a market.
	//GetSellBook(market environment.Market) ([]environment.Order, error)                          //Gets the sell(ASK) book of a market.
	//GetBuyBook(market environment.Market) ([]environment.Order, error)                           //Gets the buy(BID) book of a market.
	GetOrderBook(market environment.Market) error                                                  //Gets the order(ASK + BID) book of a market.
	BuyLimit(market environment.Market, amount float64, limit float64) (environment.Order, error)  //performs a limit buy action.
	BuyMarket(market environment.Market, amount float64) (environment.Order, error)                //performs a market buy action.
	SellLimit(market environment.Market, amount float64, limit float64) (environment.Order, error) //performs a limit sell action.
	SellMarket(market environment.Market, amount float64) (environment.Order, error)               //performs a market sell action.
}
