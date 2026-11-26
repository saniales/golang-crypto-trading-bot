package exchanges

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/saniales/golang-crypto-trading-bot/environment"
	"github.com/shopspring/decimal"
)

const coinbaseAPIURL = "https://api.exchange.coinbase.com"

// CoinbaseWrapper provides a generic wrapper for the Coinbase Pro API.
type CoinbaseWrapper struct {
	apiURL           string
	client           *http.Client
	publicKey        string
	secretKey        string
	passphrase       string
	depositAddresses map[string]string
}

// NewCoinbaseWrapper creates a generic wrapper of the Coinbase Pro API.
func NewCoinbaseWrapper(publicKey, secretKey string, depositAddresses map[string]string, passphrases ...string) ExchangeWrapper {
	var passphrase string
	if len(passphrases) > 0 {
		passphrase = passphrases[0]
	}
	return &CoinbaseWrapper{
		apiURL:           coinbaseAPIURL,
		client:           &http.Client{Timeout: 15 * time.Second},
		publicKey:        publicKey,
		secretKey:        secretKey,
		passphrase:       passphrase,
		depositAddresses: depositAddresses,
	}
}

func (wrapper *CoinbaseWrapper) Name() string {
	return "coinbase"
}

func (wrapper *CoinbaseWrapper) String() string {
	return wrapper.Name()
}

