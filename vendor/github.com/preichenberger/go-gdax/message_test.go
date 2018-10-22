package gdax

import (
	"errors"
	"testing"

	ws "github.com/gorilla/websocket"
)

func startSubscribe(wsConn *ws.Conn, message *Message) (*Message, error) {
	var receivedMessage Message

	if err := wsConn.WriteJSON(message); err != nil {
		return nil, err
	}

	for {
		if err := wsConn.ReadJSON(&receivedMessage); err != nil {
			return nil, err
		}

		if receivedMessage.Type != "subscriptions" {
			break
		}
	}

	return &receivedMessage, nil
}

func TestMessageHeartbeat(t *testing.T) {
	wsConn, err := NewTestWebsocketClient()
	if err != nil {
		t.Error(err)
	}
	defer wsConn.Close()

	subscribe := Message{
		Type: "subscribe",
		Channels: []MessageChannel{
			MessageChannel{
				Name: "heartbeat",
				ProductIds: []string{
					"BTC-USD",
				},
			},
		},
	}

	message, err := startSubscribe(wsConn, &subscribe)
	if err != nil {
		t.Error(err)
	}

	if message.Type != "heartbeat" {
		t.Error(errors.New("Invalid message type"))
	}

	props := []string{"Type", "Sequence", "LastTradeId", "ProductId", "Time"}
	if err := EnsureProperties(message, props); err != nil {
		t.Error(err)
	}
}

func TestMessageTicker(t *testing.T) {
	wsConn, err := NewTestWebsocketClient()
	if err != nil {
		t.Error(err)
	}
	defer wsConn.Close()

	subscribe := Message{
		Type: "subscribe",
		Channels: []MessageChannel{
			MessageChannel{
				Name: "ticker",
				ProductIds: []string{
					"BTC-USD",
				},
			},
		},
	}

	message, err := startSubscribe(wsConn, &subscribe)
	if err != nil {
		t.Error(err)
	}

	if message.Type != "ticker" {
		t.Error(errors.New("Invalid message type"))
	}

	props := []string{"Type", "Sequence", "ProductId", "BestBid", "BestAsk", "Price"}
	if err := EnsureProperties(message, props); err != nil {
		t.Error(err)
	}
}

func TestMessageLevel2(t *testing.T) {
	wsConn, err := NewTestWebsocketClient()
	if err != nil {
		t.Error(err)
	}
	defer wsConn.Close()

	subscribe := Message{
		Type: "subscribe",
		Channels: []MessageChannel{
			MessageChannel{
				Name: "level2",
				ProductIds: []string{
					"BTC-USD",
				},
			},
		},
	}

	message, err := startSubscribe(wsConn, &subscribe)
	if err != nil {
		t.Error(err)
	}

	if message.Type != "snapshot" {
		t.Error(errors.New("Invalid message type"))
	}

	props := []string{"ProductId", "Bids", "Asks"}
	if err := EnsureProperties(message, props); err != nil {
		t.Error(err)
	}

	/*
		l2 := false
		for i := 0; i < 10; i++ {
			message = &Message{}
			if err = wsConn.ReadJSON(&message); err != nil {
				t.Error(err)
			}

			if message.Type == "l2update" {
				l2 = true
				props := []string{"ProductId", "Changes"}
				if err := EnsureProperties(message, props); err != nil {
					t.Error(err)
				}

				break
			}
		}

		if !l2 {
			t.Error(errors.New("Did not find l2update"))
		}
	*/
}
