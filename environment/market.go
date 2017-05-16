package environment

//Market represents the environment the bot is trading in.
type Market struct {
	Name           string            //Represents the name of the market (e.g. ETH-BTC).
	BaseCurrency   string            //Represents the base currency of the market.
	MarketCurrency string            //Represents the currency to exchange by using base currency.
	WatchedChart   *CandleStickChart //Represents a map which contains all the charts watched currently by the bot.
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