type coinbaseProduct struct {
	ID            string `json:"id"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
}

// GetMarkets gets all the markets info.
func (wrapper *CoinbaseWrapper) GetMarkets() ([]*environment.Market, error) {
	var products []coinbaseProduct
	if err := wrapper.request(http.MethodGet, "/products", nil, false, &products); err != nil {
		return nil, err
	}
	markets := make([]*environment.Market, len(products))
	for i, product := range products {
		markets[i] = &environment.Market{
			Name:           product.ID,
			BaseCurrency:   product.BaseCurrency,
			MarketCurrency: product.QuoteCurrency,
		}
	}
	return markets, nil
}

type coinbaseBook struct {
	Bids [][]json.RawMessage `json:"bids"`
	Asks [][]json.RawMessage `json:"asks"`
}

// GetOrderBook gets the order book of a market.
func (wrapper *CoinbaseWrapper) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	var book coinbaseBook
	if err := wrapper.request(http.MethodGet, "/products/"+url.PathEscape(MarketNameFor(market, wrapper))+"/book?level=2", nil, false, &book); err != nil {
		return nil, err
	}
	result := &environment.OrderBook{}
	var err error
	if result.Bids, err = parseCoinbaseLevels(book.Bids); err != nil {
		return nil, err
	}
	if result.Asks, err = parseCoinbaseLevels(book.Asks); err != nil {
		return nil, err
	}
	return result, nil
}

func parseCoinbaseLevels(levels [][]json.RawMessage) ([]environment.Order, error) {
	orders := make([]environment.Order, 0, len(levels))
	for _, level := range levels {
		if len(level) < 2 {
			return nil, errors.New("invalid Coinbase order book level")
		}
		var priceText, sizeText string
		if err := json.Unmarshal(level[0], &priceText); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(level[1], &sizeText); err != nil {
			return nil, err
		}
		price, err := decimal.NewFromString(priceText)
		if err != nil {
			return nil, err
		}
		size, err := decimal.NewFromString(sizeText)
		if err != nil {
			return nil, err
		}
		orders = append(orders, environment.Order{Value: price, Quantity: size})
	}
	return orders, nil
}

type coinbaseTicker struct {
	Price string `json:"price"`
	Bid   string `json:"bid"`
	Ask   string `json:"ask"`
}

type coinbaseStats struct {
	High   string `json:"high"`
	Low    string `json:"low"`
	Volume string `json:"volume"`
}

func (wrapper *CoinbaseWrapper) getTicker(market *environment.Market) (*coinbaseTicker, error) {
	var ticker coinbaseTicker
	err := wrapper.request(http.MethodGet, "/products/"+url.PathEscape(MarketNameFor(market, wrapper))+"/ticker", nil, false, &ticker)
	return &ticker, err
}

// GetMarketSummary gets the current market summary.
func (wrapper *CoinbaseWrapper) GetMarketSummary(market *environment.Market) (*environment.MarketSummary, error) {
	ticker, err := wrapper.getTicker(market)
	if err != nil {
		return nil, err
	}
	var stats coinbaseStats
	if err := wrapper.request(http.MethodGet, "/products/"+url.PathEscape(MarketNameFor(market, wrapper))+"/stats", nil, false, &stats); err != nil {
		return nil, err
	}
	high, err := decimal.NewFromString(stats.High)
	if err != nil {
		return nil, err
	}
	low, err := decimal.NewFromString(stats.Low)
	if err != nil {
		return nil, err
	}
	volume, err := decimal.NewFromString(stats.Volume)
	if err != nil {
		return nil, err
	}
	last, err := decimal.NewFromString(ticker.Price)
	if err != nil {
		return nil, err
	}
	bid, err := decimal.NewFromString(ticker.Bid)
	if err != nil {
		return nil, err
	}
	ask, err := decimal.NewFromString(ticker.Ask)
	if err != nil {
		return nil, err
	}
	return &environment.MarketSummary{High: high, Low: low, Volume: volume, Last: last, Bid: bid, Ask: ask}, nil
}

// GetCandles gets the candle data from the exchange.
func (wrapper *CoinbaseWrapper) GetCandles(market *environment.Market) ([]environment.CandleStick, error) {
	granularity := coinbaseGranularity(MarketTimeFrameFor(market, wrapper))
	path := fmt.Sprintf("/products/%s/candles?granularity=%d", url.PathEscape(MarketNameFor(market, wrapper)), granularity)
	var rows [][]json.Number
	if err := wrapper.request(http.MethodGet, path, nil, false, &rows); err != nil {
		return nil, err
	}
	candles := make([]environment.CandleStick, 0, len(rows))
	for _, row := range rows {
		if len(row) < 6 {
			return nil, errors.New("invalid Coinbase candle")
		}
		values := make([]decimal.Decimal, 5)
		for i, index := range []int{1, 2, 3, 4, 5} {
			value, err := decimal.NewFromString(row[index].String())
			if err != nil {
				return nil, err
			}
			values[i] = value
		}
		candles = append(candles, environment.CandleStick{
			Low: values[0], High: values[1], Open: values[2], Close: values[3], Volume: values[4],
		})
	}
	return candles, nil
}

func coinbaseGranularity(timeframe string) int {
	switch strings.ToLower(timeframe) {
	case "1m", "1min", "60":
		return 60
	case "5m", "5min", "300":
		return 300
	case "15m", "15min", "900":
		return 900
	case "6h", "6hr", "21600":
		return 21600
	case "1d", "24h", "86400":
		return 86400
	default:
		return 3600
	}
}

// BuyLimit performs a limit buy action.
func (wrapper *CoinbaseWrapper) BuyLimit(market *environment.Market, amount, limit float64) (string, error) {
	return wrapper.createOrder(market, "buy", "limit", amount, limit)
}

// SellLimit performs a limit sell action.
func (wrapper *CoinbaseWrapper) SellLimit(market *environment.Market, amount, limit float64) (string, error) {
	return wrapper.createOrder(market, "sell", "limit", amount, limit)
}

// BuyMarket performs a market buy action.
func (wrapper *CoinbaseWrapper) BuyMarket(market *environment.Market, amount float64) (string, error) {
	return wrapper.createOrder(market, "buy", "market", amount, 0)
}

// SellMarket performs a market sell action.
func (wrapper *CoinbaseWrapper) SellMarket(market *environment.Market, amount float64) (string, error) {
	return wrapper.createOrder(market, "sell", "market", amount, 0)
}

func (wrapper *CoinbaseWrapper) createOrder(market *environment.Market, side, orderType string, amount, limit float64) (string, error) {
	order := map[string]string{
		"product_id": MarketNameFor(market, wrapper),
		"side":       side,
		"type":       orderType,
		"size":       strconv.FormatFloat(amount, 'f', -1, 64),
	}
	if orderType == "limit" {
		order["price"] = strconv.FormatFloat(limit, 'f', -1, 64)
	}
	var response struct {
		ID string `json:"id"`
	}
	if err := wrapper.request(http.MethodPost, "/orders", order, true, &response); err != nil {
		return "", err
	}
	return response.ID, nil
}

// CalculateTradingFees calculates the trading fees for an order on a market.
func (wrapper *CoinbaseWrapper) CalculateTradingFees(market *environment.Market, amount, limit float64, orderType TradeType) float64 {
	if orderType != MakerTrade && orderType != TakerTrade {
		panic("Unknown trade type")
	}
	return amount * limit * 0.005
}

// CalculateWithdrawFees calculates the withdrawal fees on a market.
func (wrapper *CoinbaseWrapper) CalculateWithdrawFees(market *environment.Market, amount float64) float64 {
	return 0
}

// GetBalance gets the available balance of the specified currency.
func (wrapper *CoinbaseWrapper) GetBalance(symbol string) (*decimal.Decimal, error) {
	var accounts []struct {
		Currency  string `json:"currency"`
		Available string `json:"available"`
	}
	if err := wrapper.request(http.MethodGet, "/accounts", nil, true, &accounts); err != nil {
		return nil, err
	}
	for _, account := range accounts {
		if account.Currency == symbol {
			balance, err := decimal.NewFromString(account.Available)
			if err != nil {
				return nil, err
			}
			return &balance, nil
		}
	}
	return nil, errors.New("Symbol not found")
}

// GetDepositAddress gets the configured deposit address for a coin.
func (wrapper *CoinbaseWrapper) GetDepositAddress(coinTicker string) (string, bool) {
	address, exists := wrapper.depositAddresses[coinTicker]
	return address, exists
}

// FeedConnect connects to the exchange feed.
func (wrapper *CoinbaseWrapper) FeedConnect(markets []*environment.Market) error {
	return ErrWebsocketNotSupported
}

// Withdraw transfers funds to a destination address.
func (wrapper *CoinbaseWrapper) Withdraw(destinationAddress, coinTicker string, amount float64) error {
	body := map[string]string{
		"amount":         strconv.FormatFloat(amount, 'f', -1, 64),
		"currency":       coinTicker,
		"crypto_address": destinationAddress,
	}
	return wrapper.request(http.MethodPost, "/withdrawals/crypto", body, true, nil)
}

func (wrapper *CoinbaseWrapper) request(method, path string, body interface{}, authenticated bool, result interface{}) error {
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	request, err := http.NewRequest(method, wrapper.apiURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "golang-crypto-trading-bot")
	if authenticated {
		if wrapper.publicKey == "" || wrapper.secretKey == "" || wrapper.passphrase == "" {
			return errors.New("Coinbase API key, secret, and passphrase are required")
		}
		secret, err := base64.StdEncoding.DecodeString(wrapper.secretKey)
		if err != nil {
			return fmt.Errorf("decode Coinbase API secret: %w", err)
		}
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		message := timestamp + method + request.URL.RequestURI() + string(payload)
		mac := hmac.New(sha256.New, secret)
		_, _ = mac.Write([]byte(message))
		request.Header.Set("CB-ACCESS-KEY", wrapper.publicKey)
		request.Header.Set("CB-ACCESS-SIGN", base64.StdEncoding.EncodeToString(mac.Sum(nil)))
		request.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
		request.Header.Set("CB-ACCESS-PASSPHRASE", wrapper.passphrase)
	}
	response, err := wrapper.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 4<<20))
	if err != nil {
		return err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Coinbase API returned %s: %s", response.Status, strings.TrimSpace(string(responseBody)))
	}
	if result == nil || len(responseBody) == 0 {
		return nil
	}
	return json.Unmarshal(responseBody, result)
}

var _ ExchangeWrapper = (*CoinbaseWrapper)(nil)
