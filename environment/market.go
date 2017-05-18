package environment

//Ticker provides data incoming from API Tickers, which have little amount of information
//regarding very last updates from a market.
type Ticker struct {
	Ask  float64 `json:"Ask"`  //Represents ASK value from the ticker.
	Bid  float64 `json:"Bid"`  //Represents BID value from the ticker.
	Last float64 `json:"Last"` //Represents LAST trade value from the ticker.
}

//Market represents the environment the bot is trading in.
type Market struct {
	Name           string            //Represents the name of the market (e.g. ETH-BTC).
	BaseCurrency   string            //Represents the base currency of the market.
	MarketCurrency string            //Represents the currency to exchange by using base currency.
	WatchedChart   *CandleStickChart //Represents a map which contains all the charts watched currently by the bot.
	Summary        MarketSummary     //Represents the summary of the market.
}

//MarketSummary represents the summary data of a market.
type MarketSummary struct {
	High   float64 //Represents the 24 hours maximum peak of this market.
	Low    float64 //Represents the 24 hours minimum peak of this market.
	Volume float64 //Represents the 24 volume peak of this market.
	Ask    float64 //Represents the current ASK value.
	Bid    float64 //Represents the current BID value.
	Last   float64 //Represents the value of the last trade.
}

//UpdateFromTicker updates the values of the market summary from a Ticker Data.
func (summary MarketSummary) UpdateFromTicker(ticker Ticker) {
	summary.Ask = ticker.Ask
	summary.Bid = ticker.Bid
	summary.Last = ticker.Last
}

//To add a strategy implement strategy interface functions attached to a custom market.
//
//Like this:
/*
func (market *Market) SetUpStrategy() {

}

func (market *Market) TearDownStrategy() {

}

func (market *Market) OnCandleUpdate() {

}
*/
