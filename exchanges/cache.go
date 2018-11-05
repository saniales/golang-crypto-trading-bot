package exchanges

import (
	"sync"

	"github.com/saniales/golang-crypto-trading-bot/environment"
)

// SummaryCache represents a local summary cache for every exchange. To allow dinamic polling from multiple sources (REST + Websocket)
type SummaryCache struct {
	mutex    *sync.RWMutex
	internal map[*environment.Market]*environment.MarketSummary
}

// NewSummaryCache creates a new SummaryCache Object
func NewSummaryCache() *SummaryCache {
	return &SummaryCache{
		mutex:    &sync.RWMutex{},
		internal: make(map[*environment.Market]*environment.MarketSummary),
	}
}

// Set sets a value for the specified key.
func (sc *SummaryCache) Set(market *environment.Market, summary *environment.MarketSummary) *environment.MarketSummary {
	sc.mutex.Lock()
	old := sc.internal[market]
	sc.internal[market] = summary
	sc.mutex.Unlock()
	return old
}

// Get gets the value for the specified key.
func (sc *SummaryCache) Get(market *environment.Market) (*environment.MarketSummary, bool) {
	sc.mutex.RLock()
	ret, isSet := sc.internal[market]
	sc.mutex.RUnlock()
	return ret, isSet
}

// CandlesCache represents a local candles cache for every exchange. To allow dinamic polling from multiple sources (REST + Websocket)
type CandlesCache struct {
	mutex    *sync.RWMutex
	internal map[*environment.Market][]environment.CandleStick
}

// NewCandlesCache creates a new CandlesCache Object
func NewCandlesCache() *CandlesCache {
	return &CandlesCache{
		mutex:    &sync.RWMutex{},
		internal: make(map[*environment.Market][]environment.CandleStick),
	}
}

// Set sets a value for the specified key.
func (cc *CandlesCache) Set(market *environment.Market, candles []environment.CandleStick) []environment.CandleStick {
	cc.mutex.Lock()
	old := cc.internal[market]
	cc.internal[market] = candles
	cc.mutex.Unlock()
	return old
}

// Get gets the value for the specified key.
func (cc *CandlesCache) Get(market *environment.Market) ([]environment.CandleStick, bool) {
	cc.mutex.RLock()
	ret, isSet := cc.internal[market]
	cc.mutex.RUnlock()
	return ret, isSet
}

// OrderbookCache represents a local orderbook cache for every exchange. To allow dinamic polling from multiple sources (REST + Websocket)
type OrderbookCache struct {
	mutex    *sync.RWMutex
	internal map[*environment.Market]*environment.OrderBook
}

// NewOrderbookCache creates a new OrderbookCache Object
func NewOrderbookCache() *OrderbookCache {
	return &OrderbookCache{
		mutex:    &sync.RWMutex{},
		internal: make(map[*environment.Market]*environment.OrderBook),
	}
}

// Set sets a value for the specified key.
func (cc *OrderbookCache) Set(market *environment.Market, book *environment.OrderBook) *environment.OrderBook {
	cc.mutex.Lock()
	old := cc.internal[market]
	cc.internal[market] = book
	cc.mutex.Unlock()
	return old
}

// Get gets the value for the specified key.
func (cc *OrderbookCache) Get(market *environment.Market) (*environment.OrderBook, bool) {
	cc.mutex.RLock()
	ret, isSet := cc.internal[market]
	cc.mutex.RUnlock()
	return ret, isSet
}
