package gdax

import (
	"errors"
	"testing"
)

func TestListFills(t *testing.T) {
	var fills []Fill
	client := NewTestClient()
	params := ListFillsParams{
		ProductId: "BTC-USD",
	}
	cursor := client.ListFills(params)
	for cursor.HasMore {
		if err := cursor.NextPage(&fills); err != nil {
			t.Error(err)
		}

		for _, f := range fills {
			if StructHasZeroValues(f) {
				t.Error(errors.New("Zero value"))
			}
		}
	}
}
