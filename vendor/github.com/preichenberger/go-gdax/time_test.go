package gdax

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestGetTime(t *testing.T) {
	client := NewTestClient()
	serverTime, err := client.GetTime()
	if err != nil {
		t.Error(err)
	}

	if StructHasZeroValues(serverTime) {
		t.Error(errors.New("Zero value"))
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	c := Time{}
	now := time.Now()

	jsonData, err := json.Marshal(now.Format("2006-01-02 15:04:05+00"))
	if err != nil {
		t.Error(err)
	}

	if err = c.UnmarshalJSON(jsonData); err != nil {
		t.Error(err)
	}

	if now.Equal(c.Time()) {
		t.Error(errors.New("Unmarshaled time does not equal original time"))
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	c := Time{}
	tt := time.Date(9999, 4, 12, 23, 20, 50, 0, time.UTC)
	expected := "\"9999-04-12T23:20:50Z\""

	jsonData, err := json.Marshal(tt.Format("2006-01-02 15:04:05+00"))
	if err != nil {
		t.Error(err)
	}

	if err = c.UnmarshalJSON(jsonData); err != nil {
		t.Error(err)
	}

	jsonData, err = json.Marshal(c)
	if err != nil {
		t.Error(err)
	}

	if string(jsonData) != expected {
		t.Error(errors.New("Marshaled time (" + string(jsonData) + ") does not equal original time (" + expected + ")"))
	}
}
