// Copyright © 2017 Alessandro Sanino <saninoale@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package exchanges

import (
	"errors"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/shopspring/decimal"
	"github.com/thebotguys/golang-bittrex-api/bittrex"
)

// BittrexWrapperV2 wraps Bittrex API v2.0
type BittrexWrapperV2 struct {
	PublicKey        string
	SecretKey        string
	summaries        *SummaryCache
	depositAddresses map[string]string
}

// NewBittrexV2Wrapper creates a generic wrapper of the bittrex API v2.0.
func NewBittrexV2Wrapper(publicKey string, secretKey string, depositAddresses map[string]string) ExchangeWrapper {
	return &BittrexWrapperV2{
		PublicKey:        publicKey,
		SecretKey:        secretKey,
		summaries:        NewSummaryCache(),
		depositAddresses: depositAddresses,
	}
}

// Name returns the name of the wrapped exchange.
func (wrapper *BittrexWrapperV2) Name() string {
	return "bittrex"
}

func (wrapper *BittrexWrapperV2) String() string {
	return wrapper.Name()
}

// GetMarkets gets all the markets info.
func (wrapper *BittrexWrapperV2) GetMarkets() ([]*environment.Market, error) {
	bittrexMarkets, err := bittrex.GetMarkets()
	if err != nil {
		return nil, err
	}
	wrappedMarkets := make([]*environment.Market, 0, len(bittrexMarkets))
	for _, market := range bittrexMarkets {
		if market.IsActive {
			wrappedMarkets = append(wrappedMarkets, &environment.Market{
				Name:           market.MarketName,
				BaseCurrency:   market.BaseCurrency,
				MarketCurrency: market.MarketCurrency,
			})
		}
	}
	return wrappedMarkets, nil
}

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper *BittrexWrapperV2) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	panic("GetOrderBook not implemented")
}

// BuyLimit performs a limit buy action.
func (wrapper *BittrexWrapperV2) BuyLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	return "", errors.New("BuyLimit not implemented")
}

// BuyMarket performs a market buy action.
func (wrapper *BittrexWrapperV2) BuyMarket(market *environment.Market, amount float64) (string, error) {
	return "", errors.New("BuyMarket not implemented")
}

// SellLimit performs a limit sell action.
func (wrapper *BittrexWrapperV2) SellLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	return "", errors.New("SellLimit not implemented")
}

// SellMarket performs a market sell action.
func (wrapper *BittrexWrapperV2) SellMarket(market *environment.Market, amount float64) (string, error) {
	return "", errors.New("SellMarket not implemented")
}

// GetMarketSummary gets the current market summary.
func (wrapper *BittrexWrapperV2) GetMarketSummary(market *environment.Market) (*environment.MarketSummary, error) {
	summary, err := bittrex.GetMarketSummary(market.Name)
	if err != nil {
		return nil, err
	}

	ret := &environment.MarketSummary{
		High:   summary.High,
		Low:    summary.Low,
		Volume: summary.Volume,
		Bid:    summary.Bid,
		Ask:    summary.Ask,
		Last:   summary.Last,
	}

	wrapper.summaries.Set(market, ret)

	return ret, nil
}

// GetCandles gets the candle data from the exchange.
func (wrapper *BittrexWrapperV2) GetCandles(market *environment.Market) ([]environment.CandleStick, error) {
	bittrexCandles, err := bittrex.GetTicks(MarketNameFor(market, wrapper), "30m")
	if err != nil {
		return nil, err
	}

	ret := make([]environment.CandleStick, len(bittrexCandles))

	for i, bittrexCandle := range bittrexCandles {
		ret[i] = environment.CandleStick{
			High:   bittrexCandle.High,
			Open:   bittrexCandle.Open,
			Close:  bittrexCandle.Close,
			Low:    bittrexCandle.Low,
			Volume: bittrexCandle.BaseVolume,
		}
	}

	return ret, nil
}

// GetBalance gets the balance of the user of the specified currency.
func (wrapper *BittrexWrapperV2) GetBalance(symbol string) (*decimal.Decimal, error) {
	panic("Not Implemented")
}

// GetDepositAddress gets the deposit address for the specified coin on the exchange.
func (wrapper *BittrexWrapperV2) GetDepositAddress(coinTicker string) (string, bool) {
	addr, exists := wrapper.depositAddresses[coinTicker]
	return addr, exists
}

// CalculateTradingFees calculates the trading fees for an order on a specified market.
//
//     NOTE: In Bittrex fees are hardcoded due to the inability to obtain them via API before placing an order.
func (wrapper *BittrexWrapperV2) CalculateTradingFees(market *environment.Market, amount float64, limit float64, orderType TradeType) float64 {
	var feePercentage float64
	if orderType == MakerTrade {
		feePercentage = 0.0025
	} else if orderType == TakerTrade {
		feePercentage = 0.0025
	} else {
		panic("Unknown trade type")
	}

	return amount * limit * feePercentage
}

// CalculateWithdrawFees calculates the withdrawal fees on a specified market.
func (wrapper *BittrexWrapperV2) CalculateWithdrawFees(market *environment.Market, amount float64) float64 {
	panic("Not Implemented")
}

// FeedConnect connects to the feed of the exchange.
//
//     NOTE: Not supported on Bittrex v1 API, use *BittrexWrapperV2.
func (wrapper *BittrexWrapperV2) FeedConnect(markets []*environment.Market) error {
	return ErrWebsocketNotSupported
}

// Withdraw performs a withdraw operation from the exchange to a destination address.
func (wrapper *BittrexWrapperV2) Withdraw(destinationAddress string, coinTicker string, amount float64) error {
	panic("Not Implemented")
}
