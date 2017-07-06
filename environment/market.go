package environment

import (
	"fmt"
	"strings"
)

//Ticker provides data incoming from API Tickers, which have little amount of information
//regarding very last updates from a market.
type Ticker struct {
	Ask  float64 `json:"Ask"`  //Represents ASK value from the ticker.
	Bid  float64 `json:"Bid"`  //Represents BID value from the ticker.
	Last float64 `json:"Last"` //Represents LAST trade value from the ticker.
}

//Market represents the environment the bot is trading in.
type Market struct {
	Name           string            `json:"name,required"`            //Represents the name of the market (e.g. ETH-BTC).
	BaseCurrency   string            `json:"baseCurrency,omitempty"`   //Represents the base currency of the market.
	MarketCurrency string            `json:"marketCurrency,omitempty"` //Represents the currency to exchange by using base currency.
	WatchedChart   *CandleStickChart `json:"chart,omitempty"`          //Represents a map which contains all the charts watched currently by the bot.
	Summary        MarketSummary     `json:"summary,required"`         //Represents the summary of the market.
}

func (m Market) String() string {
	ret := fmt.Sprintln("Market", m.Name)
	ret += fmt.Sprintln("Summary :")
	ret += fmt.Sprintln(m.Summary)
	return strings.TrimSpace(ret)
}

//MarketSummary represents the summary data of a market.
type MarketSummary struct {
	High   float64 `json:"high,required"`   //Represents the 24 hours maximum peak of this market.
	Low    float64 `json:"low,required"`    //Represents the 24 hours minimum peak of this market.
	Volume float64 `json:"volume,required"` //Represents the 24 volume peak of this market.
	Ask    float64 `json:"ask,required"`    //Represents the current ASK value.
	Bid    float64 `json:"bid,required"`    //Represents the current BID value.
	Last   float64 `json:"last,required"`   //Represents the value of the last trade.
}

func (ms MarketSummary) String() string {
	ret := fmt.Sprintf("  Last: %.8f\n", ms.Last)
	ret += fmt.Sprintf("  ASK: %.8f\n", ms.Ask)
	ret += fmt.Sprintf("  BID: %.8f\n", ms.Bid)
	ret += fmt.Sprintf("  Volume: %.2f\n", ms.Volume)
	ret += fmt.Sprintf("  High: %.8f\n", ms.High)
	ret += fmt.Sprintf("  Low: %.8f\n", ms.Low)
	return ret
}

//UpdateFromTicker updates the values of the market summary from a Ticker Data.
func (ms *MarketSummary) UpdateFromTicker(ticker Ticker) {
	ms.Ask = ticker.Ask
	ms.Bid = ticker.Bid
	ms.Last = ticker.Last
}
