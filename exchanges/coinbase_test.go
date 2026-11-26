package exchanges

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saniales/golang-crypto-trading-bot/environment"
)

func TestCoinbaseWrapperImplementsExchangeWrapper(t *testing.T) {
	var wrapper ExchangeWrapper = NewCoinbaseWrapper("", "", nil)
	if wrapper.Name() != "coinbase" {
		t.Fatalf("unexpected exchange name %q", wrapper.Name())
	}
	if err := wrapper.FeedConnect(nil); err != ErrWebsocketNotSupported {
		t.Fatalf("FeedConnect() error = %v, want %v", err, ErrWebsocketNotSupported)
	}
}

func TestCoinbaseMarketSummaryAndCandles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/products/BTC-USD/ticker":
			_, _ = w.Write([]byte(`{"price":"101","bid":"100","ask":"102"}`))
		case "/products/BTC-USD/stats":
			_, _ = w.Write([]byte(`{"high":"110","low":"90","volume":"12"}`))
		case "/products/BTC-USD/candles":
			if r.URL.Query().Get("granularity") != "300" {
				t.Errorf("granularity = %q, want 300", r.URL.Query().Get("granularity"))
			}
			_, _ = w.Write([]byte(`[[123,1,4,2,3,5]]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	wrapper := NewCoinbaseWrapper("", "", nil).(*CoinbaseWrapper)
	wrapper.apiURL = server.URL
	wrapper.client = server.Client()
	market := &environment.Market{
		ExchangeNames:      map[string]string{"coinbase": "BTC-USD"},
		ExchangeTimeFrames: map[string]string{"coinbase": "5m"},
	}

	summary, err := wrapper.GetMarketSummary(market)
	if err != nil {
		t.Fatal(err)
	}
	if summary.High.String() != "110" || summary.Last.String() != "101" || summary.Bid.String() != "100" {
		t.Fatalf("unexpected market summary: %+v", summary)
	}

	candles, err := wrapper.GetCandles(market)
	if err != nil {
		t.Fatal(err)
	}
	if len(candles) != 1 || candles[0].Low.String() != "1" || candles[0].High.String() != "4" || candles[0].Volume.String() != "5" {
		t.Fatalf("unexpected candles: %+v", candles)
	}
}

func TestCoinbaseAuthenticatedRequests(t *testing.T) {
	secret := base64.StdEncoding.EncodeToString([]byte("secret"))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts" {
			t.Errorf("request path = %q, want /accounts", r.URL.Path)
		}
		for _, header := range []string{"CB-ACCESS-KEY", "CB-ACCESS-SIGN", "CB-ACCESS-TIMESTAMP", "CB-ACCESS-PASSPHRASE"} {
			if r.Header.Get(header) == "" {
				t.Errorf("missing authentication header %s", header)
			}
		}
		_, _ = w.Write([]byte(`[{"currency":"USD","available":"12.5"}]`))
	}))
	defer server.Close()

	wrapper := NewCoinbaseWrapper("key", secret, nil, "passphrase").(*CoinbaseWrapper)
	wrapper.apiURL = server.URL
	wrapper.client = server.Client()
	balance, err := wrapper.GetBalance("USD")
	if err != nil {
		t.Fatal(err)
	}
	if balance.String() != "12.5" {
		t.Fatalf("balance = %s, want 12.5", balance)
	}
}

func TestCoinbaseOrderBook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"bids": [][]string{{"100", "2", "3"}},
			"asks": [][]string{{"101", "4", "5"}},
		})
	}))
	defer server.Close()

	wrapper := NewCoinbaseWrapper("", "", nil).(*CoinbaseWrapper)
	wrapper.apiURL = server.URL
	wrapper.client = server.Client()
	book, err := wrapper.GetOrderBook(&environment.Market{ExchangeNames: map[string]string{"coinbase": "BTC-USD"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(book.Bids) != 1 || book.Bids[0].Value.String() != "100" || book.Bids[0].Quantity.String() != "2" {
		t.Fatalf("unexpected order book: %+v", book)
	}
}
