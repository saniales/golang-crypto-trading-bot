package exchanges

import (
	"net/http"
	"time"

	"github.com/preichenberger/go-gdax"
	"github.com/saniales/golang-crypto-trading-bot/environment"
)

type booksLevel int

func (bl booksLevel) toInt() int {
	return int(bl)
}

const (
	//Only the best bid and ask
	bookLevelBestBidsAsks booksLevel = iota + 1
	//Top 50 bids and asks (aggregated)
	bookLevelTop50BidsAsks
	//Full order book (non aggregated)
	bookLevelFullOrderBook
)

type CoinBaseConfig struct {
	Secret, Key, Passphrase string
	BooksLavel              booksLevel
}

type CoinBaseWrapper struct {
	api        *gdax.Client
	summaries  *SummaryCache
	candles    *CandlesCache
	booksLevel booksLevel
}

func NewCoinBaseWrapper(config *CoinBaseConfig) *CoinBaseWrapper {
	client := gdax.NewClient(config.Key, config.Secret, config.Passphrase)
	client.HttpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	return &CoinBaseWrapper{
		api:        client,
		booksLevel: config.BooksLavel,
	}
}

func (CoinBaseWrapper) Name() string {
	return "coin_base"
}

func (cb CoinBaseWrapper) String() string {
	return cb.Name()
}

func (cb *CoinBaseWrapper) GetMarkets() ([]*environment.Market, error) {
	products, err := cb.api.GetProducts()
	if err != nil {
		return nil, err
	}
	markets := make([]*environment.Market, 0, len(products))
	for _, v := range products {
		markets = append(markets, &environment.Market{
			Name:           v.Id,
			BaseCurrency:   v.BaseCurrency,
			MarketCurrency: v.QuoteCurrency,
		})
	}
	return markets, nil
}

func (cb *CoinBaseWrapper) GetOrderBook(market *environment.Market) (*environment.OrderBook, error) {
	book, err := cb.api.GetBook(market.Name, cb.booksLevel.toInt()) //todo: add "MarketNameFor(market, cb)" after fill interface
	if err != nil {
		return nil, err
	}

	var orderBook environment.OrderBook
	for _, v := range book.Asks {
		amount := decimal.NewFromFloat(v.Size)
		rate := decimal.NewFromFloat(v.Price)
		orderBook.Asks = append(orderBook.Asks, environment.Order{
			Quantity: amount,
			Value:    rate,
		})
	}

	for _, v := range book.Bids {
		amount := decimal.NewFromFloat(v.Size)
		rate := decimal.NewFromFloat(v.Price)
		orderBook.Bids = append(orderBook.Bids, environment.Order{
			Quantity: amount,
			Value:    rate,
		})
	}
	return &orderBook, nil

}
