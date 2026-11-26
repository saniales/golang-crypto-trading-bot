package exchanges

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/adshao/go-binance/v2"
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fiore/kucoin-go"
	"github.com/saniales/golang-crypto-trading-bot/environment"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(body string, status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestBinanceCancelOrder(t *testing.T) {
	client := binance.NewClient("key", "secret")
	requests := 0
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requests++
		if req.Method != http.MethodDelete || req.URL.Path != "/api/v3/order" {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatal(err)
		}
		if form.Get("symbol") != "ETHBTC" || form.Get("origClientOrderId") != "client-id" || form.Get("orderId") != "" {
			t.Errorf("unexpected cancellation parameters: %v", form)
		}
		if requests == 2 {
			return jsonResponse(`{"code":-2011,"msg":"Unknown order"}`, http.StatusBadRequest), nil
		}
		return jsonResponse(`{}`, http.StatusOK), nil
	})}
	wrapper := &BinanceWrapper{api: client}
	market := &environment.Market{ExchangeNames: map[string]string{"binance": "ETHBTC"}}
	if err := wrapper.CancelOrder(market, "client-id"); err != nil {
		t.Fatal(err)
	}
	if err := wrapper.CancelOrder(market, "client-id"); err == nil {
		t.Fatal("expected exchange error")
	}
	if err := wrapper.CancelOrder(market, ""); err == nil || requests != 2 {
		t.Fatal("empty ID should fail without a request")
	}
}

func TestKrakenCancelOrder(t *testing.T) {
	count := 1
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/0/private/CancelOrder" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatal(err)
		}
		if form.Get("txid") != "transaction-id" {
			t.Errorf("unexpected txid: %q", form.Get("txid"))
		}
		if count == 0 {
			return jsonResponse(`{"error":[],"result":{"count":0}}`, http.StatusOK), nil
		}
		return jsonResponse(`{"error":[],"result":{"count":1}}`, http.StatusOK), nil
	})}
	wrapper := &KrakenWrapper{api: krakenapi.NewWithClient("key", "secret", client)}
	if err := wrapper.CancelOrder(nil, "[transaction-id]"); err != nil {
		t.Fatal(err)
	}
	count = 0
	if err := wrapper.CancelOrder(nil, "transaction-id"); err == nil {
		t.Fatal("expected zero canceled orders to fail")
	}
	for _, id := range []string{"", "[]", "[one two]", "one,two"} {
		if err := wrapper.CancelOrder(nil, id); err == nil {
			t.Errorf("expected invalid ID %q to fail", id)
		}
	}
}

func TestKucoinCancelOrder(t *testing.T) {
	side := ""
	client := http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch {
		case strings.HasSuffix(req.URL.Path, "/market/open/coins-trending"):
			return jsonResponse(`{"success":true,"data":[{"coinPair":"ETH-BTC"}]}`, http.StatusOK), nil
		case strings.HasSuffix(req.URL.Path, "/order/active-map"):
			if req.URL.Query().Get("symbol") != "ETH-BTC" {
				t.Errorf("unexpected market: %s", req.URL.RawQuery)
			}
			return jsonResponse(`{"success":true,"data":{"BUY":[{"oid":"buy-id"}],"SELL":[{"oid":"sell-id"}]}}`, http.StatusOK), nil
		case strings.HasSuffix(req.URL.Path, "/cancel-order"):
			if err := req.ParseForm(); err != nil {
				t.Fatal(err)
			}
			side = req.Form.Get("type")
			if req.Form.Get("symbol") != "ETH-BTC" || req.Form.Get("orderOid") != strings.ToLower(side)+"-id" {
				t.Errorf("unexpected cancellation parameters: %v", req.Form)
			}
			return jsonResponse(`{"success":true}`, http.StatusOK), nil
		default:
			t.Fatalf("unexpected request: %s", req.URL)
			return nil, nil
		}
	})}
	wrapper := &KucoinWrapper{api: kucoin.NewCustomClient("key", "secret", client)}
	market := &environment.Market{ExchangeNames: map[string]string{"kucoin": "ETH-BTC"}}
	for _, want := range []string{"BUY", "SELL"} {
		if err := wrapper.CancelOrder(market, strings.ToLower(want)+"-id"); err != nil {
			t.Fatal(err)
		}
		if side != want {
			t.Errorf("got side %q, want %q", side, want)
		}
	}
	side = ""
	if err := wrapper.CancelOrder(market, "missing"); err == nil || side != "" {
		t.Fatal("unknown order should not be canceled")
	}
}

func TestCancelOrderUnavailableAndInvalid(t *testing.T) {
	for _, wrapper := range []ExchangeWrapper{&HitBtcWrapperV2{}, &BittrexWrapperV2{}} {
		if err := wrapper.CancelOrder(nil, "order"); !errors.Is(err, ErrCancelOrderNotSupported) {
			t.Errorf("%T: expected unsupported error, got %v", wrapper, err)
		}
	}
	simulator := NewExchangeWrapperSimulator(nil, nil)
	if err := simulator.CancelOrder(nil, "fake-order"); err == nil {
		t.Fatal("simulated orders cannot be canceled")
	}
	for _, wrapper := range []ExchangeWrapper{&BitfinexWrapper{}, &PoloniexWrapper{}, &BittrexWrapper{}} {
		if err := wrapper.CancelOrder(nil, "../invalid"); err == nil {
			t.Errorf("%T accepted an invalid ID", wrapper)
		}
	}
}
