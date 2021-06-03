package exchanges

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/juju/errors"
	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/shopspring/decimal"
)

// ExchangeWrapperSimulator wraps another wrapper and returns simulated balances and orders.
type ExchangeWrapperSimulator struct {
	innerWrapper ExchangeWrapper
	balances     map[string]decimal.Decimal
}

// NewExchangeWrapperSimulator creates a new simulated wrapper from another wrapper and an initial balance.
func NewExchangeWrapperSimulator(mockedWrapper ExchangeWrapper, initialBalances map[string]decimal.Decimal) *ExchangeWrapperSimulator {
	return &ExchangeWrapperSimulator{
		innerWrapper: mockedWrapper,
		balances:     initialBalances,
	}
}

// String returns a string representation of the exchange simulator.
func (wrapper *ExchangeWrapperSimulator) String() string {
	return wrapper.Name()
}

// Name gets the name of the exchange.
func (wrapper *ExchangeWrapperSimulator) Name() string {
	return fmt.Sprint(wrapper.innerWrapper.Name(), "mock")
}

// GetCandles gets the candle data from the exchange.
func (wrapper *ExchangeWrapperSimulator) GetCandles(market *environment.Market) ([]environment.CandleStick, error) {
	return wrapper.innerWrapper.GetCandles(market)
}

// GetMarketSummary gets the current market summary.
func (wrapper *ExchangeWrapperSimulator) GetMarketSummary(market *environment.Market) (*environment.MarketSummary, error) {
	return wrapper.innerWrapper.GetMarketSummary(market)
}

// GetOrderBook gets the order(ASK + BID) book of a market.
func (wrapper *ExchangeWrapperSimulator) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	return wrapper.innerWrapper.GetOrderBook(market)
}

// BuyLimit here is just to implement the ExchangeWrapper Interface, do not use, use BuyMarket instead.
func (wrapper *ExchangeWrapperSimulator) BuyLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	return "", errors.New("BuyLimit operation is not mockable")
}

// SellLimit here is just to implement the ExchangeWrapper Interface, do not use, use SellMarket instead.
func (wrapper *ExchangeWrapperSimulator) SellLimit(market *environment.Market, amount float64, limit float64) (string, error) {
	return "", errors.New("SellLimit operation is not mockable")
}

// BuyMarket performs a FAKE market buy action.
func (wrapper *ExchangeWrapperSimulator) BuyMarket(market *environment.Market, amount float64) (string, error) {
	baseBalance, _ := wrapper.GetBalance(market.BaseCurrency)
	quoteBalance, _ := wrapper.GetBalance(market.MarketCurrency)

	orderbook, err := wrapper.GetOrderBook(market)
	if err != nil {
		return "", errors.Annotate(err, "Cannot market buy without orderbook knowledge")
	}

	totalQuote := decimal.Zero
	remainingAmount := decimal.NewFromFloat(amount)
	expense := decimal.Zero

	for _, ask := range orderbook.Asks {
		if remainingAmount.LessThanOrEqual(ask.Quantity) {
			totalQuote = totalQuote.Add(remainingAmount)
			expense = expense.Add(remainingAmount.Mul(ask.Value))
			if expense.GreaterThan(*baseBalance) {
				return "", fmt.Errorf("cannot Buy not enough %s balance", market.BaseCurrency)
			}
			break
		}
		totalQuote = totalQuote.Add(ask.Quantity)
		expense = expense.Add(ask.Quantity.Mul(ask.Value))
		if expense.GreaterThan(*baseBalance) {
			return "", fmt.Errorf("cannot Buy not enough %s balance", market.BaseCurrency)
		}
	}

	wrapper.balances[market.BaseCurrency] = baseBalance.Sub(expense)
	wrapper.balances[market.MarketCurrency] = quoteBalance.Add(totalQuote)

	orderFakeID, err := uuid.NewV4()
	if err != nil {
		return "", errors.Annotate(err, "UUID Generation")
	}
	return fmt.Sprintf("FAKE_BUY-%s", orderFakeID), nil
}

// SellMarket performs a FAKE market buy action.
func (wrapper *ExchangeWrapperSimulator) SellMarket(market *environment.Market, amount float64) (string, error) {
	baseBalance, _ := wrapper.GetBalance(market.BaseCurrency)
	quoteBalance, _ := wrapper.GetBalance(market.MarketCurrency)

	orderbook, err := wrapper.GetOrderBook(market)
	if err != nil {
		return "", errors.Annotate(err, "Cannot market buy without orderbook knowledge")
	}

	totalQuote := decimal.Zero
	remainingAmount := decimal.NewFromFloat(amount)
	gain := decimal.Zero

	if quoteBalance.LessThan(remainingAmount) {
		return "", fmt.Errorf("Cannot Sell: not enough %s balance", market.MarketCurrency)
	}

	for _, bid := range orderbook.Bids {
		if remainingAmount.LessThanOrEqual(bid.Quantity) {
			totalQuote = totalQuote.Add(remainingAmount)
			gain = gain.Add(remainingAmount.Mul(bid.Value))
			break
		}
		totalQuote = totalQuote.Add(bid.Quantity)
		gain = gain.Add(bid.Quantity.Mul(bid.Value))
	}

	wrapper.balances[market.BaseCurrency] = baseBalance.Add(gain)
	wrapper.balances[market.MarketCurrency] = quoteBalance.Sub(totalQuote)

	orderFakeID, err := uuid.NewV4()
	if err != nil {
		return "", errors.Annotate(err, "UUID Generation")
	}
	return fmt.Sprintf("FAKE_SELL-%s", orderFakeID), nil
}

// CalculateTradingFees calculates the trading fees for an order on a specified market.
func (wrapper *ExchangeWrapperSimulator) CalculateTradingFees(market *environment.Market, amount float64, limit float64, orderType TradeType) float64 {
	return wrapper.innerWrapper.CalculateTradingFees(market, amount, limit, orderType)
}

// CalculateWithdrawFees calculates the withdrawal fees on a specified market.
func (wrapper *ExchangeWrapperSimulator) CalculateWithdrawFees(market *environment.Market, amount float64) float64 {
	return wrapper.innerWrapper.CalculateWithdrawFees(market, amount)
}

// GetBalance gets the balance of the user of the specified currency.
func (wrapper *ExchangeWrapperSimulator) GetBalance(symbol string) (*decimal.Decimal, error) {
	bal, exists := wrapper.balances[symbol]
	if !exists {
		wrapper.balances[symbol] = decimal.Zero
		var bal = decimal.Zero
		return &bal, nil
	}
	return &bal, nil
}

// GetDepositAddress gets the deposit address for the specified coin on the exchange.
func (wrapper *ExchangeWrapperSimulator) GetDepositAddress(coinTicker string) (string, bool) {
	return "", false
}

// FeedConnect connects to the feed of the exchange.
func (wrapper *ExchangeWrapperSimulator) FeedConnect(markets []*environment.Market) error {
	return wrapper.innerWrapper.FeedConnect(markets)
}

// Withdraw performs a FAKE withdraw operation from the exchange to a destination address.
func (wrapper *ExchangeWrapperSimulator) Withdraw(destinationAddress string, coinTicker string, amount float64) error {
	if amount <= 0 {
		return errors.New("Withdraw amount must be > 0")
	}

	bal, exists := wrapper.balances[coinTicker]
	amt := decimal.NewFromFloat(amount)
	if !exists || amt.GreaterThan(bal) {
		return errors.New("Not enough balance")
	}

	wrapper.balances[coinTicker] = bal.Sub(amt)

	return nil
}
