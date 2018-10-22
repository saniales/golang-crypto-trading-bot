package gdax

import (
	"errors"
	"testing"
)

func TestGetProducts(t *testing.T) {
	client := NewTestClient()
	products, err := client.GetProducts()
	if err != nil {
		t.Error(err)
	}

	for _, p := range products {
		if StructHasZeroValues(p) {
			t.Error(errors.New("Zero value"))
		}
	}
}

func TestGetBook(t *testing.T) {
	client := NewTestClient()
	_, err := client.GetBook("BTC-USD", 1)
	if err != nil {
		t.Error(err)
	}
	_, err = client.GetBook("BTC-USD", 2)
	if err != nil {
		t.Error(err)
	}
	_, err = client.GetBook("BTC-USD", 3)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTicker(t *testing.T) {
	client := NewTestClient()
	ticker, err := client.GetTicker("BTC-USD")
	if err != nil {
		t.Error(err)
	}

	if StructHasZeroValues(ticker) {
		t.Error(errors.New("Zero value"))
	}

	ticker, err = client.GetTicker("ETH-BTC")
	if err != nil {
		t.Error(err)
	}

	if StructHasZeroValues(ticker) {
		t.Error(errors.New("Zero value"))
	}

}

func TestListTrades(t *testing.T) {
	var trades []Trade
	client := NewTestClient()
	cursor := client.ListTrades("BTC-USD")

	if err := cursor.NextPage(&trades); err != nil {
		t.Error(err)
	}

	for _, a := range trades {
		if StructHasZeroValues(a) {
			t.Error(errors.New("Zero value"))
		}
	}
}

func TestGetHistoricRates(t *testing.T) {
	client := NewTestClient()
	params := GetHistoricRatesParams{
		Granularity: 3600,
	}

	historicRates, err := client.GetHistoricRates("BTC-USD", params)
	if err != nil {
		t.Error(err)
	}

	props := []string{"Time", "Low", "High", "Open", "Close", "Volume"}
	if err := EnsureProperties(historicRates[0], props); err != nil {
		t.Error(err)
	}
}

func TestGetStats(t *testing.T) {
	client := NewTestClient()
	stats, err := client.GetStats("BTC-USD")
	if err != nil {
		t.Error(err)
	}

	props := []string{"Low", "Open", "Volume", "Last", "Volume_30Day"}
	if err := EnsureProperties(stats, props); err != nil {
		t.Error(err)
	}
}
