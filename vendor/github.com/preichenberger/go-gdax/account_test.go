package gdax

import (
	"errors"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	client := NewTestClient()
	accounts, err := client.GetAccounts()
	if err != nil {
		t.Error(err)
	}

	// Check for decoding issues
	for _, a := range accounts {
		if StructHasZeroValues(a) {
			t.Error(errors.New("Zero value"))
		}
	}
}

func TestGetAccount(t *testing.T) {
	client := NewTestClient()
	accounts, err := client.GetAccounts()
	if err != nil {
		t.Error(err)
	}

	for _, a := range accounts {
		account, err := client.GetAccount(a.Id)
		if err != nil {
			t.Error(err)
		}

		// Check for decoding issues
		if StructHasZeroValues(account) {
			t.Error(errors.New("Zero value"))
		}
	}
}
func TestListAccountLedger(t *testing.T) {
	var ledgers []LedgerEntry
	client := NewTestClient()
	accounts, err := client.GetAccounts()
	if err != nil {
		t.Error(err)
	}

	for _, a := range accounts {
		cursor := client.ListAccountLedger(a.Id)
		for cursor.HasMore {
			if err := cursor.NextPage(&ledgers); err != nil {
				t.Error(err)
			}

			for _, ledger := range ledgers {
				props := []string{"Id", "CreatedAt", "Amount", "Balance", "Type"}
				if err := EnsureProperties(ledger, props); err != nil {
					t.Error(err)
				}

				if ledger.Type == "match" || ledger.Type == "fee" {
					if err := Ensure(ledger.Details); err != nil {
						t.Error(errors.New("Details is missing"))
					}
				}
			}
		}
	}
}

func TestListHolds(t *testing.T) {
	var holds []Hold
	client := NewTestClient()
	accounts, err := client.GetAccounts()
	if err != nil {
		t.Error(err)
	}

	for _, a := range accounts {
		cursor := client.ListHolds(a.Id)
		for cursor.HasMore {
			if err := cursor.NextPage(&holds); err != nil {
				t.Error(err)
			}

			for _, h := range holds {
				// Check for decoding issues
				if StructHasZeroValues(h) {
					t.Error(errors.New("Zero value"))
				}
			}
		}
	}
}
