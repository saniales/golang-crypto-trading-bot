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
		if form.Get("txid") != "OABC23-DEF45-GHI789" {
			t.Errorf("unexpected txid: %q", form.Get("txid"))
		}
		if count == 0 {
			return jsonResponse(`{"error":[],"result":{"count":0}}`, http.StatusOK), nil
		}
		return jsonResponse(`{"error":[],"result":{"count":1}}`, http.StatusOK), nil
	})}
	wrapper := &KrakenWrapper{api: krakenapi.NewWithClient("key", "secret", client)}
	if err := wrapper.CancelOrder(nil, "[OABC23-DEF45-GHI789]"); err != nil {
		t.Fatal(err)
	}
	count = 0
	if err := wrapper.CancelOrder(nil, "OABC23-DEF45-GHI789"); err == nil {
		t.Fatal("expected zero canceled orders to fail")
	}
	for _, id := range []string{"", "[]", "[one two]", "one,two", "123", "-42", "OABC23-DEF45-GHI789,OABC23-DEF45-GHI780"} {
		if err := wrapper.CancelOrder(nil, id); err == nil {
			t.Errorf("expected invalid ID %q to fail", id)
		}
	}
}

func TestCancelOrderUnavailableAndInvalid(t *testing.T) {
	for _, wrapper := range []ExchangeWrapper{&HitBtcWrapperV2{}, &BittrexWrapperV2{}, &KucoinWrapper{}} {
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
